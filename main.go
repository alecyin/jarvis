package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-ini/ini"
	"github.com/golang/glog"
	"os"
	"os/signal"
	"syscall"
)

const (
	configFilePath  = "config/config_new.ini"
	recordFilePath  = "record/record.txt"
	ssrNodeFilePath = "config/ssr.yml"
	cronCmdFilePath = "config/cron_cmd.conf"
	cronJobFilePath = "config/cron_job.conf"
)

func ParseConfig() {
	cfg = new(Cfg)
	cfgFile, err := ini.Load(configFilePath)
	if err != nil {
		glog.Fatal("fail to read file: ", err)
	}
	cfg.HttpPort = cfgFile.Section("http server").Key("port").String()
	// sc
	sckey := cfgFile.Section("sc").Key("SCKEY").String()
	cfg.Sc = Sc{Sckey: sckey}
	// qq mail
	fromAccount := cfgFile.Section("qq_mail").Key("from_account").String()
	toAccount := cfgFile.Section("qq_mail").Key("to_account").String()
	authCode := cfgFile.Section("qq_mail").Key("auth_code").String()
	cfg.QqMail = QqMail{FromAccount: fromAccount, ToAccount: toAccount, AuthCode: authCode}
	//tg bot
	proxyAddr := cfgFile.Section("tg_bot").Key("proxy_addr").String()
	token := cfgFile.Section("tg_bot").Key("token").String()
	chatId, err := cfgFile.Section("tg_bot").Key("chat_id").Int64()
	if err != nil {
		glog.Error("parse chatId error,", err)
	}
	cfg.TgBot = TgBot{ProxyAddr: proxyAddr, Token: token, ChatId: chatId, TgApiUrl: "https://api.telegram.org/bot" + token}
	//node file
	cfg.SsrConfigFile = cfgFile.Section("ssr").Key("config_file").String()
	//cron
	cronCmdFileContent, err := ReadTotalFile(cronCmdFilePath)
	if err != nil {
		glog.Error("load cron config error,", err)
	}
	var procInfos []ProcInfo
	if err = json.Unmarshal(cronCmdFileContent, &procInfos); err != nil {
		glog.Error("convert cron json array to struct error", err)
	}
	u := map[string]*ProcInfo{}
	for _, procInfo := range procInfos {
		p := procInfo
		u[procInfo.Name] = &p
	}
	cronJobFileContent, err := ReadTotalFile(cronJobFilePath)
	if err != nil {
		glog.Error("load cron config error,", err)
	}
	var jobs map[string]map[string]string
	if err = json.Unmarshal(cronJobFileContent, &jobs); err != nil {
		glog.Error(err)
	}
	cfg.Mcron = NewMcron(u, jobs)
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
	defer glog.Flush()
	ParseConfig()
	go RunHttpApi()
	go cfg.TgBot.Run()
	waitExit()
}
