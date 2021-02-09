package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/golang/glog"
	"net/http"
	"net/url"
	"os"
	"time"
)

type TgBot struct {
	ProxyAddr string
	Token     string
	ChatId    int64
	TgApiUrl  string
	*tgbotapi.BotAPI
}

func setProxy() {
	os.Setenv("http_proxy", cfg.TgBot.ProxyAddr)
	os.Setenv("https_proxy", cfg.TgBot.ProxyAddr)
}
func unsetProxy() {
	os.Unsetenv("http_proxy")
	os.Unsetenv("https_proxy")
}
func (tgBot *TgBot) Run() {
	//setProxy()
	//defer unsetProxy()
	// new bot api
	transport := &http.Transport{Proxy: func(_ *http.Request) (*url.URL, error) {
		return url.Parse(cfg.TgBot.ProxyAddr)
	}}
	bot, err := tgbotapi.NewBotAPIWithClient(cfg.TgBot.Token, &http.Client{Transport: transport, Timeout: 10 * time.Second})
	if err != nil {
		glog.Error(err)
		return
	}
	bot.Debug = true
	tgBot.BotAPI = bot
	glog.Info("Authorized on account %s", bot.Self.UserName)
	tgBot.listen()
}

func (tgBot *TgBot) listen() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := tgBot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		glog.Info("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		if update.Message.Text == "切换节点" {
			if ChangeNode() {
				msg.Text = "success"
			} else {
				msg.Text = "failure"
			}
		}
		//msg.ReplyToMessageID = update.Message.MessageID
		tgBot.Send(msg)
	}
	//tgbotapi.NewMessage(cfg.BotApi.ChatId, "推送")
}
