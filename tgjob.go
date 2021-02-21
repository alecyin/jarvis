package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/golang/glog"
	"net/http"
	"time"
)

type TgJob struct {
	Name     string `json:"name"`
	Schedule string `json:"schedule"`
}

func newTgJob(name string, schedule string) *TgJob {
	return &TgJob{
		Name:     name,
		Schedule: schedule,
	}
}

func (tgJob *TgJob) startJob() {
	glog.Info("start job ", tgJob.Name)
}

func (tgJob *TgJob) endJob() {
	glog.Info("end job ", tgJob.Name)
}

// combine
type DailyWordJob struct {
	*TgJob
}

func NewDailyWordJob(name string, schedule string) *DailyWordJob {
	return &DailyWordJob{
		TgJob: newTgJob(name, schedule),
	}
}

func (dailyWordJob *DailyWordJob) Run() {
	dailyWordJob.startJob()
	defer dailyWordJob.endJob()
	res, err := Get("http://open.iciba.com/dsapi/?date="+time.Now().Format("2006-01-02"), nil, nil)
	if err != nil {
		glog.Error("daily english error", err)
		GetTgBotIns().SendToMe(tgbotapi.NewMessage(GetTgBotIns().ChatId, fmt.Sprintf("daily english error, %v", err)))
		return
	}
	fmtRes, err := ParseResponse(res)
	if err != nil {
		glog.Error("parse daily english resp error:", err)
		GetTgBotIns().SendToMe(tgbotapi.NewMessage(GetTgBotIns().ChatId, fmt.Sprintf("parse daily english resp error, %v", err)))
		return
	}
	data := make(map[string]interface{})
	data["caption"] = fmt.Sprintf("%v\n%v", fmtRes["content"], fmtRes["note"])
	data["photo"] = fmtRes["fenxiang_img"]
	data["chat_id"] = GetTgBotIns().ChatId

	url := GetTgBotIns().TgApiUrl + "/sendPhoto"
	jsonData, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	client := GetTgBotIns().Client
	resp, err := client.Do(req)
	if err != nil {
		glog.Error(err)
		return
	}
	defer resp.Body.Close()
	parseRes, _ := ParseResponse(resp)
	if parseRes["ok"].(bool) == true {
		glog.Info("send daily english success")
	} else {
		glog.Info("send daily english failure")
	}
}
