package Handler

import (
	"bytes"
	"encoding/base64"
	"errors"
	"github.com/kataras/iris/v12"
	"github.com/ooncn/common/constant"
	"github.com/ooncn/common/gifCaptcha"
	"github.com/ooncn/common/obj"
	"github.com/ooncn/common/oredis"
	"github.com/ooncn/common/util"
	"gorm.io/gorm"
	"image/color"
	"image/gif"
	"net/url"
	"reflect"
	"strings"
	"time"
)

type Object interface{}

type Context struct {
	iris.Context
	Ip           string     // 客户端请求IP地址
	Form         url.Values // 客户端请求IP地址
	Resp         interface{}
	Err          interface{}
	PageSizeKey  string
	PageNumKey   string
	PageCountKey string
	PageNum      int
	PageSize     int
	Page         map[string]interface{}
	UserType     int    //用户类型
	UserId       string //用户类型
}

//获取用户token
func (c *Context) SetCxt(cxt iris.Context) {
	c.Context = cxt
}

func (c *Context) PageGetCxt(cxt iris.Context) {
	c.Context = cxt
	c.PageGet()
}
func (c *Context) PageGet() {
	if c.Page == nil {
		c.Page = make(map[string]interface{})
	}
	c.PageNumKey = "page_num"
	c.PageCountKey = "page_count"
	pageNum, err := c.URLParamInt("page_num")
	if err != nil {
		pageNum, err := c.URLParamInt("pageNum")
		c.PageCountKey = "pageCount"
		c.PageNumKey = "pageNum"
		if err != nil || pageNum <= 0 {
			pageNum = 1
		}
	}
	c.PageSizeKey = "page_size"
	pageSize, err := c.URLParamInt("page_size")
	if err != nil {
		pageSize, err = c.URLParamInt("pageSize")
		c.PageSizeKey = "pageSize"
		if err != nil {
			pageSize = 10
		}
	}
	c.PageNum = pageNum
	c.PageSize = pageSize
	page := make(map[string]interface{})
	page[c.PageCountKey] = 0
	page[c.PageNumKey] = pageNum
	page[c.PageSizeKey] = pageSize
	c.Page = page
}
func (c *Context) PageSetCount(count int64) map[string]interface{} {
	c.Page[c.PageCountKey] = count
	return c.Page
}
func (c *Context) PageMapListSetCount(count int64) *obj.MapList {
	m := new(obj.MapList)
	m.Pager = c.PageSetCount(count)
	return m
}
func (c *Context) Offset(count int64) *obj.MapList {
	m := new(obj.MapList)
	m.Pager = c.PageSetCount(count)
	return m
}
func (c *Context) PageMapList() *obj.MapList {
	m := new(obj.MapList)
	m.Pager = c.Page
	return m
}

func (c *Context) Token() (token string) {
	token = c.GetHeader("token")
	if len(token) != 32 {
		token = c.URLParam("token")
		if len(token) != 32 {
			token = c.GetHeader("m_token")
			if len(token) != 32 {
				token = c.URLParam("m_token")
				if len(token) != 32 {
					token = c.GetHeader("mtoken")
					if len(token) != 32 {
						token = c.URLParam("mtoken")
					}
				}
			}
		}
	}
	return
}

func (c *Context) GetUser() (user obj.VoUserToken, err error) {
	err = oredis.Ser.GetType(constant.REDIS_USER_SESSION+c.Token(), &user)
	if err != nil {
		c.RespErrorD("LOGOUT")
	}
	return
}
func (c *Context) GetAuthMenu() ([]map[string]interface{}, error) {
	a, err := c.GetUser()
	if err != nil {
		return nil, err
	}
	list := a.Auth
	if len(list) < 1 {
		return nil, nil
	}
	m := make([]map[string]interface{}, len(list))
	for k, v := range list {
		a := make(map[string]interface{})
		a["id"] = v.Id
		a["path"] = v.Path
		a["name"] = v.Name
		a["icon"] = v.Icon
		a["status"] = v.Status
		a["pid"] = v.Pid
		m[k] = a
	}
	return m, nil
}
func (c *Context) GetAuthList() ([]string, error) {
	a, err := c.GetUser()
	if err != nil {
		return nil, err
	}
	list := a.Auth
	if len(list) < 1 {
		return nil, nil
	}
	authList := make([]string, len(list))
	for k, v := range list {
		authList[k] = v.Path
	}
	return authList, nil
}

func (c *Context) GetUserId() string {
	token := c.Token()
	id, err := oredis.Ser.Get(constant.REDIS_USER_ID + token)
	if err != nil {
		id, err = oredis.Ser.Get(constant.REDIS_MANAGER_USER_ID + token)
	}
	return id
}
func (c *Context) GetAdminId() (string, error) {
	token := c.Token()
	id, err := oredis.Ser.Get(constant.REDIS_MANAGER_USER_ID + token)
	if err != nil {
		c.RespErrorD("LOGOUT")
		return "", errors.New("LOGOUT")
	}
	return id, nil
}
func (c *Context) GetIP() string {
	ip := c.RemoteAddr()
	switch ip {
	case "127.0.0.1", "localhost":
		ip = c.GetHeader("X-Forwarded-For")
		break
	}
	if ip == "" {
		ip = c.GetHeader("X-Forwarded-For")
		if ip == "" {
			ip = c.GetHeader("X-real-ip")
		}
	}
	return ip
}
func (c *Context) GetIPOne() string {
	ip := c.GetIP()
	return strings.Split(ip, ",")[0]
}
func (c *Context) GetUserAgent() string {
	return c.GetHeader("User-Agent")
}

func (c *Context) Body() []byte {
	req, _ := c.GetBody()
	return req
}

func (c *Context) BodyStr() string {
	return string(c.Body())
}

func (c *Context) GetToken() string {
	_, _, token := util.UserAgentAndIpToToken(c.Request())
	return token
}

func (c *Context) RespSuccess() {
	c.RespSuccessCM(200, "success")
}
func (c *Context) RespSuccessD(data interface{}) {
	c.RespSuccessMD("success", data)
}
func (c *Context) RespSuccessMD(msg string, data interface{}) {
	c.RespSuccessCMD(200, msg, data)
}
func (c *Context) RespSuccessCM(code interface{}, msg string) {
	c.RespSuccessCMD(code, msg, nil)
}

func (c *Context) RespSuccessCMD(code interface{}, msg string, data interface{}) {
	r := make(map[string]interface{})
	r["code"] = code
	r["msg"] = msg
	if !util.IsBlank(data) {
		r["data"] = &data
	}
	if c.Err == nil {
		if c.Resp == nil {
			c.Resp = r
			c.JSON(c.Resp)
		}
	}
	c.StopExecution()
}

func (c *Context) RespError() {
	c.RespErrorCMD(500, "error", nil)
}
func (c *Context) RespErrorD(data interface{}) {
	c.RespErrorCMD(500, "error", data)
}
func (c *Context) RespErrorM(msg string) {
	c.RespErrorCMD(500, msg, nil)
}
func (c *Context) RespErrorMD(msg string, data interface{}) {
	c.RespErrorCMD(500, msg, data)
}
func (c *Context) RespErrorCM(code interface{}, msg string) {
	c.RespErrorCMD(code, msg, nil)
}
func (c *Context) RespErrorCMD(code interface{}, msg string, data interface{}) {
	if c.IsStopped() {
		return
	}
	r := make(map[string]interface{})
	r["code"] = code
	r["msg"] = msg
	if data != nil {
		t := reflect.ValueOf(data)
		s := t.Type().String()
		if s == "*errors.fundamental" {
			t = reflect.Indirect(t)
			f0 := t.FieldByName("msg") //获取结构体s中第一个元素a
			r["data"] = f0.String()
		} else if s == "*pq.Error" {
			t = reflect.Indirect(t)
			f0 := t.FieldByName("msg") //获取结构体s中第一个元素a
			r["data"] = util.StrError(errors.New(f0.String()))
		} else if !util.IsBlank(data) {
			r["data"] = data
		}
	}
	if c.Err == nil {
		c.Err = r
	}
	c.JSON(c.Err)
	c.StopExecution()
}

func (c *Context) BodyMap() map[string]interface{} {
	ty := new(map[string]interface{})
	util.JsonToType(c.BodyStr(), &ty)
	return *ty
}

func (c *Context) ImgCode() (img obj.ImgCode, err error) {
	var captcha = gifCaptcha.New()
	captcha.SetFont(util.GetCurrentDirectory() + "/font/COLONNA.TTF")
	captcha.SetFrontColor(color.Black,
		color.Black,
		color.Black,
		color.Black)
	_, _, token := util.UserAgentAndIpToToken(c.Request())
	gifData, code := captcha.RangCaptchaNum(4)
	buffer := new(bytes.Buffer)
	gif.EncodeAll(buffer, gifData)
	encodeString := base64.StdEncoding.EncodeToString(buffer.Bytes())
	img.Key = token
	img.Src = "data:image/gif;base64," + encodeString
	err = oredis.Ser.Set(constant.RedisImgCode+token, strings.ToUpper(code), 120*time.Second)
	if err != nil {
		c.RespErrorD("LOGOUT")
		c.StopExecution()
	}
	return
} // 图片验证码
func (c *Context) IsFormData() bool {
	return strings.Contains(c.GetHeader("Content-Type"), "multipart/form-data")
}
func (c *Context) IsFormUrlencoded() bool {
	return strings.Contains(c.GetHeader("Content-Type"), "application/x-www-form-urlencoded")
}
func (c *Context) IsForm() bool {
	header := c.GetHeader("Content-Type")
	return strings.Contains(header, "multipart/form-data") ||
		strings.Contains(header, "multipart/form-data")
}
func (c *Context) Logout() {
	token := c.Token()
	_, _ = oredis.Ser.Del(constant.REDIS_USER_SESSION + token)
	_, _ = oredis.Ser.Del(constant.REDIS_USER_TOKEN + token)
	_, _ = oredis.Ser.Del(constant.REDIS_USER_ID + token)
	_, _ = oredis.Ser.Del(constant.REDIS_USER_TYPE + token)
}
func (c *Context) LogoutManager() {
	token := c.Token()
	_, _ = oredis.Ser.Del(constant.REDIS_MANAGER_USER_SESSION + token)
	_, _ = oredis.Ser.Del(constant.REDIS_MANAGER_USER_TOKEN + token)
	_, _ = oredis.Ser.Del(constant.REDIS_USER_ID + token)
	_, _ = oredis.Ser.Del(constant.REDIS_MANAGER_USER_TYPE + token)
}
func (c *Context) GetUserType() int {
	token := c.Token()
	var i int
	oredis.Ser.GetType(constant.REDIS_USER_TYPE+token, &i)
	return i
}

// 短信验证码
// 修改密码
// 邮箱激活码
func (c *Context) Login(form obj.LoginForm, sqlSer *gorm.DB, getInfo func(sqlSer *gorm.DB) (*obj.UserInfoVo, error)) (*obj.UserInfoVo, error) {
	vo, err := getInfo(sqlSer)
	if err != nil {
		return nil, err
	}
	_, _, token := util.UserAgentAndIpToToken(c.Request())
	vo.Token = token
	uid := vo.Id
	UserLoginRedis(form.RememberMe, constant.REDIS_USER_TOKEN+uid, constant.REDIS_USER_ID+token, constant.REDIS_USER_SESSION, uid, token, vo)
	return vo, nil
}
func UserLoginRedis(rememberMe int, redisUserToken, RedisUserId, redisUserSession, uid, token string, obj interface{}) {
	t := 2 * time.Hour
	if rememberMe == 1 {
		t = 0
	}
	lodToken, err := oredis.Ser.Get(redisUserToken)
	if err == nil {
		_, err = oredis.Ser.Del(redisUserSession + lodToken)
	}
	err = oredis.Ser.Set(redisUserToken, token, t)
	err = oredis.Ser.Set(RedisUserId, uid, t)
	err = oredis.Ser.SetType(redisUserSession+token, obj, t)
}
