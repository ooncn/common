package util

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

type TimeObj time.Time

var TimeUtil TimeObj

func Day() time.Duration {
	return 24 * time.Hour
} // 一天
func Month() time.Duration {
	return time.Duration(TimeObj{}.Days(time.Now())) * Day()
} // 当月多少天
func Year() time.Duration {
	y := time.Now().Year()
	if y%400 == 0 || (y%4 == 0 && y%100 != 0) {
		return 366 * Day()
	} else {
		return 365 * Day()
	}
} // 一年时间

//获取time格式化

//yyyy-MM-dd HH:mm:ss
//2006-01-02 15:04:05
func (t *TimeObj) Sep(sep string) string {
	//获取当前时间
	nt := time.Now()
	sep = strings.ReplaceAll(sep, "yyyy", "2006")
	sep = strings.ReplaceAll(sep, "MM", "01")
	sep = strings.ReplaceAll(sep, "dd", "02")
	sep = strings.ReplaceAll(sep, "HH", "15")
	sep = strings.ReplaceAll(sep, "mm", "04")
	sep = strings.ReplaceAll(sep, "ss", "05")
	return nt.Format(sep)
} // 存在分隔符

func (t *TimeObj) StrToDate(str, sep string) (stamp time.Time, err error) {
	//获取当前时间
	stamp, err = time.ParseInLocation(sep, str, time.Local)
	if err != nil {
		return
	}
	return stamp, err
} // 字符串转换时间戳 单位：秒

func (t *TimeObj) StrToDateUnix(str, sep string) (i int64, err error) {
	//获取当前时间
	stamp, err := t.StrToDate(str, sep)
	if err != nil {
		return 0, err
	}
	i = stamp.Unix()
	return
}

func (t *TimeObj) DateToyMdHms() string {
	format := "20060102150405"
	return t.Sep(format)
}

func (t *TimeObj) NowTimeToM() int64 {
	n := t.DateToyMdHms()
	n = n[8:]
	s := n[4:]
	sum := StrIs0(s)
	s = n[2:4]
	sum += StrIs0(s) * 60
	s = n[:2]
	sum += StrIs0(s) * 3600
	return sum
} // 当前时间转换成秒

func StrIs0(str string) int64 {
	n, _ := strconv.ParseInt(str, 10, 64)
	return n
}

func (t *TimeObj) DateToyMdHm() string {
	format := "200601021504"
	return t.Sep(format)
}
func (t *TimeObj) DateToyMdH() string {
	format := "2006010215"
	return t.Sep(format)
}
func (t *TimeObj) DateToyMd() string {
	format := "20060102"
	return t.Sep(format)
}
func (t *TimeObj) DateToyM() string {
	format := "200601"
	return t.Sep(format)
}
func (t *TimeObj) DateToyMdSep() string {
	format := "2006-01-02"
	return t.Sep(format)
}
func (t *TimeObj) DateToyMdHmsSep() string {
	f := "2006-01-02 15:04:05"
	return t.Sep(f)
}
func (t *TimeObj) DateToyMdHmSep() string {
	f := "2006-01-02 15:04"
	return t.Sep(f)
}
func (t *TimeObj) DateToyMdHSep() string {
	f := "2006-01-02 15"
	return t.Sep(f)
}
func (t *TimeObj) DateToyMdHmsDate() string {
	format := "20060102150405"
	return t.Sep(format)
}

func (t *TimeObj) DateToyMdHmsSepTo(time time.Time) string {
	//获取当前时间
	format := "2006-01-02 15:04:05"
	return time.Format(format)
}

func (t *TimeObj) DateToyMdHmsTo(time time.Time) string {
	format := "20060102150405"
	//获取当前时间
	return time.Format(format)
}

func (t *TimeObj) TimeUnixToMdHms(i int64) string {
	tm := time.Unix(i, 0)
	return tm.Format("2006-01-02 15:04:05")
}

type Time time.Time

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	timeFormart := "2006-01-02 15:04:05"
	now, err := time.ParseInLocation(`"`+timeFormart+`"`, string(data), time.Local)
	*t = Time(now)
	return
}

func (t Time) MarshalJSON() ([]byte, error) {
	timeFormart := "2006-01-02 15:04:05"
	b := make([]byte, 0, len(timeFormart)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, timeFormart)
	b = append(b, '"')
	return b, nil
}

func (t *TimeObj) TimeGet() int64 {
	return time.Now().Unix()
} //获取当前时间戳 秒

func (t *TimeObj) Now() *time.Time {
	a := time.Now()
	return &a
} // 获取当前时间戳单位秒
func (t *TimeObj) NowUnix() int64 {
	return time.Now().Unix()
} // 获取当前时间戳单位秒 10位
func (t *TimeObj) NowUnixStr() string {
	return fmt.Sprintf(`%d`, time.Now().Unix())
} // 获取当前时间戳单位秒
func (t *TimeObj) NowUnixMs() int64 {
	return StrIs0(t.NowUnixMsStr())
} // 获取当前时间戳单位纳秒 13位
func (t *TimeObj) NowUnixMsStr() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())[:13]
} // 获取当前时间戳单位纳秒 16位
func (t *TimeObj) NowUnixMicroseconds() int64 {
	return StrIs0(t.NowUnixMicrosecondsStr())
} // 获取当前时间戳单位纳秒
func (t *TimeObj) NowUnixMicrosecondsStr() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())[:16]
} // 获取当前时间戳单位纳秒

func (t *TimeObj) NowSecond() int64 {
	return time.Now().Unix()
}
func (t *TimeObj) TimeGetAdd(s int64) *int {
	n := time.Now().Unix() + s
	return (*int)(unsafe.Pointer(&n))
}
func (t Time) String() string {
	timeFormart := "2006-01-02 15:04:05"
	return time.Time(t).Format(timeFormart)
}

func (t *TimeObj) UnixToDateStr(num int64) string {
	tm := time.Unix(num, 0)
	//获取当前时间
	return tm.Format("2006-01-02 15:04:05")
} // 时间戳(秒)转日期字符串

func (t *TimeObj) DateToTimestamp(time string) (i int64, err error) {
	ts := strings.Split(time, ":")
	if len(ts) != 3 {
		err = errors.New("格式错误，无法完成转换")
		return
	}
	hh := ts[0]
	if hh != "00" {
		h, _ := strconv.ParseInt(hh, 10, 64)
		h = h * 3600000
		i += h
	}
	mm := ts[1]
	if mm != "00" {
		m, _ := strconv.ParseInt(mm, 10, 64)
		m = m * 60000
		i += m
	}
	ss := ts[2]
	if ss != "00" {
		s, _ := strconv.ParseFloat(ss, 64)
		s = s * 1000
		i += int64(s)
	}
	return
} // 将时间（00:00:00）转换成时间戳（毫秒）

func (t *TimeObj) TimestampToDate(n int64) (i string) {
	//ts := strings.Split(time, ":")
	//if len(ts) != 3 {
	//	err = errors.New("格式错误，无法完成转换")
	//	return
	//}
	//hh := ts[0]
	//if hh != "00" {
	//	h, _ := strconv.ParseInt(hh, 10, 64)
	//	i += h * 60 * 60
	//}
	//mm := ts[1]
	//if mm != "00" {
	//	m, _ := strconv.ParseInt(hh, 10, 64)
	//	i += m * 60
	//}
	//i = i * 1000
	//ss := ts[2]
	//if ss != "00" {
	//	s, _ := strconv.ParseFloat(ss, 64)
	//	i += int64(s * 1000)
	//}

	// yMdHms
	//s := strconv.FormatInt(time, 10)
	//求小时
	h := "00"
	m := "00"
	s := "00"
	H := n / 3600000
	if H > 0 && H < 10 {
		h = "0" + strconv.FormatInt(H, 10)
	} else {
		h = strconv.FormatInt(H, 10)
	}
	n = n % 3600000
	M := n / 60000
	if M > 0 && M < 10 {
		m = "0" + strconv.FormatInt(M, 10)
	} else {
		m = strconv.FormatInt(M, 10)
	}
	n = n % 60000
	S := n / 1000
	if S > 0 && S < 10 {
		s = "0" + strconv.FormatInt(S, 10)
	} else {
		s = strconv.FormatInt(S, 10)
	}
	Z := n % 1000
	i = h + ":" + m + ":" + s
	if Z > 0 {
		z := strconv.FormatInt(Z, 10)
		i += "." + z
	}
	return
} // 时间戳转换成字符串
func (t *TimeObj) FormatStr(n string) (time.Time, error) {
	return time.ParseInLocation("20060102", n, time.Local)
}

func (t *TimeObj) TimeToZero() int64 {
	nowTime := time.Now()
	return nowTime.AddDate(0, 0, 1).Unix() - nowTime.Unix()
} //计算现在到凌晨还有多长时间

func (t *TimeObj) TimeEndOfMonth() int {
	nowTime := time.Now()
	return t.Days(nowTime) - nowTime.Day()
} //计算指定时间到月底还剩几天

func (t TimeObj) Days(nowTime time.Time) int {
	y := nowTime.Year()
	month := nowTime.Month()
	m := int(month)
	days := 30
	switch m {
	case 1:
	case 3:
	case 5:
	case 7:
	case 8:
	case 10:
	case 12:
		days = 31
		break
	case 2:
		days = feb(y) //2月份单独处理/判断是否为闰年
		break
	}
	return days
} //计算指定月份有多少天

func (t *TimeObj) Format() {

} //指定日期到过年还有几天

func feb(y int) int {
	//系统自带函数(返回指定的年份是否办闰年)
	//System.DateTime.IsLeapYear(y);
	//四年一闰，百年不闰;四百年再闰
	if y%400 == 0 || (y%4 == 0 && y%100 != 0) {
		return 29
	} else {
		return 28
	}
} //判断2月的天数
