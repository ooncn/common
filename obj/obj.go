package obj

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"github.com/ooncn/common/constant"
	"github.com/ooncn/common/oredis"
	"github.com/ooncn/common/util"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

type LoginForm struct {
	Account    string `json:"account"`    //账号
	Password   string `json:"password"`   //密码
	Imgkey     string `json:"imgKey"`     //用户类型
	ImgCode    string `json:"imgCode"`    //验证码
	RememberMe int    `json:"rememberMe"` //记住我
}

func (w *LoginForm) Verification() interface{} {
	// 判断防伪码生成规则
	m := make(map[string]string)
	ms := make([]string, 0)
	if !util.CheckRegexp(w.Account, `^[A-Za-z0-9]{4,36}$`) {
		m["account"] = "账号必须为英文或数字，长度为4-36字符"
		ms = append(ms, m["account"])
	}
	if !util.CheckRegexp(w.Imgkey, `^[A-Za-z0-9]{6,36}$`) && !util.CheckRegexp(w.ImgCode, `^[A-Za-z0-9]{4,10}$`) {
		m["imgCode"] = "验证码错误"
		ms = append(ms, m["imgCode"])
	}
	var img = ImgCode{Key: w.Imgkey}
	if !img.Verification(w.ImgCode) {
		m["imgCode"] = "验证码错误"
		ms = append(ms, m["imgCode"])
	}
	if !util.CheckRegexp(w.Password, `^[A-Za-z0-9]{6,36}$`) {
		m["password"] = "密码必须为英文或数字，长度为6-36字符"
		ms = append(ms, m["password"])
	}
	if w.RememberMe < 0 || w.RememberMe > 1 {
		m["rememberMe"] = "记住我错误"
		ms = append(ms, m["rememberMe"])
	}
	if len(ms) > 0 {
		//return [3]interface{}{m, ms, strings.Join(ms, "\n")}
		return strings.Join(ms, "\n")
	}
	return nil
}

type VoUserToken struct {
	Id       string   `json:"id"`       //用户主键ID
	Token    string   `json:"token"`    //密码
	Username string   `json:"username"` //账号
	Type     int      `json:"type"`     //用户类型
	Status   int      `json:"status"`   //状态
	Auth     []VOMenu `json:"auth"`     //登录权限
}

type VOMenu struct {
	Id          string `json:"id"`          //主键ID
	Sort        *int   `json:"sort"`        //排序
	Name        string `json:"name"`        //分类ID
	Description string `json:"description"` //商品ID
	Icon        string `json:"icon"`        //标识
	Pid         string `json:"pid"`         //上级
	Path        string `json:"path"`        //商品ID
	Status      *int   `json:"status"`      //商品ID
}
type SqlCount struct {
	Id       string `json:"id"`        //排序
	Count    int64  `json:"count"`     //排序
	Size     int64  `json:"size"`      //排序
	TimeLong int64  `json:"time_long"` //排序
}

type GroupTree struct {
	Id       string      `json:"id"`
	Code     string      `json:"code"`
	Name     string      `json:"name"`
	Sort     int         `json:"sort"`
	Children []GroupTree `json:"children"`
}

type IdCard struct {
	Id        string `json:"id"`        //id 身份证
	Photo     string `json:"photo"`     //photo 照片
	SnDN      string `json:"snDn"`      //身份证物理卡号
	IcCode    string `json:"icCode"`    //身份证磁卡号
	Name      string `json:"name"`      //name 姓名
	Gender    string `json:"gender"`    //gender 性别
	Nation    string `json:"nation"`    //nation 民族
	Address   string `json:"address"`   //address 地址
	BirthDate string `json:"birthDate"` //birth_date 出生日期
	IssuedBy  string `json:"issuedBy"`  //issued_by 签发机关
	StartTime string `json:"startTime"` //start_time 开始时间
	EndTime   string `json:"endTime"`   //end_time 到期时间
}

type ReqSync struct {
	Cid  string `json:"cid"`  // 企业编号
	Pid  string `json:"pid"`  // 项目编号
	Md5  string `json:"md5"`  // md5(cid + pid)
	Data string `json:"data"` //加密文件部分
	// +密文
}

type UpUrl struct {
	Src    string //原文
	Cipher string //密文
	Key    string //秘钥
}

func (u *UpUrl) Encrypt() string {
	if len(u.Key) != 16 {
		u.Key = Key()
	}
	u.Cipher = util.AesCBCEncrypt(u.Src, u.Key)
	return u.Cipher
}

func (u *UpUrl) Decrypt() string {
	if len(u.Key) != 16 {
		u.Key = Key()
	}
	u.Src = util.AesCBCDecrypt(u.Cipher, u.Key)
	return u.Src
}

func Key() string {
	//获取当前时间
	return time.Now().Format("2006-01-02_15:04")
}

type DBSetting struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
	DataBase string `json:"dataBase"`
	Url      string `json:"url"`
	Port     string `json:"port"`
	LogMode  bool   `json:"logMode"`
}
type RedisSetting struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}
type EmailSetting struct {
	UserName string `json:"userName"` //用户名
	Account  string `json:"account"`  //账号
	Pass     string `json:"pass"`     //密码
	Host     string `json:"host"`     //服务器地址
	Port     string `json:"Port"`     //端口号
}

type ServiceConfig struct {
	AuthTime     int64
	DBSetting    *DBSetting
	RedisSetting *RedisSetting
	Email        *EmailSetting
}

type ImgCode struct {
	Src string `json:"src"` //base64字符串
	Key string `json:"key"` //Key值
}

func (i *ImgCode) Verification(imgCode string) (b bool) {
	key := constant.RedisImgCode + i.Key
	code, _ := oredis.Ser.Get(key)
	b = code == strings.ToUpper(imgCode)
	_, _ = oredis.Ser.Del(key)
	return
}

// 系统设置配置文件
type SystemSetting struct {
	Logo interface{}
}

type Obj string

func (p Obj) Get() string {
	return string(p)
}
func (p Obj) GetInt64() int64 {
	i, _ := strconv.ParseInt(p.Get(), 10, 64)
	return i
}
func (p Obj) GetInt() int {
	return int(p.GetInt64())
}
func (p Obj) GetFloat64() float64 {
	i, _ := strconv.ParseFloat(p.Get(), 64)
	return i
}

type Map map[string]string

func (m Map) New() {
	m = make(map[string]string)
}

// map本来已经是引用类型了，所以不需要 *Params
func (p Map) SetString(k, s string) Map {
	p[k] = s
	return p
}
func (p Map) Set(k, s string) Map {
	p[k] = s
	return p
}

func (p Map) GetString(k string) string {
	s, _ := p[k]
	return s
}
func (p Map) Get(k string) string {
	s, _ := p[k]
	return s
}
func (p Map) SetInt64(k string, i int64) Map {
	p[k] = strconv.FormatInt(i, 10)
	return p
}

func (p Map) GetInt64(k string) int64 {
	i, _ := strconv.ParseInt(p.GetString(k), 10, 64)
	return i
}

func (p Map) GetInt(k string) int {
	return int(p.GetInt64(k))
}
func (p Map) SetFloat64(k string, i float64) Map {
	p[k] = strconv.FormatFloat(i, 'f', -1, 64)
	return p
}
func (p Map) GetFloat64(k string) float64 {
	i, _ := strconv.ParseFloat(p.GetString(k), 64)
	return i
}

// 判断key是否存在
func (p Map) ContainsKey(key string) bool {
	_, ok := p[key]
	return ok
}

// 签名
func (p *Map) SignGet(key, signType string) string {
	m := *p
	delete(m, "sign")
	str := MapAndStr(m)
	str += "&key=" + key
	switch signType {
	case "MD5":
		str = util.Md5Str(str) //需转换成切片
	case "HMACSHA256":
		dataSha256 := sha256.Sum256([]byte(str))
		str = hex.EncodeToString(dataSha256[:])
	}
	str = strings.ToUpper(str)
	return str
}

func (p *Map) SignValid(key, signType, sign string) bool {
	if !p.ContainsKey(sign) {
		return false
	}
	return p.Get(sign) == p.SignGet(key, signType)
}

func AllToStr(i interface{}) (s string) {
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
		s = strconv.FormatFloat(v.Float(), 'f', -1, 32)
	}
	return
}

func MapAndStr(m map[string]string) string {
	var data []string
	for k, v := range m {
		if v == "" {
			continue
		}
		data = append(data, fmt.Sprintf(`%s=%s`, k, v))
	}
	sort.Strings(data)
	return strings.Join(data, "&")
}

func XmlToMap(xmlStr string) Map {

	params := make(Map)
	decoder := xml.NewDecoder(strings.NewReader(xmlStr))

	var (
		key   string
		value string
	)

	for t, err := decoder.Token(); err == nil; t, err = decoder.Token() {
		switch token := t.(type) {
		case xml.StartElement: // 开始标签
			key = token.Name.Local
		case xml.CharData: // 标签内容
			content := string([]byte(token))
			value = content
		}
		if key != "xml" {
			if value != "\n" {
				params.SetString(key, value)
			}
		}
	}

	return params
}

func MapToXml(params Map) string {
	var buf bytes.Buffer
	buf.WriteString(`<xml>`)
	for k, v := range params {
		buf.WriteString(`<`)
		buf.WriteString(k)
		buf.WriteString(`><![CDATA[`)
		buf.WriteString(v)
		buf.WriteString(`]]></`)
		buf.WriteString(k)
		buf.WriteString(`>`)
	}
	buf.WriteString(`</xml>`)

	return buf.String()
}

type Resp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
type RespStr struct {
	Resp
	Data string `json:"data"`
}
type RespInterface struct {
	Resp
	Data interface{} `json:"data"`
}
type RespList struct {
	Resp
	Data []interface{} `json:"data"`
}
type RespMList struct {
	Resp
	Data MapList `json:"data"`
}

type Page map[string]interface{}
type MapList struct {
	Models interface{} `json:"models"`
	Pager  Page        `json:"pager"`
}
