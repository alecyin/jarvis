package main

import "github.com/go-ini/ini"

var (
	cfg *Cfg
)

var success = map[string]int{"code": 1}
var failure = map[string]int{"code": 0}

type Consumer struct {
	way Way
}

func (consumer *Consumer) setWay(way Way) {
	consumer.way = way
}

func (consumer *Consumer) Send(message Message) interface{} {
	return consumer.way.ConsumeMsg(message)
}

type Way interface {
	ConsumeMsg(message Message) interface{}
}

type Cfg struct {
	cfgFile       *ini.File
	HttpPort      string
	SsrConfigFile string
}

type Message struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Original string `json:"original"`
	MailName string `json:"mailName"`
	Way      string `json:"way"`
}
