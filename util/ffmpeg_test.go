package util

import (
	"fmt"
	"strings"
	"testing"
)

func TestFFmpegAppendToFile(t *testing.T) {
	filePath := GetCurrentDirectory() + `/001.wav`
	ffmpeg := Ffmpeg{FilePath: filePath}
	ffmpeg.GetInfo()
	fmt.Println(ffmpeg.WavInfo)

	str := ExeFun("ffmpeg", "-i", filePath)
	str = str[strings.Index(str, "Input #0,"):]
	strArr := strings.Split(str, "\r\n")
	fmt.Println(strArr)
	str = strArr[1]
	//  Duration: 00:01:11.40, bitrate: 128 kb/s
	str = strings.ReplaceAll(str, " ", "")
	str = strings.ReplaceAll(str, "Duration:", "")
	str = strings.ReplaceAll(str, "bitrate:", "")
	arr := strings.Split(str, ",")
	fmt.Println(arr)

	/*str := `ffmpeg version N-94300-gaf5f770113 Copyright (c) 2000-2019 the FFmpeg developers
	  built with gcc 9.1.1 (GCC) 20190621
	  configuration: --disable-static --enable-shared --enable-gpl --enable-version3 --enable-sdl2 --enable-fontconfig --enable-gnutls --enable-iconv --enable-libass --enable-libdav1d --enable-libbluray --enable-libfreetype --enable-libmp3lame --enable-libopencore-amrnb --enable-libopencore-amrwb --enable-libopenjpeg --enable-libopus --enable-libshine --enable-libsnappy --enable-libsoxr --enable-libtheora --enable-libtwolame --enable-libvpx --enable-libwavpack --enable-libwebp --enable-libx264 --enable-libx265 --enable-libxml2 --enable-libzimg --enable-lzma --enable-zlib --enable-gmp --enable-libvidstab --enable-libvorbis --enable-libvo-amrwbenc --enable-libmysofa --enable-libspeex --enable-libxvid --enable-libaom --enable-libmfx --enable-amf --enable-ffnvcodec --enable-cuvid --enable-d3d11va --enable-nvenc --enable-nvdec --enable-dxva2 --enable-avisynth --enable-libopenmpt
	  libavutil      56. 30.100 / 56. 30.100
	  libavcodec     58. 53.101 / 58. 53.101
	  libavformat    58. 28.102 / 58. 28.102
	  libavdevice    58.  7.100 / 58.  7.100
	  libavfilter     7. 56.101 /  7. 56.101
	  libswscale      5.  4.101 /  5.  4.101
	  libswresample   3.  4.100 /  3.  4.100
	  libpostproc    55.  4.100 / 55.  4.100
	Guessed Channel Layout for Input Stream #0.0 : mono
	Input #0, wav, from 'C:/Users/Administrator/Desktop/20190727/001.wav':
	  Duration: 00:01:11.40, bitrate: 128 kb/s
	    Stream #0:0: Audio: pcm_s16le ([1][0][0][0] / 0x0001), 8000 Hz, mono, s16, 128 kb/s
	At least one output file must be specified`

		str1 :=`ffmpeg version 4.2.1 Copyright (c) 2000-2019 the FFmpeg developers
	  built with gcc 4.8.5 (GCC) 20150623 (Red Hat 4.8.5-39)
	  configuration: --enable-shared --prefix=/usr/local/ffmpeg --disable-yasm
	  libavutil      56. 31.100 / 56. 31.100
	  libavcodec     58. 54.100 / 58. 54.100
	  libavformat    58. 29.100 / 58. 29.100
	  libavdevice    58.  8.100 / 58.  8.100
	  libavfilter     7. 57.100 /  7. 57.100
	  libswscale      5.  5.100 /  5.  5.100
	  libswresample   3.  5.100 /  3.  5.100
	Guessed Channel Layout for Input Stream #0.0 : mono
	Input #0, wav, from '001.wav':
	  Duration: 00:01:11.40, bitrate: 128 kb/s
	    Stream #0:0: Audio: pcm_s16le ([1][0][0][0] / 0x0001), 8000 Hz, mono, s16, 128 kb/s`*/
}
