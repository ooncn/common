package util

import (
	"gopkg.in/gomail.v2"
	"strconv"
)

type EmailAcc struct {
	UserName string   `json:"userName"` //用户名
	Account  string   `json:"account"`  //账号
	Pass     string   `json:"pass"`     //密码
	Host     string   `json:"host"`     //服务器地址
	Port     string   `json:"Port"`     //端口号
	ToEmail  []string `json:"toEmail"`  //发送到邮箱集 发送给多个用户
	Subject  string   `json:"subject"`  // 设置邮件主题,
	Body     string   `json:"body"`     // 设置邮件正文,
}

func (e *EmailAcc) SendMail() error {
	port := 465 //转换端口类型为int
	//定义邮箱服务器连接信息，如果是阿里邮箱 pass填密码，qq邮箱填授权码
	if len(e.Port) > 0 {
		p, err := strconv.Atoi(e.Port)
		if err == nil {
			port = p
		}
	}
	username := e.Subject
	if len(e.UserName) > 0 {
		username = e.UserName
	}

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(e.Account, username))
	//这种方式可以添加别名，即“XD Game”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code> 读者可以自行实验下效果
	m.SetHeader("To", e.ToEmail...)
	m.SetHeader("Subject", e.Subject)
	m.SetBody("text/html", e.Body)

	d := gomail.NewDialer(e.Host, port, e.Account, e.Pass)
	err := d.DialAndSend(m)
	return err

}
