package util

import (
	"fmt"
	"runtime"
	"testing"
)

// 单用户授权
func TestBase64(t *testing.T) {
	base64, err := FileToBase64(`D:\projects\id100\src\ODevice\update\A1\2019\10\20191011100357488272.jpeg`)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(base64)
}

func TestOS(t *testing.T) {
	fmt.Println(runtime.GOOS)
}
