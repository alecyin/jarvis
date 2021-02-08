package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/golang/glog"
	"net/http"
	"net/url"
	"os"
	"time"
)

type BotApi struct {
	ProxyAddr string
	Token     string
}

func setProxy() {
	os.Setenv("http_proxy", cfg.BotApi.ProxyAddr)
	os.Setenv("https_proxy", cfg.BotApi.ProxyAddr)
}
func unsetProxy() {
	os.Unsetenv("http_proxy")
	os.Unsetenv("https_proxy")
}
func RunTgBotApi() {
	//setProxy()
	//defer unsetProxy()
	transport := &http.Transport{Proxy: func(_ *http.Request) (*url.URL, error) {
		return url.Parse(cfg.BotApi.ProxyAddr)
	}}
	bot, err := tgbotapi.NewBotAPIWithClient(cfg.BotApi.Token, &http.Client{Transport: transport, Timeout: 10 * time.Second})
	if err != nil {
		glog.Error(err)
		return
	}

	bot.Debug = true

	glog.Info("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		glog.Info("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		if update.Message.Text == "切换节点" {
			if ChangeNode() {
				msg.Text = "done"
			} else {
				msg.Text = "failure"
			}
		}
		//msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
	}
}
