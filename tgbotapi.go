package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/golang/glog"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"sync"
	"time"
)

type TgBot struct {
	ProxyAddr string
	Token     string
	ChatId    int64
	TgApiUrl  string
	*tgbotapi.BotAPI
}

var tgBot *TgBot
var tgOnce sync.Once

func GetTgBotIns() *TgBot {
	tgOnce.Do(func() {
		tgBot = newTgBot()
	})
	return tgBot
}

func newTgBot() *TgBot {
	proxyAddr := cfg.cfgFile.Section("tg_bot").Key("proxy_addr").String()
	token := cfg.cfgFile.Section("tg_bot").Key("token").String()
	chatId, _ := cfg.cfgFile.Section("tg_bot").Key("chat_id").Int64()
	t := &TgBot{
		ProxyAddr: proxyAddr,
		Token:     token,
		ChatId:    chatId,
		TgApiUrl:  "https://api.telegram.org/bot" + token,
	}
	return t
}

func (tgBot *TgBot) setProxy() {
	os.Setenv("http_proxy", tgBot.ProxyAddr)
	os.Setenv("https_proxy", tgBot.ProxyAddr)
}
func (tgBot *TgBot) unsetProxy() {
	os.Unsetenv("http_proxy")
	os.Unsetenv("https_proxy")
}
func (tgBot *TgBot) initBotConn() {
	//setProxy()
	//defer unsetProxy()
	// new bot api
	transport := &http.Transport{Proxy: func(_ *http.Request) (*url.URL, error) {
		return url.Parse(tgBot.ProxyAddr)
	}}
	waitTimes := 1
	for {
		bot, err := tgbotapi.NewBotAPIWithClient(tgBot.Token, &http.Client{Transport: transport})
		if err != nil {
			glog.Errorf("start tg bot err, %s", err)
			time.Sleep(time.Duration(waitTimes) * time.Second)
			waitTimes += waitTimes
			continue
		}
		bot.Debug = true
		tgBot.BotAPI = bot
		glog.Infof("Authorized on account [%s]", bot.Self.UserName)
		break
	}
}

func (tgBot *TgBot) Run() {
	tgBot.initBotConn()
	tgBot.addTgJobsToCron()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := tgBot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		glog.Infof("recive message from bot, [%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		switch update.Message.Text {
		case "切换节点":
			if GetSsrIns().ChangeNode() {
				msg.Text = "success"
			} else {
				msg.Text = "failure"
			}
		case "glassinfo", "glasserror", "jarvisinfo", "jarviserror":
			logCmd := ""
			if update.Message.Text == "glassinfo" {
				logCmd = "tail -12 /home/glass/log/main.INFO"
			} else if update.Message.Text == "glasserror" {
				logCmd = "tail -12 /home/glass/log/main.ERROR"
			} else if update.Message.Text == "jarvisinfo" {
				logCmd = "tail -12 /home/jarvis/log/jarvis.INFO"
			} else if update.Message.Text == "jarviserror" {
				logCmd = "tail -12 /home/jarvis/log/jarvis.ERROR"
			}
			cmd := exec.Command("/bin/bash", "-c", logCmd)
			out, err := cmd.Output()
			if err != nil {
				glog.Error(err)
				msg.Text = fmt.Sprintf("%s", err)
			} else {
				msg.Text = string(out)
			}
		default:
			msg.Text = update.Message.Text
		}
		//msg.ReplyToMessageID = update.Message.MessageID
		tgBot.Send(msg)
	}
}

func (tgBot *TgBot) SendToMe(messageConfig tgbotapi.MessageConfig) {
	if _, err := tgBot.Send(messageConfig); err != nil {
		glog.Error("tg send to me error,", err)
		return
	}
	glog.Info("tg send to me success")
	//tgbotapi.NewMessage(cfg.BotApi.ChatId, "推送")
}

func (tgBot *TgBot) SendToMeStr(message string) error {
	if _, err := tgBot.Send(tgbotapi.NewMessage(tgBot.ChatId, message)); err != nil {
		return err
	}
	return nil
}

func (tgBot *TgBot) addTgJobsToCron() {
	jobs := GetMcronIns().jobs
	dwj := NewDailyWordJob(jobs["TG_DailyWordJob"]["name"], jobs["TG_DailyWordJob"]["schedule"])
	if _, err := GetMcronIns().cronEngine.AddJob(dwj.Schedule, dwj); err != nil {
		glog.Error("add job ", dwj.Name, " error:", err)
	} else {
		glog.Info(dwj.Name, " has been added to mcron")
	}
}

func (tgBot *TgBot) ConsumeMsg(param interface{}) interface{} {
	message, ok := param.(Message)
	if !ok {
		return false
	}
	if err := tgBot.SendToMeStr(fmt.Sprintf("%s\n%s", message.Title, message.Content)); err != nil {
		glog.Error("tg send to me str error,", err)
		return false
	}
	glog.Info("tg send to me str success")
	return true
}
