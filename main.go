package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-ini/ini"
	"github.com/golang/glog"
)

func ParseConfig() {
	cfg = new(Cfg)
	cfgFile, err := ini.Load("config/config_new.ini")
	if err != nil {
		glog.Fatal("Fail to read file: ", err)
	}
	cfg.HttpPort = cfgFile.Section("http server").Key("port").String()
	sc := new(Sc)
	sc.Sckey = cfgFile.Section("sc").Key("SCKEY").String()
	cfg.Scs = append(cfg.Scs, *sc)
	qqMail := new(QqMail)
	qqMail.FromAccount = cfgFile.Section("qq_mail").Key("from_account").String()
	qqMail.ToAccount = cfgFile.Section("qq_mail").Key("to_account").String()
	qqMail.AuthCode = cfgFile.Section("qq_mail").Key("auth_code").String()
	cfg.QqMails = append(cfg.QqMails, *qqMail)
}

func main() {
	flag.Parse()
	glog.Info("start")
	ParseConfig()
	fmt.Println("port:", cfg.HttpPort)
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
		}

		var result interface{}
		if way == "sc" {
			var sc Sc
			result = sc.ConsumeMsg(message)
			if original == "0" {
				c.JSON(200, result)
				return
			}
			c.String(200, fmt.Sprintf("%v", result))
			return
		}

		if way == "qqmail" {
			mailName := c.DefaultQuery("name", "")
			message.MailName = mailName
			var qqMail QqMail
			result = qqMail.ConsumeMsg(message)
			c.JSON(200, result)
			return
		}
	})
	r.Run(":" + cfg.HttpPort) // listen and serve on 0.0.0.0:8080
}
