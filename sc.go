package main

import (
	"fmt"
	"github.com/golang/glog"
)

const scApiUrl = "https://sctapi.ftqq.com/"

type Sc struct {
	Sckey string
}

func NewSc() *Sc {
	sckey := cfg.cfgFile.Section("sc").Key("SCKEY").String()
	return &Sc{
		Sckey: sckey,
	}
}

func (sc *Sc) ConsumeMsg(param interface{}) interface{} {
	message, ok := param.(Message)
	if !ok {
		return false
	}
	scUrl := scApiUrl + sc.Sckey + ".send"
	params := make(map[string]string)
	params["text"] = message.Title
	params["desp"] = message.Content
	res, err := Get(scUrl, params, nil)
	if err != nil {
		glog.Error("send to sc error", err)
		return false
	}
	//dismiss original result
	r, _ := ParseResponse(res)
	return fmt.Sprintf("%v", r["error"]) == "SUCCESS"
}
