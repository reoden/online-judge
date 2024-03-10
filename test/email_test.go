package test

import (
	"crypto/tls"
	"github.com/jordan-wright/email"
	"net/smtp"
	"testing"
)

func TestSendEmail(t *testing.T) {
	e := email.NewEmail()
	e.From = "Get <staraino_o@163.com>"
	e.To = []string{"1224074124@qq.com"}
	e.Subject = "验证码发送测试"
	e.HTML = []byte("您的验证码：<b>123456</b>")
	// 返回 EOF 时，关闭SSL重试
	err := e.SendWithTLS("smtp.163.com:465",
		smtp.PlainAuth("", "staraino_o@163.com", "SIYICTUUUKWYNGKO", "smtp.163.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"})
	if err != nil {
		t.Fatal(err)
	}
}
