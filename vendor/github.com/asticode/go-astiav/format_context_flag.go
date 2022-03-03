package astiav

//#cgo pkg-config: libavformat
//#include <libavformat/avformat.h>
import "C"

type FormatContextFlag int

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavformat/avformat.h#L1321
const (
	FormatContextFlagGenPts         = FormatContextFlag(C.AVFMT_FLAG_GENPTS)
	FormatContextFlagIgnidx         = FormatContextFlag(C.AVFMT_FLAG_IGNIDX)
	FormatContextFlagNonblock       = FormatContextFlag(C.AVFMT_FLAG_NONBLOCK)
	FormatContextFlagIgnDts         = FormatContextFlag(C.AVFMT_FLAG_IGNDTS)
	FormatContextFlagNofillin       = FormatContextFlag(C.AVFMT_FLAG_NOFILLIN)
	FormatContextFlagNoparse        = FormatContextFlag(C.AVFMT_FLAG_NOPARSE)
	FormatContextFlagNobuffer       = FormatContextFlag(C.AVFMT_FLAG_NOBUFFER)
	FormatContextFlagCustomIo       = FormatContextFlag(C.AVFMT_FLAG_CUSTOM_IO)
	FormatContextFlagDiscardCorrupt = FormatContextFlag(C.AVFMT_FLAG_DISCARD_CORRUPT)
	FormatContextFlagFlushPackets   = FormatContextFlag(C.AVFMT_FLAG_FLUSH_PACKETS)
	FormatContextFlagBitexact       = FormatContextFlag(C.AVFMT_FLAG_BITEXACT)
	FormatContextFlagMp4ALatm       = FormatContextFlag(C.AVFMT_FLAG_MP4A_LATM)
	FormatContextFlagSortDts        = FormatContextFlag(C.AVFMT_FLAG_SORT_DTS)
	FormatContextFlagPrivOpt        = FormatContextFlag(C.AVFMT_FLAG_PRIV_OPT)
	FormatContextFlagKeepSideData   = FormatContextFlag(C.AVFMT_FLAG_KEEP_SIDE_DATA)
	FormatContextFlagFastSeek       = FormatContextFlag(C.AVFMT_FLAG_FAST_SEEK)
	FormatContextFlagShortest       = FormatContextFlag(C.AVFMT_FLAG_SHORTEST)
	FormatContextFlagAutoBsf        = FormatContextFlag(C.AVFMT_FLAG_AUTO_BSF)
)
