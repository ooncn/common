package util

import (
	"bytes"
	"fmt"
	"github.com/axgle/mahonia"
	"math/rand"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

/*
	go 语言GBK 转UTF8
	str := "乱码的字符串变量"
    str = ConvertToString(str, "gbk", "utf-8")
    fmt.Println(str)
*/
type StringUtil struct{}

/**
GBK转UTF8
*/
func GBKConvertUTF8(src, srcCode, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

/**
获取唯一ID
*/
func GetIdToDateAndStr() string {
	return GetId()
}
func GetId() string {
	t := TimeUtil.DateToyMdHms()[2:]
	return t + strings.ToUpper(RandAaInt(6))
}

func randLen(l *int) {
	a := reflect.ValueOf(l)
	a = a.Elem()
	if *l < 4 {
		a.SetInt(4)
	} else if *l > 105 {
		a.SetInt(105)
	}
} //随机数的长度

func RandInt(l int) string {
	var s string
	rand.Seed(time.Now().UnixNano())
	rand.Intn(9)
	for i := 0; i < l; i++ {
		a := strconv.Itoa(rand.Intn(9))
		s += a
	}
	return s
} // 获取数字 随机数
func GetIdInt(l int) string {
	rand.Seed(time.Now().UnixNano())
	return RandIdNum(1) + RandInt(l)
} // 获取数字 随机数

// 获取出字母随机数
func RandAa(l int) string {
	rand.Seed(time.Now().UnixNano())
	Str := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	b := []byte(Str)
	s := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		s = append(s, b[r.Intn(len(b))])
	}
	return string(s)
}
func RandIdNum(l int) string {
	rand.Seed(time.Now().UnixNano())
	Str := "123456789"
	b := []byte(Str)
	s := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		s = append(s, b[r.Intn(len(b))])
	}
	return string(s)
}

// 获取出字母+数字随机数
func RandAaInt(l int) string {
	rand.Seed(time.Now().UnixNano())
	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz"
	b := []byte(str)
	var s []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		s = append(s, b[r.Intn(len(b))])
	}
	return string(s)
}

/**
  for i := 0; i < 10; i++ {
      a := rand.Int()
      fmt.Println(a)
  }
  for i := 0; i < 10; i++ {
      a := rand.Intn(100)
      fmt.Println(a)
  }
  for i := 0; i < 10; i++ {
      a := rand.Float32()
      fmt.Println(a)
  }
*/

func CamelCase(s string) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == '_' {
		t = append(t, 'X')
		i++
	}
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && isASCIIUpper(s[i+1]) {
			continue
		}
		if isASCIIDigit(c) {
			t = append(t, c)
			continue
		}

		if isASCIIUpper(c) {
			c ^= ' '
		}
		t = append(t, c)

		for i+1 < len(s) && isASCIIUpper(s[i+1]) {
			i++
			t = append(t, '_')
			t = append(t, bytes.ToLower([]byte{s[i]})[0])
		}
	}
	return string(t)
}
func isASCIIUpper(c byte) bool {
	return 'A' <= c && c <= 'Z'
}
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

//email verify
func IsEmail(s string) bool {
	if s == "" {
		return false
	}
	//pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	regular := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`
	return CheckRegexp(s, regular)
}

//判断手机号
func IsPhone(s string) bool {
	if s == "" {
		return false
	}
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,7,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	return CheckRegexp(s, regular)
}

//判断中文
func IsCh(s string) bool {
	if s == "" {
		return false
	}
	return CheckRegexp(s, `[\p{Han}]+`)
}

//判断身份证号
func IsIdCard(s string) bool {
	if s == "" {
		return false
	}
	regular := "^(\\d{15,15}|\\d{16,16}|\\d{17,17}|\\d{18,18}|\\d{19,19}|(\\d{17,17}[x|X]))$"
	return CheckRegexp(s, regular)
}

//判断IpV4
func IsIpV4(s string) bool {
	if s == "" {
		return false
	}
	regular := "^(1\\d{2}|2[0-4]\\d|25[0-5]|[1-9]\\d|[1-9])(\\.(1\\d{2}|2[0-4]\\d|25[0-5]|[1-9]\\d|\\d)){3}$"
	return CheckRegexp(s, regular)
}
func IsIp4(s string) bool {
	if s == "" {
		return false
	}
	regular := "`^((2[0-4]\\d|25[0-5]|[01]?\\d\\d?)\\.){3}(2[0-4]\\d|25[0-5]|[01]?\\d\\d?)$`"
	return CheckRegexp(s, regular)
}

//判断是否是数字
func IsNum(s string) bool {
	if s == "" {
		return false
	}
	return CheckRegexp(s, "^(1\\d{2}|2[0-4]\\d|25[0-5]|[1-9]\\d|[1-9])(\\.(1\\d{2}|2[0-4]\\d|25[0-5]|[1-9]\\d|\\d)){3}$")
}
func IsLandline(s string) bool {
	if s == "" {
		return false
	}
	return CheckRegexp(s, "[0-9]{3,4}[-][0-9]{8}")
}

//判断是否是数字
func CheckRegexp(s, r string) bool {
	return regexp.MustCompile(r).MatchString(s)
}

// Mail _
type Mail struct {
	string
}

func MatchesPass(userPassword string) bool {
	if IsNoBlank(userPassword) {
		userPassword = strings.Replace(userPassword, " ", "", -1)
		// 去除换行符
		userPassword = strings.Replace(userPassword, "\n", "", -1)
		userPassword = strings.Replace(userPassword, "\t", "", -1)
		//以字母开头，长度在6-18之间，只能包含字符、数字和下划线。
		matched, err := regexp.MatchString(`^[a-zA-Z0-9]\w{5,17}$`, userPassword)
		fmt.Println(err)
		return matched
	} else {
		return false
	}
}

// 将字符串数组化
func StrToArr(s string) []string {
	return strings.Split(s, ",")
}

// 将数组引号组合

func ArrSep(a []string) string {
	if len(a) > 0 {
		for k, v := range a {
			a[k] = "'" + v + "'"
		}
	}
	return strings.Join(a, ",")
}

func UnicodeToString(sText string) string {
	sUnicode := strings.Split(sText, "\\u")
	var e string
	for k, v := range sUnicode {
		if k == 0 {
			e += v
			continue
		}
		if len(v) < 1 {
			continue
		} else if len(v) < 4 {
			e += v
		}
		var s string
		if len(v) > 4 {
			s = v[4:]
			v = v[0:4]
		}
		temp, err := strconv.ParseInt(v, 16, 32)
		if err != nil {
			panic(err)
		}
		e += fmt.Sprintf("%c", temp)
		e += s
	}
	return e
}

const MAX = int(^uint(0) >> 1)
const MIN = int(^MAX)

/**
计算最小值
*/
func IntMin(a ...int) int {
	m := MAX
	for _, i := range a {
		if i < m {
			m = i
		}
	}
	return m
}

/**
计算最大值
*/
func IntMax(a ...int) int {
	m := MIN
	for _, i := range a {
		if i > m {
			m = i
		}
	}
	return m
}

func PhoneEncode(phone string) string {
	if IsPhone(phone) {
		return phone[0:3] + "****" + phone[7:]
	}
	return phone
}

func URLAddQuery(u string, m map[string]string) (uri *url.URL) {
	uri, _ = url.Parse(u)
	if m == nil {
		return
	}
	v := url.Values{}
	for k, a := range m {
		v.Set(k, a)
	}
	q := uri.RawQuery
	if q == "" {
		uri.RawQuery = v.Encode()
	} else {
		uri.RawQuery = q + "&" + v.Encode()
	}
	return
}

func CheckPassword(password string, i int) bool {
	var c = 0
	l := len(password)
	if l > 6 || l < 36 {
		c += 1
	}
	return false
}

func MapInterfaceToString(m map[string]interface{}) map[string]string {
	var v = make(map[string]string)
	for k, val := range m {
		t := reflect.TypeOf(val)
		switch t.Kind() {
		case reflect.String:
			v[k] = val.(string)
			break
		case reflect.Bool:
			if val.(bool) {
				v[k] = "true"
			} else {
				v[k] = "false"
			}
			break
		case reflect.Int:
			v[k] = strconv.FormatInt(int64(val.(int)), 10)
			break
		case reflect.Float32:
			v[k] = strconv.FormatFloat(float64(val.(float32)), 'f', -1, 32)
			break
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v[k] = strconv.FormatInt(val.(int64), 10)
			break
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			v[k] = strconv.FormatUint(val.(uint64), 10)
			break
		case reflect.Float64:
			v[k] = strconv.FormatFloat(val.(float64), 'f', -1, 32)
			break
		}
	}
	return v
}
func PuStrLen(str string) (s, c int64) {
	for _, t := range str {
		if unicode.Is(unicode.Han, t) {
			c++
		} else {
			s++
		}
	}
	s += c
	return
}
