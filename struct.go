package main

var (
	cfg *Cfg
)

type Cfg struct {
	HttpPort string
	Scs      []Sc
}

type Sc struct {
	Sckey string
}

type Message struct {
	Title   string
	Content string
}
