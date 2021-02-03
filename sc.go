package main

import "github.com/golang/glog"

const url = "https://sc.ftqq.com/"

func DoConsumeMsg(message Message) {
	if cfg == nil || len(cfg.Scs) == 0 {
		glog.Fatal("none of sc config")
	}
	sckey := cfg.Scs[0].Sckey
	scUrl := url + sckey + ".send"
	params := make(map[string]string)
	params["text"] = message.Title
	params["desp"] = message.Content
	res, err := Get(scUrl, params, nil)
	if err != nil {
		glog.Error("send to sc error", err)
	}
	glog.Info(res)
}
