package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-ini/ini"
	"github.com/golang/glog"
	"log"
)

func ParseConfig() {
	cfg = new(Cfg)
	cfgFile, err := ini.Load("config/config.ini")
	if err != nil {
		log.Fatal("Fail to read file: ", err)
	}
	cfg.HttpPort = cfgFile.Section("http server").Key("port").String()
	sc := new(Sc)
	sc.Sckey = cfgFile.Section("sc").Key("SCKEY").String()
	cfg.Scs = append(cfg.Scs, *sc)
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
	r.GET("/message", func(c *gin.Context) {
		way := c.DefaultQuery("way", "sc")
		title := c.DefaultQuery("title", "")
		content := c.DefaultQuery("content", "")

		if title == "" {
			c.JSON(400, gin.H{
				"msg": "Required fields are missing",
			})
		}

		message := Message{
			Title:   title,
			Content: content,
		}

		if way == "sc" {
			DoConsumeMsg(message)
		}
		c.JSON(200, gin.H{
			"msg": "ok",
		})
	})
	r.Run(":" + cfg.HttpPort) // listen and serve on 0.0.0.0:8080
}
