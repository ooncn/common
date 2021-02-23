package util

import (
	"fmt"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	// 30天
	ss := "1000255:59:59.01"
	i, _ := TimeUtil.DateToTimestamp(ss)
	fmt.Println(ss)
	fmt.Println(i)
	s := TimeUtil.TimestampToDate(i)
	fmt.Println(s)

	fmt.Println(TimeUtil.NowSecond())

	timeStr := time.Now().Format("2006-01-02")
	fmt.Println(timeStr)

	//使用Parse 默认获取为UTC时区 需要获取本地时区 所以使用ParseInLocation
	t1, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr+" 23:59:59", time.Local)
	t2, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)

	nowTime := time.Now()
	fmt.Println(t1.Unix() + 1)
	fmt.Println(t2.AddDate(0, 0, 1).Unix())
	fmt.Println(nowTime.AddDate(0, 0, 1).Unix())
	fmt.Println(nowTime.AddDate(0, 0, 1).Format("2006-01-02"))
	fmt.Println(t2.AddDate(1, 0, 0).Unix())
	fmt.Println(nowTime.AddDate(1, 0, 0).Unix())
	fmt.Println(nowTime.AddDate(1, 0, 0).Format("2006-01-02"))
}
