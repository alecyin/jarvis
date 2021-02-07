package main

var (
	cfg *Cfg
)

var success = map[string]int{"code": 1}
var failure = map[string]int{"code": 0}

type Consumer interface {
	ConsumeMsg()
}

type Cfg struct {
	HttpPort string
	Scs      []Sc
	QqMails  []QqMail
}

type Sc struct {
	Sckey string
}

type QqMail struct {
	FromAccount string
	AuthCode    string
	ToAccount   string
}

type Message struct {
	Title    string
	Content  string
	Original string
	MailName string
}
