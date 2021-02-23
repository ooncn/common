package util

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/Comdex/imgo"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	//"github.com/axgle/mahonia"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

/**
截取字符串
*/
func Substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}
func IsPictureFormat(path string) (ext string, filename string, err error) {
	temp := strings.Split(path, ".")
	if len(temp) <= 1 {
		err = errors.New("文件路径不符合")
		return
	}
	mapRule := make(map[string]int64)
	mapRule["jpg"] = 1
	mapRule["png"] = 1
	mapRule["jpeg"] = 1
	/** 添加其他格式 */
	if mapRule[temp[1]] == 1 {
		ext = temp[1]
		filename = temp[0]
		return
	}
	err = errors.New("文件路径不符合")
	return
}

/**
获取当前父级路径
*/
func GetParentDirectory(dirctory string) string {
	return Substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

//region 获取当前执行路径
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

//endregion

//region 图片与Base64字符串互转

//文件转base64
func FileToBase64(src string) (base string, err error) {
	return imgo.Img2Base64(src)
}

//Img2Base64 produce a base64 string from a image file.
func Img2Base64(file io.Reader) (encodeString string, err error) {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}
	encodeString = base64.StdEncoding.EncodeToString(data)
	return
}

//base64实例
func base64Test() {
	//读原图片
	ff, _ := os.Open("b.png")
	defer ff.Close()
	sourceBuffer := make([]byte, 500000)
	n, _ := ff.Read(sourceBuffer)
	//base64压缩
	sourceString := base64.StdEncoding.EncodeToString(sourceBuffer[:n])

	//写入临时文件
	ioutil.WriteFile("a.png.txt", []byte(sourceString), 0667)
	//读取临时文件
	cc, _ := ioutil.ReadFile("a.png.txt")

	//解压
	dist, _ := base64.StdEncoding.DecodeString(string(cc))
	//写入新文件
	f, _ := os.OpenFile("xx.png", os.O_RDWR|os.O_CREATE, os.ModePerm)
	defer f.Close()
	f.Write(dist)

}

//base64转文件
func Base64ToFileByte(base64Str string) (toMd5 string, dist []byte, e error) {
	dist, e = base64.StdEncoding.DecodeString(base64Str)
	if e != nil {
		return
	}
	toMd5 = fileMd5(dist)
	return
}
func FileSaveByte(dist []byte, src string) (size int, e error) {
	//写入新文件
	f, e := os.OpenFile(src, os.O_RDWR|os.O_CREATE, os.ModePerm)
	defer f.Close()
	size, e = f.Write(dist)
	return
}

func Base64ToFile(base64Str, src string) (toMd5 string, size int, e error) {
	dist, e := base64.StdEncoding.DecodeString(base64Str)
	toMd5 = fileMd5(dist)
	//写入新文件
	f, e := os.OpenFile(src, os.O_RDWR|os.O_CREATE, os.ModePerm)
	defer f.Close()
	size, e = f.Write(dist)
	return
}

func Base64ToFileMd5(base64Str, src, md5Str string) (toMd5 string, size int, e error) { //读原图片
	//解压
	dist, _ := base64.StdEncoding.DecodeString(base64Str)
	toMd5 = fileMd5(dist)
	if toMd5 != strings.ToUpper(md5Str) {
		e = errors.New("FILE_MD5_ERROR") //文件md5值错误
		return
	}
	//写入新文件
	f, _ := os.OpenFile(src, os.O_RDWR|os.O_CREATE, os.ModePerm)
	defer f.Close()
	size, e = f.Write(dist)
	return
}
func fileMd5(dist []byte) (toMd5 string) {
	h := md5.New()
	h.Write(dist)
	cipherStr := h.Sum(nil)
	toMd5 = strings.ToUpper(hex.EncodeToString(cipherStr))
	return
}

/**
base64获取文件类型和后缀名
fileType string 	文件类型
ext 	string 		文件后缀
eStr 	string		base64密文
err 	error		错误
*/
func Base64ToFileTypeAndExt(file string) (fileType string, ext string, eStr string, err error) {
	if file == "" {
		err = errors.New("file is null")
		return
	}
	str := file[:strings.Index(file, ",")]
	fileType = str[strings.Index(str, ":")+1 : strings.Index(str, "/")]
	ext = str[strings.Index(str, "/")+1 : strings.Index(str, ";base64")]
	eStr = file[strings.Index(file, ",")+1:]
	return
}

//base64转换文件保存到指定的位置
func Base64ToFileSetPath(file string, path string) (toMd5 string, fileType string, ext string, filePath string, size int, err error) {
	fileType, ext, Estr, err := Base64ToFileTypeAndExt(file)
	if err != nil {
		return
	}
	filePath = path
	_ = os.MkdirAll(filePath, os.ModePerm)
	filePath += GetIdToDateAndStr() + "." + ext
	toMd5, size, err = Base64ToFile(Estr, filePath)
	return
}

func Base64ToFileMd5SetPath(file, path, md5Str string) (toMd5 string, fileType string, ext string, filePath string, size int, err error) {
	fileType, ext, Estr, err := Base64ToFileTypeAndExt(file)
	if err != nil {
		return
	}
	filePath = path
	_ = os.MkdirAll(filePath, os.ModePerm)
	filePath += GetIdToDateAndStr() + "." + ext
	toMd5, size, err = Base64ToFileMd5(Estr, filePath, md5Str)
	return toMd5, fileType, ext, filePath, size, err
}

//endregion

//region 判断文件是否存在  存在返回 true 不存在返回false
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

//endregion

func FileShow(filename string) (*os.File, error) {
	if CheckFileIsExist(filename) { //如果文件存在
		f, err1 := os.OpenFile(filename, os.O_APPEND, 0666) //打开文件
		defer f.Close()
		fmt.Println("文件存在")
		return f, err1
	}
	return nil, nil
}

// 读取文件到[]byte中
func File2Bytes(filename string) ([]byte, error) {
	// File
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// FileInfo:
	stats, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// []byte
	data := make([]byte, stats.Size())
	count, err := file.Read(data)
	if err != nil {
		return nil, err
	}
	fmt.Printf("read file %s len: %d \n", filename, count)
	return data, nil
}

//

//region 获取远程图片
func GetUrlImg(fileDir, url string) (p string, n int64, err error) {
	path := strings.Split(url, "/")
	var name string
	if len(path) > 1 {
		name = path[len(path)-1]
	}
	p = GetCurrentDirectory() + fileDir + "/"
	_ = os.MkdirAll(p, os.ModePerm)
	p += name
	//History
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer f.Close()
	resp, err := GETUrlToByte(url)
	if err != nil {
		return
	}
	n, err = io.Copy(f, bytes.NewReader(resp))
	return
}

//endregion

func ReaderFile(path string) (str string, err error) {
	path = GetCurrentDirectory() + "/" + path
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	str = string(b)
	str = PuDecrypt(str)
	return
}
func ReadFile(path string) (str string, err error) {
	path = GetCurrentDirectory() + "/" + path
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	str = string(b)
	return
}

func ReaderFileByte(path string) (b []byte, err error) {
	path = GetCurrentDirectory() + "/" + path
	b, err = ioutil.ReadFile(path)
	return
}

// 保护配置文件
func SaveConfigFile(fileName, content string) (err error) {
	content = PuEncrypt(content)
	path := GetCurrentDirectory() + "/" + fileName
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()
	if err != nil {
		return
	} else {
		_, err = f.Write([]byte(content))
	}
	return
}
func SaveFile(fileName string, content []byte) (err error) {
	path := GetCurrentDirectory() + "/" + fileName
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()
	if err != nil {
		return
	} else {
		_, err = f.Write(content)
	}
	return
}
func SaveConfigFileByte(fileName string, content []byte) (err error) {
	path := GetCurrentDirectory() + "/" + fileName
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()
	if err != nil {
		return
	} else {
		_, err = f.Write(content)
	}
	return
}

func LogInfo(i interface{}) {
	var s string
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	switch t.Kind() {
	case reflect.String:
		s = v.String()
	case reflect.Bool:
		if v.Bool() {
			s = "true"
		} else {
			s = "false"
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s = strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		s = strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		s = strconv.FormatFloat(v.Float(), 'f', 5, 32)
	case reflect.Ptr:
		s = JsonToStr(i)
	}
	s = TimeUtil.DateToyMdHmsSep() + "\t" + s
	LogReader("log/", "log_"+TimeUtil.DateToyMdSep()+".log", s+"\n")
	fmt.Println(s)
}

func LogReader(path, fileName, c string) {
	path = GetCurrentDirectory() + "/" + path
	_ = os.MkdirAll(path, os.ModePerm)
	path = path + fileName
	b, err := ioutil.ReadFile(path)
	if err != nil {
		LogSaveFile(path, c)
	} else {
		LogSaveFile(path, string(b)+c)
	}
	return
}

// 保护配置文件
func LogSaveFile(path, content string) {
	//var enc mahonia.Decoder
	//enc = mahonia.NewDecoder("gbk")
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()
	if err != nil {
		return
	} else {
		_, err = f.WriteString(content)
	}
	return
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
} //判断目录或文件是否存在 存在true.不存在false 如果返回的错误为其它类型,则不确定是否在存在
func MkDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
} //创建多级目录

func AppendToFile(fileName string, content string) error {
	return AppendToFileByte(fileName, []byte(content))
}
func AppendToFileByte(fileName string, content []byte) error {
	// 以只写的模式，打开文件
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()
	if err == nil {
		// 查找文件末尾的偏移量
		n, _ := f.Seek(0, io.SeekEnd)
		// 从末尾的偏移量开始写入内容
		_, err = f.WriteAt(content, n)
	}
	return err
}

func RemoveDirAndFiles(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
} //删除目录以及下级所有文件

/*func main() {
}*/

func ImgCompress(
	getReadSizeFile func() (io.Reader, error),
	getDecodeFile func() (*os.File, error),
	to string,
	Quality,
	base int,
	format string) (err error) {
	/** 读取文件 */
	fileOrigin, err := getDecodeFile()
	defer fileOrigin.Close()
	if err != nil {
		fmt.Println("os.Open(file)错误")
		return
	}
	var origin image.Image
	var config image.Config
	var temp io.Reader
	/** 读取尺寸 */
	temp, err = getReadSizeFile()
	if err != nil {
		fmt.Println("os.Open(temp)")
		return
	}
	var typeImage int64
	format = strings.ToLower(format)
	/** jpg 格式 */
	if format == "jpg" || format == ".jpg" || format == ".jpeg" || format == "jpeg" {
		typeImage = 1
		origin, err = jpeg.Decode(fileOrigin)
		if err != nil {
			fmt.Println("jpeg.Decode(fileOrigin)")
			return
		}
		config, err = jpeg.DecodeConfig(temp)
		if err != nil {
			fmt.Println("jpeg.DecodeConfig(temp)")
			return
		}
	} else if format == "png" || format == ".png" {
		typeImage = 0
		origin, err = png.Decode(fileOrigin)
		if err != nil {
			fmt.Println("png.Decode(fileOrigin)")
			return
		}
		config, err = png.DecodeConfig(temp)
		if err != nil {
			fmt.Println("png.DecodeConfig(temp)")
			return
		}
	} else {
		err = errors.New("格式错误")
		return
	}
	/** 做等比缩放 */
	var width = uint(config.Width)
	var height = uint(config.Height)
	if config.Height > config.Width && config.Height > base {
		/** 做等比缩放 */
		height = uint(base) /** 基准 */
		width = uint(base * config.Width / config.Height)
	} else if config.Width > base {
		/** 做等比缩放 */
		width = uint(base) /** 基准 */
		height = uint(base * config.Height / config.Width)
	}

	canvas := resize.Thumbnail(width, height, origin, resize.Lanczos3)
	fileOut, err := os.Create(to)
	defer fileOut.Close()
	if err != nil {
		return
	}
	if typeImage == 0 {
		err = png.Encode(fileOut, canvas)
		if err != nil {
			fmt.Println("压缩图片失败")
			return
		}
	} else {
		err = jpeg.Encode(fileOut, canvas, &jpeg.Options{Quality})
		if err != nil {
			fmt.Println("压缩图片失败")
			return
		}
	}

	return
}
