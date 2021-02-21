package main

import (
	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
)

const scApiUrl = "https://sc.ftqq.com/"

type Sc struct {
	Sckey string
}

func NewSc() *Sc {
	sckey := cfg.cfgFile.Section("sc").Key("SCKEY").String()
	return &Sc{
		Sckey: sckey,
	}
}

func (sc *Sc) ConsumeMsg(message Message) interface{} {
	scUrl := scApiUrl + sc.Sckey + ".send"
	params := make(map[string]string)
	params["text"] = message.Title
	params["desp"] = message.Content
	res, err := Get(scUrl, params, nil)
	if err != nil {
		glog.Error("send to sc error", err)
		return nil
	}
	//dismiss original result
	if message.Original == "0" {
		res, _ := ParseResponse(res)
		if fmt.Sprintf("%v", res["errno"]) == "0" { //success
			return success
		}
		return failure
	}

	body, _ := ioutil.ReadAll(res.Body)
	return string(body)
}
