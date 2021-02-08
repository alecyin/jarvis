package main

import (
	"flag"
	"fmt"
	"github.com/go-ini/ini"
	"github.com/golang/glog"
	"os"
	"os/signal"
	"syscall"
)

const (
	configFilePath = "./config/config_new.ini"
	recordFilePath = "record/record.txt"
)

func ParseConfig() {
	cfg = new(Cfg)
	cfgFile, err := ini.Load(configFilePath)
	if err != nil {
		glog.Fatal("Fail to read file: ", err)
	}
	cfg.HttpPort = cfgFile.Section("http server").Key("port").String()
	// sc
	sckey := cfgFile.Section("sc").Key("SCKEY").String()
	cfg.Scs = append(cfg.Scs, Sc{Sckey: sckey})
	// qq mail
	fromAccount := cfgFile.Section("qq_mail").Key("from_account").String()
	toAccount := cfgFile.Section("qq_mail").Key("to_account").String()
	authCode := cfgFile.Section("qq_mail").Key("auth_code").String()
	cfg.QqMails = append(cfg.QqMails, QqMail{FromAccount: fromAccount, ToAccount: toAccount, AuthCode: authCode})
	//tg bot
	proxyAddr := cfgFile.Section("tg_bot").Key("proxy_addr").String()
	token := cfgFile.Section("tg_bot").Key("token").String()
	cfg.BotApi = BotApi{ProxyAddr: proxyAddr, Token: token}
}

func waitExit() {
	c := make(chan os.Signal)
	//ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for s := range c {
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			glog.Info("exit", s)
			fmt.Print("exit")
		default:
			glog.Info("other", s)
		}
	}
}

func main() {
	flag.Parse()
	ParseConfig()
	//go RunHttpApi()
	go RunTgBotApi()
	waitExit()
}
