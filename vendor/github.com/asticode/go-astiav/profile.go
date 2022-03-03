package astiav

//#cgo pkg-config: libavcodec
//#include <libavcodec/avcodec.h>
import "C"

type Profile int

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavcodec/avcodec.h#L1526
const (
	ProfileAacEld                            = Profile(C.FF_PROFILE_AAC_ELD)
	ProfileAacHe                             = Profile(C.FF_PROFILE_AAC_HE)
	ProfileAacHeV2                           = Profile(C.FF_PROFILE_AAC_HE_V2)
	ProfileAacLd                             = Profile(C.FF_PROFILE_AAC_LD)
	ProfileAacLow                            = Profile(C.FF_PROFILE_AAC_LOW)
	ProfileAacLtp                            = Profile(C.FF_PROFILE_AAC_LTP)
	ProfileAacMain                           = Profile(C.FF_PROFILE_AAC_MAIN)
	ProfileAacSsr                            = Profile(C.FF_PROFILE_AAC_SSR)
	ProfileAv1High                           = Profile(C.FF_PROFILE_AV1_HIGH)
	ProfileAv1Main                           = Profile(C.FF_PROFILE_AV1_MAIN)
	ProfileAv1Professional                   = Profile(C.FF_PROFILE_AV1_PROFESSIONAL)
	ProfileDnxhd                             = Profile(C.FF_PROFILE_DNXHD)
	ProfileDnxhr444                          = Profile(C.FF_PROFILE_DNXHR_444)
	ProfileDnxhrHq                           = Profile(C.FF_PROFILE_DNXHR_HQ)
	ProfileDnxhrHqx                          = Profile(C.FF_PROFILE_DNXHR_HQX)
	ProfileDnxhrLb                           = Profile(C.FF_PROFILE_DNXHR_LB)
	ProfileDnxhrSq                           = Profile(C.FF_PROFILE_DNXHR_SQ)
	ProfileDts                               = Profile(C.FF_PROFILE_DTS)
	ProfileDts9624                           = Profile(C.FF_PROFILE_DTS_96_24)
	ProfileDtsEs                             = Profile(C.FF_PROFILE_DTS_ES)
	ProfileDtsExpress                        = Profile(C.FF_PROFILE_DTS_EXPRESS)
	ProfileDtsHdHra                          = Profile(C.FF_PROFILE_DTS_HD_HRA)
	ProfileDtsHdMa                           = Profile(C.FF_PROFILE_DTS_HD_MA)
	ProfileH264Baseline                      = Profile(C.FF_PROFILE_H264_BASELINE)
	ProfileH264Cavlc444                      = Profile(C.FF_PROFILE_H264_CAVLC_444)
	ProfileH264Constrained                   = Profile(C.FF_PROFILE_H264_CONSTRAINED)
	ProfileH264ConstrainedBaseline           = Profile(C.FF_PROFILE_H264_CONSTRAINED_BASELINE)
	ProfileH264Extended                      = Profile(C.FF_PROFILE_H264_EXTENDED)
	ProfileH264High                          = Profile(C.FF_PROFILE_H264_HIGH)
	ProfileH264High10                        = Profile(C.FF_PROFILE_H264_HIGH_10)
	ProfileH264High10Intra                   = Profile(C.FF_PROFILE_H264_HIGH_10_INTRA)
	ProfileH264High422                       = Profile(C.FF_PROFILE_H264_HIGH_422)
	ProfileH264High422Intra                  = Profile(C.FF_PROFILE_H264_HIGH_422_INTRA)
	ProfileH264High444                       = Profile(C.FF_PROFILE_H264_HIGH_444)
	ProfileH264High444Intra                  = Profile(C.FF_PROFILE_H264_HIGH_444_INTRA)
	ProfileH264High444Predictive             = Profile(C.FF_PROFILE_H264_HIGH_444_PREDICTIVE)
	ProfileH264Intra                         = Profile(C.FF_PROFILE_H264_INTRA)
	ProfileH264Main                          = Profile(C.FF_PROFILE_H264_MAIN)
	ProfileH264MultiviewHigh                 = Profile(C.FF_PROFILE_H264_MULTIVIEW_HIGH)
	ProfileH264StereoHigh                    = Profile(C.FF_PROFILE_H264_STEREO_HIGH)
	ProfileHevcMain                          = Profile(C.FF_PROFILE_HEVC_MAIN)
	ProfileHevcMain10                        = Profile(C.FF_PROFILE_HEVC_MAIN_10)
	ProfileHevcMainStillPicture              = Profile(C.FF_PROFILE_HEVC_MAIN_STILL_PICTURE)
	ProfileHevcRext                          = Profile(C.FF_PROFILE_HEVC_REXT)
	ProfileJpeg2000CstreamNoRestriction      = Profile(C.FF_PROFILE_JPEG2000_CSTREAM_NO_RESTRICTION)
	ProfileJpeg2000CstreamRestriction0       = Profile(C.FF_PROFILE_JPEG2000_CSTREAM_RESTRICTION_0)
	ProfileJpeg2000CstreamRestriction1       = Profile(C.FF_PROFILE_JPEG2000_CSTREAM_RESTRICTION_1)
	ProfileJpeg2000Dcinema2K                 = Profile(C.FF_PROFILE_JPEG2000_DCINEMA_2K)
	ProfileJpeg2000Dcinema4K                 = Profile(C.FF_PROFILE_JPEG2000_DCINEMA_4K)
	ProfileMjpegHuffmanBaselineDct           = Profile(C.FF_PROFILE_MJPEG_HUFFMAN_BASELINE_DCT)
	ProfileMjpegHuffmanExtendedSequentialDct = Profile(C.FF_PROFILE_MJPEG_HUFFMAN_EXTENDED_SEQUENTIAL_DCT)
	ProfileMjpegHuffmanLossless              = Profile(C.FF_PROFILE_MJPEG_HUFFMAN_LOSSLESS)
	ProfileMjpegHuffmanProgressiveDct        = Profile(C.FF_PROFILE_MJPEG_HUFFMAN_PROGRESSIVE_DCT)
	ProfileMjpegJpegLs                       = Profile(C.FF_PROFILE_MJPEG_JPEG_LS)
	ProfileMpeg2422                          = Profile(C.FF_PROFILE_MPEG2_422)
	ProfileMpeg2AacHe                        = Profile(C.FF_PROFILE_MPEG2_AAC_HE)
	ProfileMpeg2AacLow                       = Profile(C.FF_PROFILE_MPEG2_AAC_LOW)
	ProfileMpeg2High                         = Profile(C.FF_PROFILE_MPEG2_HIGH)
	ProfileMpeg2Main                         = Profile(C.FF_PROFILE_MPEG2_MAIN)
	ProfileMpeg2Simple                       = Profile(C.FF_PROFILE_MPEG2_SIMPLE)
	ProfileMpeg2SnrScalable                  = Profile(C.FF_PROFILE_MPEG2_SNR_SCALABLE)
	ProfileMpeg2Ss                           = Profile(C.FF_PROFILE_MPEG2_SS)
	ProfileMpeg4AdvancedCoding               = Profile(C.FF_PROFILE_MPEG4_ADVANCED_CODING)
	ProfileMpeg4AdvancedCore                 = Profile(C.FF_PROFILE_MPEG4_ADVANCED_CORE)
	ProfileMpeg4AdvancedRealTime             = Profile(C.FF_PROFILE_MPEG4_ADVANCED_REAL_TIME)
	ProfileMpeg4AdvancedScalableTexture      = Profile(C.FF_PROFILE_MPEG4_ADVANCED_SCALABLE_TEXTURE)
	ProfileMpeg4AdvancedSimple               = Profile(C.FF_PROFILE_MPEG4_ADVANCED_SIMPLE)
	ProfileMpeg4BasicAnimatedTexture         = Profile(C.FF_PROFILE_MPEG4_BASIC_ANIMATED_TEXTURE)
	ProfileMpeg4Core                         = Profile(C.FF_PROFILE_MPEG4_CORE)
	ProfileMpeg4CoreScalable                 = Profile(C.FF_PROFILE_MPEG4_CORE_SCALABLE)
	ProfileMpeg4Hybrid                       = Profile(C.FF_PROFILE_MPEG4_HYBRID)
	ProfileMpeg4Main                         = Profile(C.FF_PROFILE_MPEG4_MAIN)
	ProfileMpeg4NBit                         = Profile(C.FF_PROFILE_MPEG4_N_BIT)
	ProfileMpeg4ScalableTexture              = Profile(C.FF_PROFILE_MPEG4_SCALABLE_TEXTURE)
	ProfileMpeg4Simple                       = Profile(C.FF_PROFILE_MPEG4_SIMPLE)
	ProfileMpeg4SimpleFaceAnimation          = Profile(C.FF_PROFILE_MPEG4_SIMPLE_FACE_ANIMATION)
	ProfileMpeg4SimpleScalable               = Profile(C.FF_PROFILE_MPEG4_SIMPLE_SCALABLE)
	ProfileMpeg4SimpleStudio                 = Profile(C.FF_PROFILE_MPEG4_SIMPLE_STUDIO)
	ProfileReserved                          = Profile(C.FF_PROFILE_RESERVED)
	ProfileSbcMsbc                           = Profile(C.FF_PROFILE_SBC_MSBC)
	ProfileUnknown                           = Profile(C.FF_PROFILE_UNKNOWN)
	ProfileVc1Advanced                       = Profile(C.FF_PROFILE_VC1_ADVANCED)
	ProfileVc1Complex                        = Profile(C.FF_PROFILE_VC1_COMPLEX)
	ProfileVc1Main                           = Profile(C.FF_PROFILE_VC1_MAIN)
	ProfileVc1Simple                         = Profile(C.FF_PROFILE_VC1_SIMPLE)
	ProfileVp90                              = Profile(C.FF_PROFILE_VP9_0)
	ProfileVp91                              = Profile(C.FF_PROFILE_VP9_1)
	ProfileVp92                              = Profile(C.FF_PROFILE_VP9_2)
	ProfileVp93                              = Profile(C.FF_PROFILE_VP9_3)
)