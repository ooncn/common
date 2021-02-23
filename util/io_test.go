package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAppendToFile(t *testing.T) {
	AppendToFile(GetCurrentDirectory()+"a.txt", TimeUtil.DateToyMdHms())
}
func TestRemoveDirAndFiles(t *testing.T) {
	dir := strings.TrimSuffix(filepath.Base(os.Args[0]), filepath.Ext(os.Args[0]))
	dir = filepath.Join(os.TempDir(), dir)
	dirs := filepath.Join(dir, `tmpdir`)
	err := os.MkdirAll(dirs, 0777)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	file := filepath.Join(dir, `tmpfile`)
	f, err := os.Create(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	f.Close()
	//di,err := os.Create(file)
	//file = filepath.Join(di)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	f.Close()

	err = RemoveDirAndFiles(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
