package main

import (
	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
)

const url = "https://sc.ftqq.com/"

func DoConsumeMsg(message Message) interface{} {
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
		return nil
	}
	//dismiss original result
	if message.Original == "0" {
		returnMap, _ := ParseResponse(res)
		if fmt.Sprintf("%v", returnMap["errno"]) == "0" { //success
			formatMap := map[string]int{"code": 0}
			return formatMap
		}
		formatMap := map[string]int{"code": -1}
		return formatMap
	}

	body, _ := ioutil.ReadAll(res.Body)
	return string(body)
}
