package main

import (
	"github.com/golang/glog"
	"net/smtp"
	"strings"
)

const qqSmtpHost = "smtp.qq.com"
const qqSmtpPort = ":587"

type QqMail struct {
	FromAccount string
	AuthCode    string
	ToAccount   string
}

func (*QqMail) ConsumeMsg(message Message) interface{} {
	if cfg == nil || cfg.QqMail == (QqMail{}) {
		glog.Error("none of qq mail config")
		return failure
	}

	qqMail := cfg.QqMail
	fromAccount := qqMail.FromAccount
	toAccount := qqMail.ToAccount
	authCode := qqMail.AuthCode

	to := []string{toAccount}
	contentType := "Content-Type: text/plain; charset=UTF-8"
	if message.MailName == "" {
		message.MailName = fromAccount
	}
	msg := []byte("To: " + strings.Join(to, ",") + "\r\nFrom: " + message.MailName +
		"<" + fromAccount + ">\r\nSubject: " + message.Title + "\r\n" + contentType + "\r\n\r\n" + message.Content)

	auth := smtp.PlainAuth("", fromAccount, authCode, qqSmtpHost)
	err := smtp.SendMail(qqSmtpHost+qqSmtpPort, auth, fromAccount, to, msg)
	if err != nil {
		glog.Error("send mail error: %v", err)
		return failure
	}
	return success
}
