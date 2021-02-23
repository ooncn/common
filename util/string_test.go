package util

import (
	"bytes"
	"common/logs"
	"encoding/base64"
	"fmt"
	"github.com/golang/freetype"
	"html"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestString(t *testing.T) {
	//sText := `{"data":"","code":-6,"message":"DISTRICTCODE\u6216DISTRICTCCODE\u5b57\u6bb5\u5168\u56fd\u884c\u653f\u533a\u5212code\u5fc5\u987b\u4e3a\u53bb\u9664\u65e0\u6548\u96f6\u7684\u5076\u6570,\u53c2\u6570\u4e3a\uff1a410701000000410781000000","requestcode":"2019102114054798100555"}`

	//sText := "中文"
	//textQuoted := strconv.QuoteToASCII(sText)
	//textUnquoted := textQuoted[1 : len(textQuoted)-1]
	//fmt.Println(UnicodeToString(sText))
	//fmt.Println(IntMax(1, 2, 33, 11, 50, 0))

	ch := "Ad213s"
	ch = "aaaaa"
	ch = "01234561"
	ch = "AAAA"
	ch = "玩笑开大"
	//regular := `^[A-Za-z0-9]+$`    //英文和数字
	//regular = `^[A-Za-z0-9]{2,5}$` //英文和数字 限制大小2-5个字符
	//regular = `^[A-Za-z0-9]{5}$`   //英文和数字 限制5个字符
	regular := `^[\p{Han}]{1,5}$`

	reg := regexp.MustCompile(regular)
	if reg.MatchString(ch) {
		fmt.Println("true")
	} else {
		fmt.Println("false")
	}

	password := "wer.1327"
	l := len(password)

	if CheckRegexp(password, fmt.Sprintf(`^[A-Za-z0-9]{%d}$`, l)) {
		fmt.Println(password)
	} else {
		fmt.Println("错误")
	}
}

/**
替换img src
*/

var src = `<p>
    <img class="" src="https://upload-images.jianshu.io/upload_images/13150128-17411ce520a9f9e8.jpg?imageMogr2/auto-orient/strip|imageView2/2/w/1200"/>
</p>`

func TestSrc(t *testing.T) {
	str := html.EscapeString(src)
	//fmt.Println(str)

	imgs := strings.Split(str, "&lt;img ")
	var content string
	for _, v := range imgs {
		src := "src=&#34;"
		i := strings.Index(v, src)
		// 判断是否符合图片请求
		if i >= 0 {
			content += v[:i] + src
			i += len(src)
			//证明该字符串中存在图片文件
			jc := v[i:]
			j := strings.Index(jc, "&#34;")
			if j >= 0 {
				imgUrl := jc[:j]
				h := HttpOk{Url: imgUrl}
				err := h.QueryGet()
				if err != nil {
					fmt.Println(err)
					return
				}
				resp := h.ResponseObj
				hander := resp.Header
				fileType := hander["Content-Type"][0]
				if len(fileType) < 6 || strings.Index(fileType, "image") < 0 {
					fmt.Println("ERROR_HTTP_CONTENTTYPE")
					return
				}
				f := h.Response
				//img(f)
				imgBase64 := base64.StdEncoding.EncodeToString(f)
				imgBase64 = "data:" + fileType + ";base64," + imgBase64
				content += imgBase64 + jc[j:]
			} else {
				content += jc
			}
		} else {
			content += v + "&lt;img "
		}
	}
	fmt.Println(html.UnescapeString(content))
}

// TODO 图片加文字水印
func img(file []byte) {
	//需要加水印的图片
	//imgfile, _ := os.Open("u=488179422,3251067872&fm=200&gp=0.jpg")
	//defer imgfile.Close()
	jpgimg, _ := jpeg.Decode(io.Reader(bytes.NewReader(file)))

	img := image.NewNRGBA(jpgimg.Bounds())

	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			img.Set(x, y, jpgimg.At(x, y))
		}
	}
	//拷贝一个字体文件到运行目录
	fontBytes, err := ioutil.ReadFile(GetCurrentDirectory() + "/msyh.ttc")
	if err != nil {
		log.Println(err)
	}

	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
	}

	f := freetype.NewContext()
	f.SetDPI(72)
	f.SetFont(font)
	f.SetFontSize(12)
	f.SetClip(jpgimg.Bounds())
	f.SetDst(img)
	f.SetSrc(image.NewUniform(color.RGBA{R: 255, G: 0, B: 0, A: 255}))

	pt := freetype.Pt(img.Bounds().Dx()-200, img.Bounds().Dy()-12)
	_, err = f.DrawString("中文 string 255.43,232.12312 老纪", pt)

	//draw.Draw(img,jpgimg.Bounds(),jpgimg,image.ZP,draw.Over)

	//保存到新文件中
	newfile, _ := os.Create(TimeUtil.DateToyMdHms() + "aaa.jpg")
	defer newfile.Close()

	err = jpeg.Encode(newfile, img, &jpeg.Options{100})
	if err != nil {
		fmt.Println(err)
	}
}

func TestDg(t *testing.T) {
	uri := "http://www.baidu.com"
	uri = url.QueryEscape(uri) // 加密
	fmt.Println(uri)
	fmt.Println(url.QueryUnescape(uri)) // 解密
	/*
		str, err := util.ReaderFile("systemConfig.o")
		if err != nil {
			fmt.Println("数据库配置文件获取失败" + err.Error())
			return
		}
		fmt.Println(str)*/
	if strings.Contains("multipart/form-data; boundary=--------------------------362807300093195501603772", "multipart/form-data") {
		fmt.Println("true")
	} else {
		fmt.Println("false")
	}
}
func TestS(t *testing.T) {
	var a int64
	var b float64
	a = 63
	b = 1.54
	fmt.Println(fmt.Sprintf("%d", a))
	fmt.Println(fmt.Sprintf("%.2f", b))
}
