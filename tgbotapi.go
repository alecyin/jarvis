package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"net/http"
	"net/url"
)

type BotApi struct {
	ProxyAddr string
	Token     string
}

func RunTgBotApi() {
	transport := &http.Transport{Proxy: func(_ *http.Request) (*url.URL, error) {
		return url.Parse(cfg.BotApi.ProxyAddr)
	}}
	bot, err := tgbotapi.NewBotAPIWithClient(cfg.BotApi.Token, &http.Client{Transport: transport})
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
