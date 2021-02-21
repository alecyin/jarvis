package main

import (
	"flag"
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
	cfg.cfgFile = cfgFile
	if err != nil {
		glog.Fatal("fail to read file: ", err)
	}
	//node file
	cfg.SsrConfigFile = cfgFile.Section("ssr").Key("config_file").String()
}

func waitExit() {
	c := make(chan os.Signal)
	//ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for s := range c {
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			glog.Info("exit,", s)
		default:
			glog.Info("other", s)
		}
	}
}

func main() {
	flag.Parse()
	defer glog.Flush()
	ParseConfig()
	go GetSsrIns().ServiceabilityTest()
	go NewHttpApi().Run()
	go GetTgBotIns().Run()
	waitExit()
}
