package wechat

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/ooncn/common/constant"
	"github.com/ooncn/common/oredis"
	"github.com/ooncn/common/util"
	"io"
	"net/url"
	"sort"
	"strings"
	"time"
)

func MakeSignature(token, timestamp, nonce string) string {
	sl := []string{token, timestamp, nonce}
	sort.Strings(sl)
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))
	return fmt.Sprintf("%x", s.Sum(nil))
}

type WxToken struct {
	MqAppID         string `json:"mq_app_id" gorm:"primary_key;size:33"`      //MqAppID
	MqAppSecret     string `json:"mq_app_secret" gorm:"unique_index;size:33"` //MqAppSecret
	URL             string `json:"URL" gorm:"size:300"`                       //必须以http://或https://开头，分别支持80端口和443端口。
	Host            string `json:"host" gorm:"size:300"`                      //必须以http://或https://开头，分别支持80端口和443端口。
	EncodingAESKey  string `json:"encoding_aes_key" gorm:"size:43"`           //消息加密密钥由43位字符组成，可随机修改，字符范围为A-Z，a-z，0-9。
	EncodingAESType int    `json:"encoding_aes_type" gorm:"DEFAULT:0"`        // 加密方式：0、明文 1、兼容模式 2、安全模式
	Token           string `json:"token" gorm:"size:43"`                      //必须为英文或数字，长度为3-32字符。

	OpenAppid     string `json:"open_appid"`      //微信开发平台
	OpenAppSecret string `json:"open_app_secret"` //微信开发平台秘钥

	CreateTime *time.Time `json:"create_time"` // 创建时间
}
type WxPayNotify struct {
	TransactionId string     `json:"transaction_id" gorm:"primary_key;size:33"` //微信支付订单号
	OutTradeNo    string     `json:"out_trade_no" gorm:"size:33"`               //商户订单号
	Appid         string     `json:"appid" gorm:"size:68"`                      //交易类型
	MchId         string     `json:"mch_id" gorm:"size:20"`                     //交易类型
	Openid        string     `json:"openid" gorm:"size:68"`                     //交易类型
	TradeType     string     `json:"trade_type" gorm:"size:10"`                 //交易类型
	ResultCode    string     `json:"result_code" gorm:"size:10"`                //交易类型
	ReturnCode    string     `json:"return_code" gorm:"size:10"`                //交易类型
	ErrCode       string     `json:"err_code" gorm:"size:68"`                   //交易类型
	ErrCodeDes    string     `json:"err_code_des" gorm:"size:300"`              //错误返回的信息描述
	TotalFee      string     `json:"total_fee"  gorm:"type:int"`                //错误返回的信息描述
	Body          string     `json:"body" `
	CreateTime    *time.Time `json:"createTime" gorm:"type:timestamptz(0);DEFAULT:now()"` // 创建时间
}
type WxMqMsgData struct {
	Id           string     `json:"id" gorm:"primary_key;size:33"` //主键ID
	MsgId        string     `gorm:"size:33"`                       //消息id，64位整型
	MsgType      string     `gorm:"size:15"`                       //消息类型，文本为 text
	ToUserName   string     `gorm:"size:33"`                       //开发者微信号
	FromUserName string     `gorm:"size:33"`                       //发送方帐号（一个OpenID）
	Data         string     //消息主体（一个OpenID）
	CreateTime   int64      // 创建时间
	CreateDate   *time.Time `gorm:"type:timestamptz(0);DEFAULT:now()"` // 创建时间
}
type WxMqEventData struct {
	Id           string `json:"id" gorm:"primary_key;size:33"` //主键ID
	Event        string `gorm:"size:33"`                       //事件类型
	MsgType      string `gorm:"size:15"`                       //消息类型，文本为 text
	ToUserName   string `gorm:"size:33"`                       //开发者微信号
	FromUserName string `gorm:"size:33"`                       //发送方帐号（一个OpenID）
	Data         string `json:"data"`                          //消息主体（一个OpenID）
	CreateTime   int64  // 创建时间
}
type WxUserInfo struct {
	Openid         string     `json:"openid" gorm:"size:33;primary_key"`        //用户的唯一标识
	Unionid        string     `json:"unionid" gorm:"size:33;index"`             //只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段。
	Uid            string     `json:"uid" gorm:"index;size:33"`                 //主键ID
	Nickname       string     `json:"nickname" gorm:"size:10"`                  //用户昵称
	Sex            int        `json:"sex" gorm:"type:smallint;DEFAULT:0"`       //用户的性别，值为1时是男性，值为2时是女性，值为0时是未知
	Province       string     `json:"province" gorm:"size:50"`                  //用户个人资料填写的省份
	City           string     `json:"city" gorm:"size:50"`                      //普通用户个人资料填写的城市
	Country        string     `json:"country" gorm:"size:50"`                   //国家，如中国为CN
	Headimgurl     string     `json:"headimgurl" gorm:"size:300"`               //用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。
	Privilege      string     `json:"privilege"`                                //用户特权信息，json 数组，如微信沃卡用户为（chinaunicom）
	SubscribTime   int64      `json:"subscribe_time"`                           //用户关注时间，为时间戳。如果用户曾多次关注，则取最后关注时间
	SubscribeScene string     `json:"subscribe_scene"`                          //返回用户关注的渠道来源，ADD_SCENE_SEARCH 公众号搜索，ADD_SCENE_ACCOUNT_MIGRATION 公众号迁移，ADD_SCENE_PROFILE_CARD 名片分享，ADD_SCENE_QR_CODE 扫描二维码，ADD_SCENE_PROFILE_ LINK 图文页内名称点击，ADD_SCENE_PROFILE_ITEM 图文页右上角菜单，ADD_SCENE_PAID 支付后关注，ADD_SCENE_OTHERS 其他
	QrScene        string     `json:"qr_scene"`                                 //返回用户关注的渠道来源，ADD_SCENE_SEARCH 公众号搜索，ADD_SCENE_ACCOUNT_MIGRATION 公众号迁移，ADD_SCENE_PROFILE_CARD 名片分享，ADD_SCENE_QR_CODE 扫描二维码，ADD_SCENE_PROFILE_ LINK 图文页内名称点击，ADD_SCENE_PROFILE_ITEM 图文页右上角菜单，ADD_SCENE_PAID 支付后关注，ADD_SCENE_OTHERS 其他
	QrSceneStr     string     `json:"qr_scene_str"`                             //返回用户关注的渠道来源，ADD_SCENE_SEARCH 公众号搜索，ADD_SCENE_ACCOUNT_MIGRATION 公众号迁移，ADD_SCENE_PROFILE_CARD 名片分享，ADD_SCENE_QR_CODE 扫描二维码，ADD_SCENE_PROFILE_ LINK 图文页内名称点击，ADD_SCENE_PROFILE_ITEM 图文页右上角菜单，ADD_SCENE_PAID 支付后关注，ADD_SCENE_OTHERS 其他
	Subscribe      int        `json:"subscribe" gorm:"type:smallint;DEFAULT:0"` //0未关注1关注
	CreateTime     *time.Time `gorm:"type:timestamptz(0);DEFAULT:now()"`        // 创建时间
} // 微信用户信息
type WxMqUserInfo struct {
	WxUserInfo
} // 微信用户信息
type WxPcUserInfo struct {
	WxUserInfo
} // 微信用户信息

type WxTokenVo struct {
	WxToken
	Echostr   string `json:"echostr"`
	Nonce     string `json:"nonce"`
	Signature string `json:"signature"`
	Timestamp string `json:"timestamp"`
}

func (w *WxToken) Verification() interface{} {
	// 判断防伪码生成规则
	m := make(map[string]string)
	ms := make([]string, 0)
	_, err := url.Parse(w.URL)
	if err != nil || len(w.URL) > 300 {
		m["URL"] = "URL不合格，长度不能超过300"
		ms = append(ms, m["URL"])
	}

	if !util.CheckRegexp(w.MqAppID, `^[A-Za-z0-9]{16,20}$`) {
		m["MqAppID"] = "MqAppID必须为英文或数字，长度为16-20字符"
		ms = append(ms, m["mqAppID"])
	}
	if !util.CheckRegexp(w.MqAppSecret, `^[A-Za-z0-9]{32}$`) {
		m["mqAppSecret"] = "MqAppSecret必须为英文或数字，长度为32字符"
		ms = append(ms, m["mqAppSecret"])
	}
	if !util.CheckRegexp(w.Token, `^[A-Za-z0-9]{3,32}$`) {
		m["token"] = "token必须为英文或数字，长度为3-32字符"
		ms = append(ms, m["token"])
	}

	if w.EncodingAESType > 2 {
		m["encodingAESType"] = "加密方式错误"
		ms = append(ms, m["encodingAESType"])
	}

	if !util.CheckRegexp(w.EncodingAESKey, `^[A-Za-z0-9]{43}$`) {
		m["encodingAESKey"] = "消息加密密钥由43位字符组成，可随机修改，字符范围为A-Z，a-z，0-9。"
		ms = append(ms, m["encodingAESKey"])
	}

	if len(ms) > 0 {
		//return [3]interface{}{m, ms, strings.Join(ms, "\n")}
		return strings.Join(ms, "\n")
	}
	return nil
}
func (w *WxToken) GetMqAccessToken() {
	//v,err := oredis.Ser.Get("wx:mq:"+w.Id+":accessToken")
	//if err {
	//
	//}
}

/*
获取access_token
接口调用请求说明
https请求方式: GET https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=APPID&secret=APPSECRET
参数说明
参数		是否必须		说明
grant_type	是			获取access_token填写client_credential
appid		是			第三方用户唯一凭证
secret		是			第三方用户唯一凭证密钥，即appsecret
*/

type WxAccessToken struct {
	AccessToken  string `json:"access_token"`  //获取到的凭证
	Ticket       string `json:"ticket"`        //获取到的凭证
	ExpiresIn    int    `json:"expires_in"`    //凭证有效时间，单位：秒
	Errcode      int    `json:"errcode"`       //返回码
	Errmsg       string `json:"errmsg"`        //说明
	RefreshToken string `json:"refresh_token"` //填写通过access_token获取到的refresh_token参数
	Openid       string `json:"openid"`        //用户的唯一标识
	Unionid      string `json:"unionid"`       //用户的唯一标识
	Scope        string `json:"scope"`         //说明
}

type WxMqMsg struct {
	MsgId        string //消息id，64位整型
	MsgType      string //消息类型，文本为 text
	ToUserName   string //开发者微信号
	FromUserName string //发送方帐号（一个OpenID）
	CreateTime   int64  //消息创建时间 （整型）
}

type WxMqMsgAll struct {
	WxMqMsg
	Content      string //文本消息内容
	MediaId      string //图片消息媒体id，可以调用获取临时素材接口拉取数据。
	PicUrl       string //图片链接（由系统生成）
	Format       string //语音格式，如amr，speex等
	Recognition  string //语音识别结果，UTF8编码
	ThumbMediaId string //图片链接（由系统生成）f
	Location_X   string //图片链接（由系统生成）
	Location_Y   string //图片链接（由系统生成）
	Scale        string //图片链接（由系统生成）
	Label        string //图片链接（由系统生成）
	Title        string //消息标题
	Description  string //消息描述
	Url          string //消息链接
}

// 文本消息 text
type WxMqMsgText struct {
	WxMqMsg
	Content string //文本消息内容
}
type WxMqMsgMedia struct {
	WxMqMsg
	MediaId string //图片消息媒体id，可以调用获取临时素材接口拉取数据。
}

// 图片消息 image
type WxMqMsgMediaImg struct {
	WxMqMsgMedia
	PicUrl string //图片链接（由系统生成）
}

// 语音消息 voice
type WxMqMsgMediaVoice struct {
	WxMqMsgMedia
	Format      string //语音格式，如amr，speex等
	Recognition string //语音识别结果，UTF8编码
}

// 视频消息 video / shortvideo 小视频消息
type WxMqMsgMediaVideo struct {
	WxMqMsgMedia
	ThumbMediaId string //图片链接（由系统生成）
}

//地理位置消息 location
type WxMqMsgLocation struct {
	WxMqMsg
	Location_X string //地理位置维度
	Location_Y string //地理位置经度
	Scale      string //地图缩放大小
	Label      string //地理位置信息
}

//链接消息 link
type WxMqMsgLink struct {
	WxMqMsg
	Title       string //消息标题
	Description string //消息描述
	Url         string //消息链接
}

type WxMqEvent struct {
	MsgType      string //消息类型，文本为 text
	ToUserName   string //开发者微信号
	FromUserName string //发送方帐号（一个OpenID）
	CreateTime   int64  //消息创建时间 （整型）
}
type WxMqEventAll struct {
	MsgType      string //消息类型，文本为 text
	ToUserName   string //开发者微信号
	FromUserName string //发送方帐号（一个OpenID）
	CreateTime   int64  //消息创建时间 （整型）
	Event        string //事件类型，subscribe(订阅)、unsubscribe(取消订阅)
	EventKey     string //事件类型，subscribe(订阅)、unsubscribe(取消订阅)
	Ticket       string //事件类型，subscribe(订阅)、unsubscribe(取消订阅)
	Latitude     string //事件类型，subscribe(订阅)、unsubscribe(取消订阅)
	Longitude    string //事件类型，subscribe(订阅)、unsubscribe(取消订阅)
	Precision    string //事件类型，subscribe(订阅)、unsubscribe(取消订阅)
}

// 关注取消关注事件
// 自定义菜单事件\点击菜单跳转链接时的事件推送
type WxMqEventSub struct {
	WxMqEvent
	Event    string //事件类型，subscribe(订阅)、unsubscribe(取消订阅)
	EventKey string //事件类型，subscribe(订阅)、unsubscribe(取消订阅)
}

// 关注
type WxMqEventSubscribe struct {
	WxMqEventSub
	Ticket string //事件类型，subscribe(订阅)、unsubscribe(取消订阅)
}

// 上报地理位置事件
type WxMqEventLatitude struct {
	WxMqEventSub
	Latitude  string //事件类型，subscribe(订阅)、unsubscribe(取消订阅)
	Longitude string //事件类型，subscribe(订阅)、unsubscribe(取消订阅)
	Precision string //事件类型，subscribe(订阅)、unsubscribe(取消订阅)
}

var SubscribeScene = map[string]string{
	"ADD_SCENE_SEARCH":            "公众号搜索",
	"ADD_SCENE_ACCOUNT_MIGRATION": "公众号迁移",
	"ADD_SCENE_PROFILE_CARD":      "名片分享",
	"ADD_SCENE_QR_CODE":           "扫描二维码",
	"ADD_SCENE_PROFILE_ LINK":     "图文页内名称点击",
	"ADD_SCENE_PROFILE_ITEM":      "图文页右上角菜单",
	"ADD_SCENE_PAID":              "支付后关注",
	"ADD_SCENE_OTHERS":            "其他",
}

func (v *WxMqUserInfo) To() WxMqUserInfoVo {
	return WxMqUserInfoVo{Uid: v.Uid,
		Unionid:    v.Unionid,
		Openid:     v.Openid,
		Nickname:   v.Nickname,
		Sex:        v.Sex,
		Province:   v.Province,
		City:       v.City,
		Country:    v.Country,
		Headimgurl: v.Headimgurl,
		Privilege:  strings.Split(v.Privilege, ","),
	}
}

type WxMqUserInfoVo struct {
	Uid        string   `json:"uid"`        //主键ID
	Unionid    string   `json:"unionid"`    //只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段。
	Openid     string   `json:"openid"`     //用户的唯一标识
	Nickname   string   `json:"nickname"`   //用户昵称
	Sex        int      `json:"sex"`        //用户的性别，值为1时是男性，值为2时是女性，值为0时是未知
	Province   string   `json:"province"`   //用户个人资料填写的省份
	City       string   `json:"city"`       //普通用户个人资料填写的城市
	Country    string   `json:"country"`    //国家，如中国为CN
	Headimgurl string   `json:"headimgurl"` //用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。
	Privilege  []string `json:"privilege"`  //用户特权信息，json 数组，如微信沃卡用户为（chinaunicom）
}

func (v *WxMqUserInfoVo) To() WxMqUserInfo {
	return WxMqUserInfo{WxUserInfo{Uid: v.Uid,
		Unionid:    v.Unionid,
		Openid:     v.Openid,
		Nickname:   v.Nickname,
		Sex:        v.Sex,
		Province:   v.Province,
		City:       v.City,
		Country:    v.Country,
		Headimgurl: v.Headimgurl,
		Privilege:  strings.Join(v.Privilege, ",")},
	}
}

var (
	//URLHost = "http://dwx.free.idcfengye.com/"
	//  若提示“该链接无法访问”，请检查参数是否填写错误，是否拥有scope参数对应的授权作用域权限。
	URLWxMqOAuth2Authorize       = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&state=%s&scope=%s#wechat_redirect"
	URLWxCgiBinToken             = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	URLWxCgiBinUserInfo          = `https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=zh_CN`
	CgiBinMenuCreate             = `https://api.weixin.qq.com/cgi-bin/menu/create?access_token=%s`               // 菜单创建接口
	CgiBinGetCurrentSelfMenuInfo = `https://api.weixin.qq.com/cgi-bin/get_current_selfmenu_info?access_token=%s` // 菜单查询接口
	CgiBinMenuDelete             = `https://api.weixin.qq.com/cgi-bin/menu/delete?access_token=%s`               // 菜单删除接口
	URLWxMqMessageTemplateSend   = `https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s`
	URLWxJsApiGetTicket          = `https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi`
	URLWxConnectArconnect        = `https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect`
	URLWxSnsOauth2RefreshToken   = "https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=%s&grant_type=refresh_token&refresh_token=%s"
	URLWxSnsOauth2AccessToken    = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	URLWxSnsUserInfo             = "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN" //http：GET（请使用https协议）
)

/*
func Event(v model.WxMqMsgAnd) {
	switch v.Event {
	case "subscribe": // 关注
		break
	case "unsubscribe": // 取消关注
		break
	case "LOCATION": // 地理位置
		break
	case "CLICK": // 自定义菜单事件
		break
	case "VIEW": // 点击菜单跳转链接时的事件推送
		break
	}
}*/
func (w *WxToken) GetAccessToken() (accessToken WxAccessToken, err error) {
	k := w.MqAppID + constant.RedisKeyWxSnsAccessToken
	e := oredis.Ser.GetType(k, &accessToken)
	if e != nil || len(accessToken.AccessToken) < 20 {
		accessToken, err = w.SetAccessToken(k)
		if err != nil {
			return
		}
	}
	return
}
func (w *WxToken) GetAccessTokenStr() (token string, err error) {
	a, err := w.GetAccessToken()
	if err != nil {
		return
	}
	token = a.AccessToken
	return
}

func (w *WxToken) GetJsApiGetTicket() (token string, err error) {
	k := w.MqAppID + constant.RedisKeyWxJsApiGetTicket
	accessToken := WxAccessToken{}
	e := oredis.Ser.GetType(k, &accessToken)
	if e != nil || len(accessToken.AccessToken) < 20 {
		token, err = w.SetJsApiGetTicket(k)
		if err != nil {
			return
		}
		return
	}
	return accessToken.Ticket, nil
}
func (w *WxToken) GetAccessTokenNew() (accessToken WxAccessToken, err error) {
	http := util.HttpOk{Url: fmt.Sprintf(URLWxCgiBinToken, w.MqAppID, w.MqAppSecret)}
	http.Query()
	v := http.ResponseBody
	util.JsonToType(v, &accessToken)
	if len(accessToken.Errmsg) > 2 {
		err = errors.New(v)
		return
	}
	return accessToken, nil
}
func (w *WxToken) SetAccessToken(key string) (token WxAccessToken, err error) {
	accessToken, err := w.GetAccessTokenNew()
	if err != nil {
		return
	}
	err = oredis.Ser.SetType(key, accessToken, 110*time.Minute)
	if err != nil {
		return
	}
	return accessToken, nil
}

func (w *WxToken) SetJsApiGetTicket(key string) (token string, err error) {
	accessToken, err := w.GetJsApiGetTicketNew()
	if err != nil {
		return
	}
	err = oredis.Ser.SetType(key, accessToken, 110*time.Minute)
	if err != nil {
		return
	}
	return accessToken.Ticket, nil
}
func (w *WxToken) GetJsApiGetTicketNew() (accessToken WxAccessToken, err error) {
	token, err := w.GetAccessTokenStr()
	if err != nil {
		return
	}
	URL := strings.ReplaceAll(URLWxJsApiGetTicket, "ACCESS_TOKEN", token)
	http := util.HttpOk{Url: URL}
	http.QueryGet()
	v := http.ResponseBody
	util.JsonToType(v, &accessToken)
	if len(accessToken.Errmsg) > 2 {
		err = errors.New(v)
		return
	}
	//oredis.Ser.SetKV(,"")
	return accessToken, nil
}

func (w *WxToken) MenuCreate(menu interface{}) (str string, err error) {
	ACCESS_TOKEN, err := w.GetAccessTokenStr()
	if err != nil {
		return
	}
	http := util.HttpOk{Url: fmt.Sprintf(CgiBinMenuCreate, ACCESS_TOKEN), Method: "POST"}
	http.Params = menu
	http.Query()
	v := http.ResponseBody
	return v, nil
}

func (w *WxToken) MenuInfo() (str string, err error) {
	ACCESS_TOKEN, err := w.GetAccessTokenStr()
	if err != nil {
		return
	}
	http := util.HttpOk{Url: fmt.Sprintf(CgiBinGetCurrentSelfMenuInfo, ACCESS_TOKEN), Method: "GET"}
	http.Query()
	v := http.ResponseBody
	return v, nil
}
func (w *WxToken) MenuDelete() (str string, err error) {
	ACCESS_TOKEN, err := w.GetAccessTokenStr()
	if err != nil {
		return
	}
	http := util.HttpOk{Url: fmt.Sprintf(CgiBinMenuDelete, ACCESS_TOKEN), Method: "GET"}
	http.Query()
	v := http.ResponseBody
	return v, nil
}

// TODO 网页用户授权

// 第一步：用户同意授权，获取code
func (w *WxToken) GetURLWxMqOAuth2Authorize(state string, REDIRECT_URI string) string {
	return fmt.Sprintf(URLWxMqOAuth2Authorize, w.MqAppID, url.QueryEscape(REDIRECT_URI), state)
}

// 静默方式登录
func (w *WxToken) GetURLWxMqOAuth2AuthorizeBase(state, scope string) string {
	return fmt.Sprintf(URLWxMqOAuth2Authorize, w.MqAppID, url.QueryEscape(w.Host+"/wx/authorize?scope="+scope), state, scope)
}
func (w *WxToken) GetURLWxConnectArconnect(redirect_uri, state string) string {
	return fmt.Sprintf(URLWxConnectArconnect, w.MqAppID, url.QueryEscape(redirect_uri), state)
}

//第二步：通过code换取网页授权access_token
func WxSnsOAuth2AccessToken(appID, appSecret, code string) (at WxAccessToken, err error) {
	h := util.HttpOk{Url: fmt.Sprintf(URLWxSnsOauth2AccessToken, appID, appSecret, code)}
	h.Query()
	resp := h.ResponseBody
	if strings.Index(resp, "errcode") > 0 {
		err = errors.New(resp)
		return
	}
	util.JsonToType(resp, &at)
	return
}

// 获取用户信息
func GetWxSnsUserInfo(w WxAccessToken) (user WxUserInfo, err error) {
	h := util.HttpOk{Url: fmt.Sprintf(URLWxSnsUserInfo, w.AccessToken, w.Openid)}
	err = h.Query()
	if err != nil {
		return
	}
	resp := h.ResponseBody
	if strings.Index(resp, "errcode") > 0 {
		err = errors.New(resp)
		return
	}
	util.JsonToType(resp, &user)
	return
}

// 获取关注用户的信息
func GetWxCgiBinUserInfo(token, openId string) (user WxUserInfo, err error) {
	h := util.HttpOk{Url: fmt.Sprintf(URLWxCgiBinUserInfo, token, openId)}
	err = h.Query()
	if err != nil {
		return
	}
	resp := h.ResponseBody
	if strings.Index(resp, "errcode") > 0 {
		err = errors.New(resp)
		return
	}
	util.JsonToType(resp, &user)
	return
}
func GetWxPcCgiBinUserInfo(token, openId string) (user WxUserInfo, err error) {
	h := util.HttpOk{Url: fmt.Sprintf(URLWxCgiBinUserInfo, token, openId)}
	err = h.Query()
	if err != nil {
		return
	}
	resp := h.ResponseBody
	if strings.Index(resp, "errcode") > 0 {
		err = errors.New(resp)
		return
	}
	util.JsonToType(resp, &user)
	return
}

type WxMqMsgTemp struct {
	OpenId        string `json:"openId"`        // 接收信息用户的openID
	TempId        string `json:"tempId"`        // 模板ID
	Url           string `json:"url"`           // 跳转路径
	MinAppId      string `json:"minAppId"`      //小程序ID
	MinAppPath    string `json:"minAppPath"`    //小程序路径
	FirstValue    string `json:"firstValue"`    //标题值
	FirstColor    string `json:"firstColor"`    //标题颜色
	Keyword1Value string `json:"keyword1Value"` //第一个参数
	Keyword1Color string `json:"keyword1Color"` //第一个颜色
	Keyword2Value string `json:"keyword2Value"` //第二个参数
	Keyword2Color string `json:"keyword2Color"` //第二个颜色
	Keyword3Value string `json:"keyword3Value"` //第三个参数
	Keyword3Color string `json:"keyword3Color"` //第三个颜色
	RemarkValue   string `json:"remarkValue"`   // 最后值
	RemarkColor   string `json:"remarkColor"`   // 最后值颜色
}

/*
`{
           "touser":"OPENID",
           "template_id":"TEMPLATE_ID",
           "url":"URL",
           "miniprogram":{
             "appid":"MINI_APPID",
             "pagepath":"PAGEPATH"
           },
           "data":{
                   "first": {
                       "value":"FIRST_VALUE",
                       "color":"FIRST_COLOR"
                   },
                   "keyword1":{
                       "value":"KEYWORD1_VALUE",
                       "color":"KEYWORD1_COLOR"
                   },
                   "keyword2": {
                       "value":"KEYWORD2_VALUE",
                       "color":"KEYWORD2_COLOR"
                   },
                   "keyword3": {
                       "value":"KEYWORD3_VALUE",
                       "color":"KEYWORD3_COLOR"
                   },
                   "remark":{
                       "value":"REMARK_VALUE",
                       "color":"REMARK_COLOR"
                   }
           }
	}`

*/
// TODO 微信消息模板
func (w *WxToken) WxMqMessageTemplateSend(data WxMqMsgTemp) (v string, err error) {
	token, err := w.GetAccessTokenStr()
	if err != nil {
		return
	}
	URL := strings.ReplaceAll(URLWxMqMessageTemplateSend, "ACCESS_TOKEN", token)
	http := util.HttpOk{Url: URL, Method: `POST`}
	temp := `{
           "touser":"` + data.OpenId + `",
           "template_id":"` + data.TempId + `",
           "url":"` + data.Url + `",
           "miniprogram":{
             "appid":"` + data.MinAppId + `",
             "pagepath":"` + data.MinAppPath + `"
           },
           "data":{
                   "first": {
                       "value":"` + data.FirstValue + `",
                       "color":"` + data.FirstColor + `"
                   },
                   "keyword1":{
                       "value":"` + data.Keyword1Value + `",
                       "color":"` + data.Keyword1Color + `"
                   },
                   "keyword2": {
                       "value":"` + data.Keyword2Value + `",
                       "color":"` + data.Keyword2Color + `"
                   },
                   "keyword3": {
                       "value":"` + data.Keyword3Value + `",
                       "color":"` + data.Keyword3Color + `"
                   },
                   "remark":{
                       "value":"` + data.RemarkValue + `",
                       "color":"` + data.RemarkColor + `"
                   }
           }
	}`

	http.Params = temp
	fmt.Println("微信消息通知：")
	fmt.Println(util.JsonToStr(data))
	http.Query()
	v = http.ResponseBody
	fmt.Println(v)
	if strings.Index(v, `"errmsg":"ok"`) < 0 {
		err = errors.New(v)
	}
	return
}

type WxShortUrl struct {
	ErrCode  int    `json:"errcode"`
	ErrMsg   string `json:"errmsg"`
	ShortUrl string `json:"short_url"`
}

func (w *WxToken) Shorturl(langUrl string) (string, error) {
	ACCESS_TOKEN, err := w.GetAccessTokenStr()
	if err != nil {
		return "", err
	}
	uri := `https://api.weixin.qq.com/cgi-bin/shorturl?access_token=` + ACCESS_TOKEN
	http := util.HttpOk{Url: uri, Method: "POST"}
	http.QueryJsonByte([]byte(`{"action":"long2short","long_url":"` + langUrl + `"}`))
	v := http.ResponseBody
	var s WxShortUrl
	util.JsonToType(v, &s)
	if s.ErrMsg != "ok" {
		return "", errors.New(v)
	}
	return s.ShortUrl, nil
} //微信短连接生成
