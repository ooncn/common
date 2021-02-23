package util

import (
	"os"
	"time"
)

// 根据日期获取文件名
func todayFilename(name string) string {
	today := time.Now().Format("2006-01-02")
	return name + "-" + today + ".log"
}

func NewLogFile(name string) (f *os.File, err error) {
	filename := GetCurrentDirectory() + "/log/" + time.Now().Format("2006-01") + "/"
	_ = os.MkdirAll(filename, os.ModePerm)
	filename += todayFilename(name)
	//打开文件，如果服务器重新启动，它将添加到今天的文件中。
	//f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	f, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	return
}
