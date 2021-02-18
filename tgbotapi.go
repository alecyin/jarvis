package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/golang/glog"
	"net/http"
	"net/url"
	"os"
	"os/exec"
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
	for k := 0; k < 10; k++ {
		bot, err := tgbotapi.NewBotAPIWithClient(cfg.TgBot.Token, &http.Client{Transport: transport})
		if err != nil {
			glog.Errorf("start tg bot err, %s", err)
			time.Sleep(time.Second * 20)
			continue
		}
		bot.Debug = true
		tgBot.BotAPI = bot
		glog.Info("Authorized on account %s", bot.Self.UserName)
		break
	}
	tgBot.addTgJobsToCron() // ensure initialized
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

		glog.Infof("recive message from bot, [%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		if update.Message.Text == "切换节点" {
			if ChangeNode() {
				msg.Text = "success"
			} else {
				msg.Text = "failure"
			}
		} else if update.Message.Text == "glassinfo" {
			cmd := exec.Command("/bin/bash", "-c", "tail -12 /home/glass/log/main.INFO")
			out, err := cmd.Output()
			if err != nil {
				glog.Error(err)
			}
			msg.Text = string(out)
		} else if update.Message.Text == "glasserror" {
			cmd := exec.Command("/bin/bash", "-c", "tail -12 /home/glass/log/main.ERROR")
			out, err := cmd.Output()
			if err != nil {
				glog.Error(err)
			}
			msg.Text = string(out)
		}
		//msg.ReplyToMessageID = update.Message.MessageID
		tgBot.Send(msg)
	}
}

func (tgBot *TgBot) SendToMe(messageConfig tgbotapi.MessageConfig) {
	if cfg == nil || cfg.TgBot.BotAPI == nil {
		glog.Error("tg bot not initialized")
		return
	}
	tgBot.Send(messageConfig)
	//tgbotapi.NewMessage(cfg.BotApi.ChatId, "推送")
}

func (tgBot *TgBot) SendToMeStr(message string) {
	if cfg == nil || cfg.TgBot.BotAPI == nil {
		glog.Error("tg bot not initialized")
		return
	}
	tgBot.Send(tgbotapi.NewMessage(cfg.TgBot.ChatId, message))
}

func (tgBot *TgBot) addTgJobsToCron() {
	jobs := cfg.Mcron.jobs
	dwj := NewDailyWordJob(jobs["TG_DailyWordJob"]["name"], jobs["TG_DailyWordJob"]["schedule"])
	if _, err := cfg.Mcron.cronEngine.AddJob(dwj.Schedule, dwj); err != nil {
		glog.Error("add job ", dwj.Name, " error:", err)
		fmt.Print(err)
	}
}

func (tgBot *TgBot) ConsumeMsg(message Message) interface{} {
	if cfg == nil || cfg.TgBot.BotAPI == nil {
		glog.Error("tg bot not initialized")
		return failure
	}
	tgBot.SendToMe(tgbotapi.NewMessage(tgBot.ChatId, fmt.Sprintf("%s\n%s", message.Title, message.Content)))
	return success
}
