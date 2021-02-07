package main

import (
	"github.com/golang/glog"
	"net/smtp"
	"strings"
)

const qqSmtpHost = "smtp.qq.com"
const qqSmtpPort = ":587"

func (qqMail QqMail) ConsumeMsg(message Message) interface{} {
	if cfg == nil || len(cfg.Scs) == 0 {
		glog.Fatal("none of sc config")
	}

	fromAccount := cfg.QqMails[0].FromAccount
	toAccount := cfg.QqMails[0].ToAccount
	authCode := cfg.QqMails[0].AuthCode

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
