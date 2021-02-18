package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func RunHttpApi() {
	glog.Info("port:", cfg.HttpPort)
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg": "pong",
		})
	})
	r.GET("/push", func(c *gin.Context) {
		way := c.DefaultQuery("way", "sc")
		title := c.DefaultQuery("title", "")
		content := c.DefaultQuery("content", "")
		original := c.DefaultQuery("original", "0")
		glog.Info("receive message,way:", way, ",title:", title, ",content:", content)
		if title == "" {
			c.JSON(400, gin.H{
				"msg": "Required fields are missing",
			})
			return
		}
		message := Message{
			Title:    title,
			Content:  content,
			Original: original,
			Way:      way,
		}
		// choose strategic
		var consumer Consumer
		if way == "sc" {
			consumer.setWay(&Sc{})
		} else if way == "qqmail" {
			message.MailName = c.DefaultQuery("name", "")
			consumer.setWay(&QqMail{})
		} else if way == "tg" {
			consumer.setWay(&cfg.TgBot)
		}
		go RecordMessage(message)
		result := consumer.Send(message)
		if original == "0" { // dismiss original result
			c.JSON(200, result)
			return
		}
		c.String(200, fmt.Sprintf("%v", result))
		return
	})
	r.Run(":" + cfg.HttpPort)
}
