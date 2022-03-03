package astiav

//#cgo pkg-config: libavformat
//#include <libavformat/avformat.h>
import "C"

type IOContextFlag int

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavformat/avio.h#L621
const (
	IOContextFlagRead      = IOContextFlag(C.AVIO_FLAG_READ)
	IOContextFlagWrite     = IOContextFlag(C.AVIO_FLAG_WRITE)
	IOContextFlagReadWrite = IOContextFlag(C.AVIO_FLAG_READ_WRITE)
	IOContextFlagNonBlock  = IOContextFlag(C.AVIO_FLAG_NONBLOCK)
	IOContextFlagDirect    = IOContextFlag(C.AVIO_FLAG_DIRECT)
)
