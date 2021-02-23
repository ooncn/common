package util

import (
	"fmt"
	"testing"
)

func TestEmail(t *testing.T) {
	//定义收件人
	mailTo := []string{
		"qq@qq.com",
	}
	//邮件主题为"Hello"
	subject := "Hello"
	// 邮件正文
	body := "Good"
	email := EmailAcc{
		UserName: "ONCMS",
		Account:  "qq@163.com",
		Pass:     "EPFISMMAPENMCMAW",
		Host:     "smtp.163.com",
		Port:     "465",
		ToEmail:  mailTo,
		Subject:  subject,
		Body:     body,
	}
	err := email.SendMail()
	if err != nil {
		fmt.Println(err)
	}
}
