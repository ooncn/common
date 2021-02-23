package jpgCaptcha

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/ooncn/common/util"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"math"
	"math/rand"
	"time"
)

// JpgCaptcha gif 验证码
type JpgCaptcha struct {
	frontColors  []color.Color    //图片前景
	bkgColors    []color.Color    //图片背景
	disturbLevel DisturbLevel     //干扰级别
	fonts        []*truetype.Font //字体
	size         image.Point      //图片大小
}

// 验证码字符类型
type StrType int

const (
	NUM   StrType = iota // 数字
	LOWER                // 小写字母
	UPPER                // 大写字母
	ALL                  // 全部
)

// DisturbLevel 干扰级别
type DisturbLevel int

const (
	NORMAL DisturbLevel = 4
	MEDIUM DisturbLevel = 8
	HIGH   DisturbLevel = 16
)

func New() *JpgCaptcha {
	c := &JpgCaptcha{
		disturbLevel: MEDIUM,
		size:         image.Point{X: 128, Y: 48},
	}
	c.frontColors = []color.Color{color.Black}
	c.bkgColors = []color.Color{color.White}
	_ = c.AddFont(util.GetCurrentDirectory() + "/COLONNA.TTF")
	return c
}

// AddFont 添加一个字体
func (c *JpgCaptcha) AddFont(path string) error {
	fontData, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	font, err := freetype.ParseFont(fontData)
	if err != nil {
		return err
	}
	if c.fonts == nil {
		c.fonts = []*truetype.Font{}
	}
	c.fonts = append(c.fonts, font)
	return nil
}

//AddFontFromBytes allows to load font from slice of bytes, for example, load the font packed by https://github.com/jteeuwen/go-bindata
func (c *JpgCaptcha) AddFontFromBytes(contents []byte) error {
	font, err := freetype.ParseFont(contents)
	if err != nil {
		return err
	}
	if c.fonts == nil {
		c.fonts = []*truetype.Font{}
	}
	c.fonts = append(c.fonts, font)
	return nil
}

// SetFont 设置字体 可以设置多个
func (c *JpgCaptcha) SetFont(paths ...string) error {
	for _, v := range paths {
		if err := c.AddFont(v); err != nil {
			return err
		}
	}
	return nil
}

func (c *JpgCaptcha) SetDisturbance(d DisturbLevel) {
	if d > 0 {
		c.disturbLevel = d
	}
}

func (c *JpgCaptcha) SetFrontColor(colors ...color.Color) {
	if len(colors) > 0 {
		c.frontColors = c.frontColors[:0]
		for _, v := range colors {
			c.frontColors = append(c.frontColors, v)
		}
	}
}

func (c *JpgCaptcha) SetBkgColor(colors ...color.Color) {
	if len(colors) > 0 {
		c.bkgColors = c.bkgColors[:0]
		for _, v := range colors {
			c.bkgColors = append(c.bkgColors, v)
		}
	}
}

func (c *JpgCaptcha) SetSize(w, h int) {
	if w < 48 {
		w = 48
	}
	if h < 20 {
		h = 20
	}
	c.size = image.Point{w, h}
}

func (c *JpgCaptcha) randFont() *truetype.Font {
	if len(c.fonts) == 0 {
		return nil
	}
	return c.fonts[rand.Intn(len(c.fonts))]
}

// 绘制背景
func (c *JpgCaptcha) drawBkg(img *Image) {
	ra := rand.New(rand.NewSource(time.Now().UnixNano()))
	//填充主背景色
	bgColorIndex := ra.Intn(len(c.bkgColors))
	bkg := image.NewUniform(c.bkgColors[bgColorIndex])
	img.FillBkg(bkg)
}

// 绘制噪点
func (c *JpgCaptcha) drawNoises(img *Image) {
	ra := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 待绘制图片的尺寸
	size := img.Bounds().Size()
	dlen := int(c.disturbLevel)
	// 绘制干扰斑点
	for i := 0; i < dlen; i++ {
		x := ra.Intn(size.X)
		y := ra.Intn(size.Y)
		r := ra.Intn(size.Y/20) + 1
		colorIndex := ra.Intn(len(c.frontColors))
		img.DrawCircle(x, y, r, i%4 != 0, c.frontColors[colorIndex])
	}

	// 绘制干扰线
	for i := 0; i < dlen; i++ {
		x := ra.Intn(size.X)
		y := ra.Intn(size.Y)
		o := int(math.Pow(-1, float64(i)))
		w := ra.Intn(size.Y) * o
		h := ra.Intn(size.Y/10) * o
		colorIndex := ra.Intn(len(c.frontColors))
		img.DrawLine(x, y, x+w, y+h, c.frontColors[colorIndex])
	}

}

// 绘制噪点
func (c *JpgCaptcha) drawNoisesArr(img *Image, dot, line [][]int, frontColor color.Color) {

	// 绘制干扰斑点
	for i := 0; i < len(dot); i++ {
		img.DrawCircle(dot[i][0], dot[i][1], dot[i][2], i%4 != 0, frontColor)
	}

	// 绘制干扰线
	for i := 0; i < len(line); i++ {
		img.DrawLine(line[i][0], line[i][1], line[i][2], line[i][3], frontColor)
	}

}

// 绘制噪点
func (c *JpgCaptcha) getNoises() (dot, line [][]int) {
	ra := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 待绘制图片的尺寸
	size := c.size
	dlen := int(c.disturbLevel)
	// 绘制干扰斑点
	for i := 0; i < dlen; i++ {
		x := ra.Intn(size.X)
		y := ra.Intn(size.Y)
		r := ra.Intn(size.Y/20) + 1
		/*colorIndex := ra.Intn(len(c.frontColors))
		img.DrawCircle(x, y, r, i%4 != 0, c.frontColors[colorIndex])*/
		arr := []int{x, y, r}
		dot = append(dot, arr)
	}

	// 绘制干扰线
	for i := 0; i < dlen; i++ {
		x := ra.Intn(size.X)
		y := ra.Intn(size.Y)
		o := int(math.Pow(-1, float64(i)))
		w := ra.Intn(size.Y) * o
		h := ra.Intn(size.Y/10) * o
		/*colorIndex := ra.Intn(len(c.frontColors))
		img.DrawLine(x, y, x+w, y+h, c.frontColors[colorIndex])
		*/
		arr := []int{x, y, x + w, y + h}
		line = append(line, arr)
	}
	return
}

// 绘制文字
func (c *JpgCaptcha) drawString(str string, dot, line [][]int, frontColor color.Color) (tmp *Image) {
	tmp = NewImage(c.size.X, c.size.Y)

	// 文字大小为图片高度的 0.6
	fsize := int(float64(c.size.Y) * 0.6)
	// 用于生成随机角度
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 文字之间的距离
	// 左右各留文字的1/4大小为内部边距
	padding := fsize / 4
	gap := (c.size.X - padding*2) / (len(str))

	// 逐个绘制文字到图片上
	for i, char := range str {
		// 创建单个文字图片
		// 以文字为尺寸创建正方形的图形
		img := NewImage(fsize, fsize)
		// str.FillBkg(image.NewUniform(color.Black))
		// 随机取一个前景色
		colorIndex := r.Intn(len(c.frontColors))
		//随机取一个字体
		font := c.randFont()
		img.DrawString(font, c.frontColors[colorIndex], string(char), float64(fsize))

		// 转换角度后的文字图形
		rs := img.Rotate(float64(r.Intn(40) - 20))
		// 计算文字位置
		s := rs.Bounds().Size()
		left := i*gap + padding
		top := (c.size.Y - s.Y) / 2
		// 绘制到图片上
		draw.Draw(tmp, image.Rect(left, top, left+s.X, top+s.Y), rs, image.ZP, draw.Over)
	} /*
		if c.size.Y >= 48 {
			// 高度大于48添加波纹 小于48波纹影响用户识别
			tmp.distortTo(float64(fsize)/10, 200.0)
		}*/
	c.drawNoisesArr(tmp, dot, line, frontColor)
	return
}

// Create 生成一个验证码图片
func (c *JpgCaptcha) RangCaptcha() (gifData image.Image, str string) {
	str = string(c.randStr(4, int(ALL)))
	gifData = c.createGif(str)
	return
}
func (c *JpgCaptcha) RangCaptchaNum(num int) (gifData image.Image, str string) {
	str = RandAaInt(num)
	gifData = c.createGif(str)
	return
}

func RandAaInt(l int) string {
	str := "3456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	b := []byte(str)
	var s []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		s = append(s, b[r.Intn(len(b))])
	}
	return string(s)
}

// Create 生成一个验证码图片
func (c *JpgCaptcha) Create(num int, t StrType) (gifData image.Image, str string) {
	if num <= 0 {
		num = 4
	}
	str = string(c.randStr(num, int(t)))
	gifData = c.createGif(str)
	return
}

func (c *JpgCaptcha) CreateCustom(str string) image.Image {
	if len(str) == 0 {
		str = "unkown"
	}
	return c.createGif(str)
}

func (c *JpgCaptcha) createGif(str string) image.Image {
	dot, line := c.getNoises()
	// 用于生成随机角度
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	frontColor := c.frontColors[r.Intn(len(c.frontColors))]
	tmp := c.drawString(str, dot, line, frontColor)

	img := NewImage(c.size.X, c.size.Y)
	bkg := image.NewUniform(color.White)
	img.FillBkg(bkg)

	draw.Draw(img, tmp.Bounds(), tmp, image.Point{}, draw.Over)
	//draw.Draw(img, img.Bounds(), img, image.Point{}, draw.Src) //添加图片
	return img
}

var fontKinds = [][]int{{10, 48}, {26, 97}, {26, 65}}

// 生成随机字符串
// size 个数 kind 模式
func (c *JpgCaptcha) randStr(size int, kind int) []byte {
	ikind, result := kind, make([]byte, size)
	isAll := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if isAll {
			ikind = rand.Intn(3)
		}
		scope, base := fontKinds[ikind][0], fontKinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}
