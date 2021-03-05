package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"net/http"
	"sync"
)

type WeChat struct {
	AppId      string
	AppSecret  string
	OpenId     string
	TemplateId string
	Token      string
}

var wechat *WeChat
var wOnce sync.Once

func GetWeChatIns() *WeChat {
	wOnce.Do(func() {
		wechat = newWeChat()
	})
	return wechat
}

func newWeChat() *WeChat {
	appId := cfg.cfgFile.Section("wechat").Key("app_id").String()
	appSecret := cfg.cfgFile.Section("wechat").Key("app_secret").String()
	openId := cfg.cfgFile.Section("wechat").Key("open_id").String()
	templateId := cfg.cfgFile.Section("wechat").Key("template_id").String()

	w := &WeChat{
		AppId:      appId,
		AppSecret:  appSecret,
		OpenId:     openId,
		TemplateId: templateId,
	}
	w.refreshToken()
	if _, err := GetMcronIns().cronEngine.AddFunc("@every 1h30m", w.refreshToken); err != nil {
		glog.Error("add func refresh wechat token error:", err)
	} else {
		glog.Info("refresh wechat token has been added to mcron")
	}
	return w
}

func (we *WeChat) refreshToken() {
	glog.Info("refresh wechat token begin")
	url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + we.AppId + "&secret=" + we.AppSecret
	resp, err := http.Get(url)
	if err != nil {
		glog.Error("refresh wechat token fail", err)
		return
	}
	defer resp.Body.Close()
	fmtRes, err := ParseResponse(resp)
	if err != nil {
		glog.Error("refresh wechat token fail", err)
		return
	}
	if _, ok := fmtRes["errcode"]; ok {
		glog.Error("refresh wechat token fail", fmtRes["errmsg"])
		return
	}
	we.Token = fmt.Sprintf("%v", fmtRes["access_token"])
	glog.Info("refresh wechat token success")
}

func (we *WeChat) ConsumeMsg(param interface{}) interface{} {
	message, ok := param.(Message)
	if !ok {
		return false
	}
	type wmsg struct {
		touser      string
		template_id string
		topcolor    string
		data        interface{}
	}
	wm := new(wmsg)
	wm.touser = we.OpenId
	wm.template_id = we.TemplateId
	title := struct {
		value string
		color string
	}{value: message.Title, color: "#173177"}
	content := struct {
		value string
		color string
	}{value: message.Content}
	data := struct {
		title   interface{}
		content interface{}
	}{title: title, content: content}
	wm.data = data
	jsonValue, _ := json.Marshal(wm)
	url := "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + we.Token
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		glog.Error(err)
		return false
	}
	defer resp.Body.Close()
	fmtRes, err := ParseResponse(resp)
	if err != nil {
		glog.Error(err)
		return false
	}
	return fmt.Sprintf("%v", fmtRes["errmsg"]) == "ok"
}
