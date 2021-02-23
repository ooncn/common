package obj

import (
	"bytes"
	"fmt"
	"github.com/ooncn/common/util"
	"html"
	"io"
	"os"
	"strings"
	"time"
)

// 附件表
type Att struct {
	Model      `json:"-" gorm:"-"`
	Id         string     `json:"id" gorm:"primary_key;size:33"`   //主键ID
	Uid        string     `json:"uid" gorm:"size:33;index"`        //上传用户ID
	Name       string     `json:"name" gorm:"size:150;index"`      //附件名称
	Tag        string     `json:"tag" gorm:"size:33;"`             //标签
	Path       string     `json:"path"`                            //物理路径
	Src        string     `json:"src"`                             //物理路径
	Ext        string     `json:"ext" gorm:"size:20;"`             //后缀名称
	Md5        string     `json:"md5" gorm:"size:58;unique_index"` //IP地址
	Ip         string     `json:"ip" gorm:"size:33;"`              //IP地址
	Size       int64      `json:"size" gorm:"DEFAULT:0"`           //占用空间
	Status     int        `json:"status" gorm:"DEFAULT:0"`         //状态 0等待上传 1成功 -1删除、不存在
	CreateTime *time.Time `json:"createTime"`                      //上传时间
}

func (a *Att) SaveUrl(uri string) {
	path := strings.Split(uri, "/")
	var name string
	if len(path) > 1 {
		name = path[len(path)-1]
	}
	// 文件名要格式化
	a.Name = html.EscapeString(name)
	p := util.GetCurrentDirectory() + "/uploads/" + a.Uid + util.TimeUtil.Sep("/yyyy/MM/dd/")
	_ = os.MkdirAll(p, os.ModePerm)
	p += name
	a.Path = p
	//History
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer f.Close()
	resp, err := util.GETUrlToByte(uri)
	if err != nil {
		return
	}
	size, err := io.Copy(f, bytes.NewReader(resp))
	a.Size = size
}

func (a *Att) SaveBase64(s string) {

}

/**
邮件模板
*/
type Email struct {
	Model   `json:"-" gorm:"-"`
	Id      string `json:"id"`      // 主键 // id String 主键
	Tag     string `json:"tag"`     // email_type String 邮件类型
	Name    string `json:"name"`    // sign_name String 短信签名
	Param   string `json:"param"`   // template_param String 邮件参数
	Content string `json:"content"` // template_content String 邮件模板
}

/**
邮件发送记录
*/
type EmailLog struct {
	Model      `json:"-" gorm:"-"`
	Id         string     `json:"id" gorm:"primary_key;size:33"`                       // 主键 // id Integer 主键
	Eid        string     `json:"eid" gorm:"size:33;index:Eid"`                        // email_id Integer 邮件模板
	FromEmail  string     `json:"fromEmail" gorm:"size:68"`                            // from_email String 发送邮箱
	ToEmail    string     `json:"toEmail" gorm:"size:68;index:ToEmail"`                // to_email String 接收邮箱
	Ip         string     `json:"ip" gorm:"size:33;index:Ip"`                          // ip String IP地址
	Title      string     `json:"title" gorm:"size:130"`                               // content String 邮箱标题
	Content    string     `json:"content"`                                             // content String 邮件内容
	Status     *int       `json:"status" gorm:"type:smallint;DEFAULT:0"`               // status String 状态
	CreateTime *time.Time `json:"createTime" gorm:"type:timestamptz(0);DEFAULT:now()"` //上传时间
}

//发票记录
type InvoiceLog struct {
	Model       `json:"-" gorm:"-"`
	Id          string     `json:"id" gorm:"primary_key;size:33"` // 主键   // id Integer
	FormUid     *int       // form_user_id Integer 开户人
	ToUid       *int       // to_user_id Integer 收票人
	InvoiceCode string     // invoice_code String 发票代码
	InvoiceNum  string     // invoice_num String 发票号码
	invoiceTime string     // invoice_time Date  开票时间
	CheckCode   string     // check_code String 校验码
	TwoName     string     // two_name String 买方名称
	TwoCode     string     // two_code String 纳税人识别号
	GoodsData   string     // goods_data String 商品信息
	Price       *float64   // price BigDecimal 价格合计
	OneName     string     // one_name String 销售方
	OneCode     string     // one_code String 纳税人识别号
	FilePath    string     // file_path String 发票文件地址
	ToTool      *int       // to_tool Integer 发送工具:邮件：备注写邮箱、快递：备注写单号
	Msg         string     // msg String 备注：发送工具:邮件：备注写邮箱、快递：备注写单号
	Status      *int       `json:"status" gorm:"type:smallint;DEFAULT:0"`               // status Integer 状态:0表示未审核、1审核通过、2问题
	CreateTime  *time.Time `json:"createTime" gorm:"type:timestamptz(0);DEFAULT:now()"` //上传时间
}

// 日志
type Log struct {
	Model      `json:"-" gorm:"-"`
	Id         string     `json:"id" gorm:"primary_key;size:33"`                       // 主键 // id String 主键ID
	Type       string     `json:"type" gorm:"size:33;"`                                // type String 请求方式
	Uid        string     `json:"uid" gorm:"size:33;index:Uid"`                        // user_id String 用户ID
	Parameter  string     `json:"parameter"`                                           // parameter String 参数列表
	Path       string     `json:"path"`                                                // path String 路径
	Ip         string     `json:"ip" gorm:"size:33;index:Ip"`                          // ip String ip地址
	Equipment  string     `json:"equipment"`                                           // equipment String 设备信息
	CreateTime *time.Time `json:"createTime" gorm:"type:timestamptz(0);DEFAULT:now()"` //上传时间
}

// 系统公告
type Message struct {
	Model      `json:"-" gorm:"-"`
	Id         string `json:"id" gorm:"primary_key;size:33"`       // 主键 // id String 主键
	Uid        string `json:"uid" gorm:"size:33;index:Uid"`        // uid 接收方
	MsgType    *int   `json:"type" gorm:"type:smallint;DEFAULT:0"` // msg_type Integer 消息类型
	Title      string // title String 账号
	Url        string // url String 连接地址
	Msg        string // msg String 内容通知
	Status     *int   `json:"status" gorm:"type:smallint;DEFAULT:0"` // status Integer 状态：0未查阅、1已查阅
	CreateTime string // create_time Date  创建时间
}

// 校验邮箱
type VerificationEmail struct {
	Token string `json:"token"` // 授权token
	Email string `json:"email"`
}

type Code struct {
	Tag   string `json:"tag"`
	Code  string `json:"code"`
	Del   int    `json:"del"`
	Phone string `json:"phone"`
	Md5   string `json:"md5"`
}

type CodeForm struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
	SmsCode  string `json:"smsCode"`
}

func RedisKeySmsVerification(phone string) string {
	return fmt.Sprintf("sms:verification:account:%v", phone)
} //短信验证码

type UserInfoVo struct {
	Id       string     `json:"id"`       //用户主键ID
	Account  string     `json:"account"`  //账号 4-32
	Nickname string     `json:"nickname"` //昵称
	HeadImg  string     `json:"head_img"` //昵称
	Phone    string     `json:"phone"`    //手机号
	Token    string     `json:"token"`    //密码
	Birthday string     `json:"birthday"` //生日
	Sex      int        `json:"sex"`      //性别：用户的性别，值为1时是男性，值为2时是女性，值为0时是未知
	Status   int        `json:"status"`   //状态 0未激活、1激活、2原始密码、-1、删除
	Login    *UserLogin `json:"login"`    // 登录信息
	Type     int        `json:"type"`     // 登录信息
	Info     UserInfo
}
type UserInfo struct {
	Id         string     `json:"id"`          // 主键 // user_id String 会员ID
	Gid        string     `json:"gid"`         // group_id Integer 会员组
	Nickname   string     `json:"nickname"`    // nickname String 会员昵称
	Sex        int        `json:"sex"`         // sex Integer 性别 0 未知 1男 2女
	Country    string     `json:"country"`     // country String 国家
	Province   string     `json:"province"`    // province String 省份
	City       string     `json:"city"`        // city String 城市
	HeadImgUrl string     `json:"headImgUrl"`  // headimgurl String 头像
	Birthday   string     `json:"birthday"`    // birthday Date  生日
	IP         string     `json:"ip"`          //注册IP地址
	CreateTime *time.Time `json:"create_time"` //创建时间
}
type UserLogin struct {
	Id         string     `json:"id"`          //用户主键ID
	Uid        string     `json:"uid"`         //用户ID
	Ip         string     `json:"ip"`          //登录IP
	Ua         string     `json:"ua"`          //设备信息
	Tag        string     `json:"tag"`         //设备信息
	Gps        string     `json:"gps"`         //设备信息
	Status     int        `json:"status"`      //登录状态 0、失败；1、成功
	EndTime    int64      `json:"end_time"`    //登录时间
	CreateTime *time.Time `json:"create_time"` //创建时间
}

type Address struct {
	Model      `json:"-" gorm:"-"`
	Id         string     `json:"id" gorm:"primary_key;size:33"`                        // 主键 //主键ID
	Uid        string     `json:"uid" gorm:"size:33;index"`                             //user_id 用户id
	Phone      string     `json:"phone"`                                                //收件联系人
	Name       string     `json:"name"`                                                 //收件人
	Gps        string     `json:"gps"`                                                  // GPS定位
	Sheng      string     `json:"sheng"`                                                // 省
	Shi        string     `json:"shi"`                                                  // 市
	Qu         string     `json:"qu"`                                                   // 区县
	Address    string     `json:"address"`                                              // 详细地址
	Mo         *int       `json:"mo"`                                                   //默认地址：0、不是/1、是
	Tag        string     `json:"tag"`                                                  //标签
	Status     *int       `json:"status" gorm:"type:smallint;DEFAULT:0"`                //状态：0、未激活/1、激活/2、删除
	CreateTime *time.Time `json:"create_time" gorm:"type:timestamptz(0);DEFAULT:now()"` // create_time String 创建时间
}
type GPS struct {
	Type      string `json:"type"`      // gaode\qq\baidu
	Sheng     string `json:"sheng"`     // 省
	Shi       string `json:"shi"`       // 市
	Qu        string `json:"qu"`        // 区县
	Address   string `json:"address"`   // 详细地址
	Latitude  string `json:"latitude"`  //
	Longitude string `json:"longitude"` // 标签
}

type IdName struct {
	Id   string `json:"id"`   //主键
	Name string `json:"name"` //名字
} //字典//fv_dictionaries
