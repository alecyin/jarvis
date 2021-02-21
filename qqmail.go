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

func NewQqMail() *QqMail {
	fromAccount := cfg.cfgFile.Section("qq_mail").Key("from_account").String()
	toAccount := cfg.cfgFile.Section("qq_mail").Key("to_account").String()
	authCode := cfg.cfgFile.Section("qq_mail").Key("auth_code").String()
	return &QqMail{
		FromAccount: fromAccount,
		ToAccount:   toAccount,
		AuthCode:    authCode,
	}
}

func (qqMail *QqMail) ConsumeMsg(message Message) interface{} {
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
