package main

import (
	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
)

const url = "https://sc.ftqq.com/"

func (sc Sc) ConsumeMsg(message Message) interface{} {
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
		res, _ := ParseResponse(res)
		if fmt.Sprintf("%v", res["errno"]) == "0" { //success
			return success
		}
		return failure
	}

	body, _ := ioutil.ReadAll(res.Body)
	return string(body)
}
