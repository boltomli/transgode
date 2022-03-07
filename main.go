package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/asticode/go-astiav"
	"github.com/asticode/go-astikit"
	"github.com/gofiber/fiber/v2"
)

type stream struct {
	buffersinkContext *astiav.FilterContext
	buffersrcContext  *astiav.FilterContext
	decCodec          *astiav.Codec
	decCodecContext   *astiav.CodecContext
	decFrame          *astiav.Frame
	encCodec          *astiav.Codec
	encCodecContext   *astiav.CodecContext
	encPkt            *astiav.Packet
	filterFrame       *astiav.Frame
	filterGraph       *astiav.FilterGraph
	inputStream       *astiav.Stream
	outputStream      *astiav.Stream
}

var (
	supportedEncCodecs = make(map[string]string)
)

type TranscodeTask struct {
	AudioUrl   string `form:"audiourl"`
	MediaType  string `form:"mediatype"`
	Channels   int    `form:"channels"`
	SampleRate int    `form:"samplerate"`
	Success    bool
	Status     int
	Message    string `default:""`
}

func main() {
	// Handle ffmpeg logs
	astiav.SetLogLevel(astiav.LogLevelDebug)
	astiav.SetLogCallback(func(l astiav.LogLevel, msg, parent string) {
		log.Printf("ffmpeg log: %s (level: %d)\n", strings.TrimSpace(msg), l)
	})

	supportedEncCodecs = map[string]string{
		"wav": "pcm_s16le",
		"raw": "pcm_s16le",
	}

	app := fiber.New()
	app.Post("/speak/transcode", func(ct *fiber.Ctx) (err error) {
		task := new(TranscodeTask)

		if err := ct.BodyParser(task); err != nil {
			return ct.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		// default to stereo
		if task.Channels < 1 {
			task.Channels = 2
		}
		if task.Channels > 2 {
			task.Channels = 2
		}

		// default to 44100
		if task.SampleRate < 16000 {
			task.SampleRate = 44100
		}
		if task.SampleRate > 48000 {
			task.SampleRate = 48000
		}

		task.Success = false
		task.Status = http.StatusOK

		// support only PCM for now
		if v := supportedEncCodecs[task.MediaType]; v == "" {
			task.Message = fmt.Sprintf("main: codec not supported: %s", task.MediaType)
			task.Status = http.StatusUnsupportedMediaType
			return ct.JSON(task)
		}

		var (
			c                   = astikit.NewCloser()
			inputFormatContext  *astiav.FormatContext
			outputFormatContext *astiav.FormatContext
			streams             = make(map[int]*stream) // Indexed by input stream index
		)

		// We use an astikit.Closer to free all resources properly
		defer c.Close()

		// Open input file
		// Alloc input format context
		if inputFormatContext = astiav.AllocFormatContext(); inputFormatContext == nil {
			task.Message = fmt.Sprintf("main: input format context is nil")
			task.Status = http.StatusBadRequest
			return ct.JSON(task)
		}
		c.Add(inputFormatContext.Free)

		// Open input
		if err = inputFormatContext.OpenInput(task.AudioUrl, nil, nil); err != nil {
			task.Message = fmt.Sprintf("main: opening input failed: %s", err)
			task.Status = http.StatusBadRequest
			return ct.JSON(task)
		}
		c.Add(inputFormatContext.CloseInput)

		// Find stream info
		if err = inputFormatContext.FindStreamInfo(nil); err != nil {
			task.Message = fmt.Sprintf("main: finding stream info failed: %w", err)
			task.Status = http.StatusBadRequest
			return ct.JSON(task)
		}

		// Loop through streams
		for _, is := range inputFormatContext.Streams() {
			// Only process audio
			if is.CodecParameters().MediaType() != astiav.MediaTypeAudio {
				continue
			}

			// Create stream
			s := &stream{inputStream: is}

			// Find decoder
			if s.decCodec = astiav.FindDecoder(is.CodecParameters().CodecID()); s.decCodec == nil {
				err = errors.New("main: codec is nil")
				task.Status = http.StatusBadRequest
				return ct.JSON(task)
			}

			// Alloc codec context
			if s.decCodecContext = astiav.AllocCodecContext(s.decCodec); s.decCodecContext == nil {
				err = errors.New("main: codec context is nil")
				task.Status = http.StatusBadRequest
				return ct.JSON(task)
			}
			c.Add(s.decCodecContext.Free)

			// Update codec context
			if err = is.CodecParameters().ToCodecContext(s.decCodecContext); err != nil {
				task.Message = fmt.Sprintf("main: updating codec context failed: %s", err)
				task.Status = http.StatusBadRequest
				return ct.JSON(task)
			}

			// Set framerate
			if is.CodecParameters().MediaType() == astiav.MediaTypeVideo {
				s.decCodecContext.SetFramerate(inputFormatContext.GuessFrameRate(is, nil))
			}

			// Update channel layout
			s.decCodecContext.SetChannelLayout(astiav.ChannelLayout(channels2Layout(s.decCodecContext.Channels())))

			// Open codec context
			if err = s.decCodecContext.Open(s.decCodec, nil); err != nil {
				task.Message = fmt.Sprintf("main: opening codec context failed: %s", err)
				task.Status = http.StatusBadRequest
				return ct.JSON(task)
			}

			// Alloc frame
			s.decFrame = astiav.AllocFrame()
			c.Add(s.decFrame.Free)

			// Store stream
			streams[is.Index()] = s
		}

		// Open output file
		f, err := ioutil.TempFile("", fmt.Sprintf("transcode_*.%s", "wav"))
		defer os.Remove(f.Name())
		if err != nil {
			task.Message = fmt.Sprintf("main: get temp output file failed: %s", err)
			task.Status = http.StatusBadRequest
			return ct.JSON(task)
		}

		mediaType := strings.ToLower(task.MediaType)
		formatName := ""
		if strings.ToLower(mediaType) == "raw" {
			formatName = "data"
		}

		// Alloc output format context
		if outputFormatContext, err = astiav.AllocOutputFormatContext(nil, formatName, f.Name()); err != nil {
			task.Message = fmt.Sprintf("main: allocating output format context failed: %w", err)
			task.Status = http.StatusBadRequest
			return ct.JSON(task)
		} else if outputFormatContext == nil {
			err = errors.New("main: output format context is nil")
			task.Status = http.StatusBadRequest
			return ct.JSON(task)
		}
		c.Add(outputFormatContext.Free)

		// Loop through streams
		for _, is := range inputFormatContext.Streams() {
			// Get stream
			s, ok := streams[is.Index()]
			if !ok {
				continue
			}

			// Create output stream
			if s.outputStream = outputFormatContext.NewStream(nil); s.outputStream == nil {
				err = errors.New("main: output stream is nil")
				task.Status = http.StatusBadRequest
				return ct.JSON(task)
			}

			// Get codec for audio only
			if s.decCodecContext.MediaType() != astiav.MediaTypeAudio {
				err = errors.New("main: codec is not audio")
				task.Status = http.StatusBadRequest
				return ct.JSON(task)
			}

			encCodec := mediaType
			if v := supportedEncCodecs[mediaType]; v != "" {
				encCodec = v
			}

			// Find encoder
			if s.encCodec = astiav.FindEncoderByName(encCodec); s.encCodec == nil {
				err = errors.New("main: codec is nil")
				task.Status = http.StatusBadRequest
				return ct.JSON(task)
			}

			// Alloc codec context
			if s.encCodecContext = astiav.AllocCodecContext(s.encCodec); s.encCodecContext == nil {
				err = errors.New("main: codec context is nil")
				task.Status = http.StatusBadRequest
				return ct.JSON(task)
			}
			c.Add(s.encCodecContext.Free)

			// Update codec context
			if s.decCodecContext.MediaType() == astiav.MediaTypeAudio {
				channelLayout := astiav.ChannelLayout(channels2Layout(task.Channels))
				if v := s.encCodec.ChannelLayouts(); len(v) > 0 {
					result := false
					for _, x := range v {
						if x == channelLayout {
							result = true
							break
						}
					}
					if !result {
						err = errors.New("main: codec not support channel layout " + channelLayout.String())
						task.Status = http.StatusBadRequest
						return ct.JSON(task)
					}
				}
				s.encCodecContext.SetChannelLayout(channelLayout)
				s.encCodecContext.SetChannels(task.Channels)
				s.encCodecContext.SetSampleRate(task.SampleRate)

				sampleFormat := s.decCodecContext.SampleFormat()
				if v := s.encCodec.SampleFormats(); len(v) > 0 {
					result := false
					for _, x := range v {
						if x == sampleFormat {
							result = true
							break
						}
					}
					if !result {
						sampleFormat = v[0]
					}
				}
				s.encCodecContext.SetSampleFormat(sampleFormat)
				s.encCodecContext.SetTimeBase(s.decCodecContext.TimeBase())
			} else {
				s.encCodecContext.SetHeight(s.decCodecContext.Height())
				if v := s.encCodec.PixelFormats(); len(v) > 0 {
					s.encCodecContext.SetPixelFormat(v[0])
				} else {
					s.encCodecContext.SetPixelFormat(s.decCodecContext.PixelFormat())
				}
				s.encCodecContext.SetSampleAspectRatio(s.decCodecContext.SampleAspectRatio())
				s.encCodecContext.SetTimeBase(s.decCodecContext.TimeBase())
				s.encCodecContext.SetWidth(s.decCodecContext.Width())
			}

			// Update flags
			if s.decCodecContext.Flags().Has(astiav.CodecContextFlagGlobalHeader) {
				s.encCodecContext.SetFlags(s.encCodecContext.Flags().Add(astiav.CodecContextFlagGlobalHeader))
			}

			// Open codec context
			if err = s.encCodecContext.Open(s.encCodec, nil); err != nil {
				task.Message = fmt.Sprintf("main: opening codec context failed: %s", err)
				task.Status = http.StatusBadRequest
				return ct.JSON(task)
			}

			// Update codec parameters
			if err = s.outputStream.CodecParameters().FromCodecContext(s.encCodecContext); err != nil {
				task.Message = fmt.Sprintf("main: updating codec parameters failed: %s", err)
				task.Status = http.StatusBadRequest
				return ct.JSON(task)
			}

			// Update stream
			s.outputStream.SetTimeBase(s.encCodecContext.TimeBase())
		}

		// If this is a file, we need to use an io context
		if !outputFormatContext.OutputFormat().Flags().Has(astiav.IOFormatFlagNofile) {
			// Create io context
			ioContext := astiav.NewIOContext()

			// Open io context
			if err = ioContext.Open(f.Name(), astiav.NewIOContextFlags(astiav.IOContextFlagWrite)); err != nil {
				task.Message = fmt.Sprintf("main: opening io context failed: %s", err)
				task.Status = http.StatusBadRequest
				return ct.JSON(task)
			}
			c.AddWithError(ioContext.Closep)

			// Update output format context
			outputFormatContext.SetPb(ioContext)
		}

		// Write header
		if err = outputFormatContext.WriteHeader(nil); err != nil {
			task.Message = fmt.Sprintf("main: writing header failed: %s", err)
			task.Status = http.StatusBadRequest
			return ct.JSON(task)
		}

		// Init filters
		// Loop through output streams
		for _, s := range streams {
			// Alloc graph
			if s.filterGraph = astiav.AllocFilterGraph(); s.filterGraph == nil {
				err = errors.New("main: graph is nil")
				task.Status = http.StatusBadRequest
				return ct.JSON(task)
			}
			c.Add(s.filterGraph.Free)

			// Alloc outputs
			outputs := astiav.AllocFilterInOut()
			if outputs == nil {
				err = errors.New("main: outputs is nil")
				task.Status = http.StatusBadRequest
				return ct.JSON(task)
			}
			c.Add(outputs.Free)

			// Alloc inputs
			inputs := astiav.AllocFilterInOut()
			if inputs == nil {
				err = errors.New("main: inputs is nil")
				task.Status = http.StatusBadRequest
				return ct.JSON(task)
			}
			c.Add(inputs.Free)

			// Support only audio type
			args := astiav.FilterArgs{
				"channel_layout": s.decCodecContext.ChannelLayout().String(),
				"sample_fmt":     s.decCodecContext.SampleFormat().Name(),
				"sample_rate":    strconv.Itoa(s.decCodecContext.SampleRate()),
				"time_base":      s.decCodecContext.TimeBase().String(),
			}
			buffersrc := astiav.FindFilterByName("abuffer")
			buffersink := astiav.FindFilterByName("abuffersink")
			content := fmt.Sprintf("aresample=isr=%d:osr=%d:icl=%s:ocl=%s:isf=%s:osf=%s", s.decCodecContext.SampleRate(), s.encCodecContext.SampleRate(), s.decCodecContext.ChannelLayout().String(), s.encCodecContext.ChannelLayout().String(), s.decCodecContext.SampleFormat().Name(), s.encCodecContext.SampleFormat().Name())

			// Check filters
			if buffersrc == nil {
				err = errors.New("main: buffersrc is nil")
				task.Status = http.StatusBadRequest
				return ct.JSON(task)
			}
			if buffersink == nil {
				err = errors.New("main: buffersink is nil")
				task.Status = http.StatusBadRequest
				return ct.JSON(task)
			}

			// Create filter contexts
			if s.buffersrcContext, err = s.filterGraph.NewFilterContext(buffersrc, "in", args); err != nil {
				task.Message = fmt.Sprintf("main: creating buffersrc context failed: %s", err)
				task.Status = http.StatusBadRequest
				return ct.JSON(task)
			}
			if s.buffersinkContext, err = s.filterGraph.NewFilterContext(buffersink, "in", nil); err != nil {
				task.Message = fmt.Sprintf("main: creating buffersink context failed: %s", err)
				task.Status = http.StatusBadRequest
				return ct.JSON(task)
			}

			// Update outputs
			outputs.SetName("in")
			outputs.SetFilterContext(s.buffersrcContext)
			outputs.SetPadIdx(0)
			outputs.SetNext(nil)

			// Update inputs
			inputs.SetName("out")
			inputs.SetFilterContext(s.buffersinkContext)
			inputs.SetPadIdx(0)
			inputs.SetNext(nil)

			// Parse
			if err = s.filterGraph.Parse(content, inputs, outputs); err != nil {
				task.Message = fmt.Sprintf("main: parsing filter failed: %w", err)
				task.Status = http.StatusBadRequest
				return ct.JSON(task)
			}

			// Configure
			if err = s.filterGraph.Configure(); err != nil {
				task.Message = fmt.Sprintf("main: configuring filter failed: %w", err)
				task.Status = http.StatusBadRequest
				return ct.JSON(task)
			}

			// Alloc frame
			s.filterFrame = astiav.AllocFrame()
			c.Add(s.filterFrame.Free)

			// Alloc packet
			s.encPkt = astiav.AllocPacket()
			c.Add(s.encPkt.Free)
		}

		// Alloc packet
		pkt := astiav.AllocPacket()
		c.Add(pkt.Free)

		// Loop through packets
		for {
			// Read frame
			if err := inputFormatContext.ReadFrame(pkt); err != nil {
				if errors.Is(err, astiav.ErrEof) {
					break
				}
				task.Message = fmt.Sprintf("main: reading frame failed: %s", err)
				task.Status = http.StatusBadRequest
				return ct.JSON(task)
			}

			// Get stream
			s, ok := streams[pkt.StreamIndex()]
			if !ok {
				continue
			}

			// Update packet
			pkt.RescaleTs(s.inputStream.TimeBase(), s.decCodecContext.TimeBase())

			// Send packet
			if err := s.decCodecContext.SendPacket(pkt); err != nil {
				task.Message = fmt.Sprintf("main: sending packet failed: %s", err)
				task.Status = http.StatusBadRequest
				return ct.JSON(task)
			}

			// Loop
			for {
				// Receive frame
				if err := s.decCodecContext.ReceiveFrame(s.decFrame); err != nil {
					if errors.Is(err, astiav.ErrEof) || errors.Is(err, astiav.ErrEagain) {
						break
					}
					task.Message = fmt.Sprintf("main: receiving frame failed: %s", err)
					task.Status = http.StatusBadRequest
					return ct.JSON(task)
				}

				// Filter, encode and write frame
				if err := filterEncodeWriteFrame(s.decFrame, s, outputFormatContext); err != nil {
					task.Message = fmt.Sprintf("main: filtering, encoding and writing frame failed: %s", err)
					task.Status = http.StatusBadRequest
					return ct.JSON(task)
				}
			}
		}

		// Loop through streams
		for _, s := range streams {
			// Flush filter
			if err := filterEncodeWriteFrame(nil, s, outputFormatContext); err != nil {
				task.Message = fmt.Sprintf("main: filtering, encoding and writing frame failed: %s", err)
				task.Status = http.StatusBadRequest
				return ct.JSON(task)
			}

			// Flush encoder
			if err := encodeWriteFrame(nil, s, outputFormatContext); err != nil {
				task.Message = fmt.Sprintf("main: encoding and writing frame failed: %s", err)
				task.Status = http.StatusBadRequest
				return ct.JSON(task)
			}
		}

		// Write trailer
		if err := outputFormatContext.WriteTrailer(); err != nil {
			task.Message = fmt.Sprintf("main: writing trailer failed: %s", err)
			task.Status = http.StatusBadRequest
			return ct.JSON(task)
		}

		// Success
		task.Success = true
		return ct.SendFile(f.Name(), true)
	})
	app.Listen(":8080")
}

func filterEncodeWriteFrame(f *astiav.Frame, s *stream, outputFormatContext *astiav.FormatContext) (err error) {
	// Add frame
	if err = s.buffersrcContext.BuffersrcAddFrame(f, astiav.NewBuffersrcFlags(astiav.BuffersrcFlagKeepRef)); err != nil {
		err = fmt.Errorf("main: adding frame failed: %w", err)
		return
	}

	// Loop
	for {
		// Unref frame
		s.filterFrame.Unref()

		// Get frame
		if err = s.buffersinkContext.BuffersinkGetFrame(s.filterFrame, astiav.NewBuffersinkFlags()); err != nil {
			if errors.Is(err, astiav.ErrEof) || errors.Is(err, astiav.ErrEagain) {
				err = nil
				break
			}
			err = fmt.Errorf("main: getting frame failed: %w", err)
			return
		}

		// Reset picture type
		s.filterFrame.SetPictureType(astiav.PictureTypeNone)

		// Encode and write frame
		if err = encodeWriteFrame(s.filterFrame, s, outputFormatContext); err != nil {
			err = fmt.Errorf("main: encoding and writing frame failed: %w", err)
			return
		}
	}
	return
}

func encodeWriteFrame(f *astiav.Frame, s *stream, outputFormatContext *astiav.FormatContext) (err error) {
	// Unref packet
	s.encPkt.Unref()

	// Send frame
	if err = s.encCodecContext.SendFrame(f); err != nil {
		err = fmt.Errorf("main: sending frame failed: %w", err)
		return
	}

	// Loop
	for {
		// Receive packet
		if err = s.encCodecContext.ReceivePacket(s.encPkt); err != nil {
			if errors.Is(err, astiav.ErrEof) || errors.Is(err, astiav.ErrEagain) {
				err = nil
				break
			}
			err = fmt.Errorf("main: receiving packet failed: %w", err)
			return
		}

		// Update pkt
		s.encPkt.SetStreamIndex(s.outputStream.Index())
		s.encPkt.RescaleTs(s.encCodecContext.TimeBase(), s.outputStream.TimeBase())

		// Write frame
		if err = outputFormatContext.WriteInterleavedFrame(s.encPkt); err != nil {
			err = fmt.Errorf("main: writing frame failed: %w", err)
			return
		}
	}
	return
}

func channels2Layout(channels int) uint64 {
	if channels == 1 {
		// mono (0x4)
		return 4
	} else {
		// left (0x1) + right (0x2)
		return 3
	}
}
