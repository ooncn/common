package util

import (
	"errors"
	"runtime"
	"strings"
)

type Wav struct {
	Id                 string `json:"id"`                 // id String 主键ID
	From               string `json:"from"`               //文件绝对路径
	Encoder            string `json:"encoder"`            //编码器
	Duration           string `json:"duration"`           //时长
	DurationTimeMillis int64  `json:"durationTimeMillis"` //时长 毫秒
	Start              string `json:"start"`              //开始
	Bitrate            string `json:"bitrate"`            //比特率
	Audio              string `json:"audio"`              //音频
	TimeLong           int64  `json:"timeLong"`           // timeLong String 时长
	StartTime          string `json:"startTime"`          // startTime String 开始
	EndTime            string `json:"endTime"`            // endTime String 结束
	StartTimeLong      int64  `json:"startTimeLong"`      // startTime String 开始
	EndTimeLong        int64  `json:"endTimeLong"`        // endTime String 结束
}

type Ffmpeg struct {
	Exe      string //执行路径
	FilePath string // 文件路径
	WavInfo  string // 音频内容
	Wav      Wav    // 音频对象
}

// 当前系统
func (f Ffmpeg) Os() string {
	return runtime.GOOS
}

// 当前系统是否是Windows
func (f Ffmpeg) IsWind() bool {
	return f.Os() == "windows"
}

// 执行命令
func (f *Ffmpeg) Query(cmd ...string) (str string, err error) {
	if !f.IsWind() {
		// 当前系统是Linux时、将cmd命令语句更换成双引号
		for k, v := range cmd {
			cmd[k] = strings.ReplaceAll(v, "\"", "'")
		}
		s := strings.Join(cmd, " ")
		str = ExeFun("/bin/sh", "-c", " ffmpeg "+s)
	} else {
		for k, v := range cmd {
			cmd[k] = strings.ReplaceAll(v, "\"", "")
		}
		str = ExeFun("ffmpeg", cmd...)
	}
	i := strings.Index(str, "Input #0,")
	if i < 0 {
		err = errors.New(str)
		return
	}
	str = str[i:]
	return
}

// 根据开始时间(start)和时间长度（size）截取音频到指定的路径（toPath）
func (f *Ffmpeg) CutStartAndSizeToPath(start, size, toPath string) (str string, err error) {
	return f.Query("-i", "\""+f.FilePath+"\"", "-ss", start, "-t", size, "-y", "\""+toPath+"\"")
}
func (f *Ffmpeg) GetInfo() (*Wav, error) {
	path := f.FilePath
	str, err := f.Query("-i", "\""+path+"\"")
	if err != nil {
		return nil, errors.New("失败！" + err.Error())
	}
	w := ToWav(str)
	w.From = path
	return w, nil
}
func (f Ffmpeg) GetInfoPath(path string) (*Wav, error) {
	str, err := f.Query("-i", "\""+path+"\"")
	if err != nil {
		return nil, errors.New("失败！" + err.Error())
	}
	w := ToWav(str)
	w.From = path
	return w, nil
}

func ToWav(str string) *Wav {
	str = str[strings.Index(str, "Input #0,"):]
	strArr := strings.Split(str, "\r\n")
	str = strArr[1]
	//  Duration: 00:01:11.40, bitrate: 128 kb/s
	str = strings.ReplaceAll(str, " ", "")
	str = strings.ReplaceAll(str, "Duration:", "")
	str = strings.ReplaceAll(str, "bitrate:", "")
	arr := strings.Split(str, ",")
	duration := arr[0]
	durationTimeMillis, _ := TimeUtil.DateToTimestamp(duration)
	return &Wav{
		Duration:           duration,
		DurationTimeMillis: durationTimeMillis,
		Bitrate:            arr[1],
		Audio:              strArr[2],
	}
}
