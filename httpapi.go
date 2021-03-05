package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

type HttpApi struct {
	port string
}

func NewHttpApi() *HttpApi {
	port := cfg.cfgFile.Section("http server").Key("port").String()
	return &HttpApi{
		port: port,
	}
}

func (httpApi *HttpApi) Run() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg": "pong",
		})
	})
	r.GET("/changenode", func(c *gin.Context) {
		c.JSON(200, GetSsrIns().ChangeNode())
	})
	r.GET("/push", func(c *gin.Context) {
		way := c.DefaultQuery("way", "sc")
		title := c.DefaultQuery("title", "")
		content := c.DefaultQuery("content", "")
		glog.Info("receive message,way:", way, ",title:", title, ",content:", content)
		if title == "" {
			c.JSON(400, gin.H{
				"msg": "Required fields are missing",
			})
			return
		}
		message := Message{
			Title:   title,
			Content: content,
			Way:     way,
		}
		// choose strategic
		var consumer Consumer
		if way == "sc" {
			consumer.setWay(NewSc())
		} else if way == "qqmail" {
			message.MailName = c.DefaultQuery("name", "")
			consumer.setWay(NewQqMail())
		} else if way == "tg" {
			consumer.setWay(GetTgBotIns())
		}
		go RecordMessage(message)
		result := consumer.Send(message)
		c.String(200, fmt.Sprintf("%v", result))
		return
	})
	glog.Info("http run on port:", httpApi.port)
	r.Run(":" + httpApi.port)
}
