package main

import (
	"encoding/json"
	"github.com/golang/glog"
	"io"
	"os"
)

func prepare() (*os.File, error) {
	if _, err := os.Stat(recordFilePath); os.IsNotExist(err) {
		os.Create(recordFilePath)
	}
	return os.OpenFile(recordFilePath, os.O_APPEND, 0666)
}

func RecordMessage(message Message) {
	f, err := prepare()
	if err != nil {
		glog.Error("record file open error, %v", err)
		return
	}
	defer f.Close()
	msg, err := json.Marshal(message)
	if err != nil {
		glog.Error("convert json error, %v", err)
		return
	}
	if _, err = io.WriteString(f, string(msg)+"\n"); err != nil {
		glog.Error("write record file error, %v", err)
		return
	}
	glog.Info("write to record file success")
}
