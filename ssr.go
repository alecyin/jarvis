package main

import (
	"encoding/json"
	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type Ssr struct {
	Server              string      `json:"server"`
	ServerIpv6          string      `json:"server_ipv6"`
	ServerPort          int         `json:"server_port"`
	LocalAddress        string      `json:"local_address"`
	LocalPort           int         `json:"local_port"`
	Password            string      `json:"password"`
	Method              string      `json:"method"`
	Protocol            string      `json:"protocol"`
	ProtocolParam       string      `json:"protocol_param"`
	Obfs                string      `json:"obfs"`
	ObfsParam           string      `json:"obfs_param"`
	SpeedLimitPerCon    int         `json:"speed_limit_per_con"`
	SpeedLimitPerUser   int         `json:"speed_limit_per_user"`
	AdditionalPorts     interface{} `json:"additional_ports"`
	AdditionalPortsOnly bool        `json:"additional_ports_only"`
	Timeout             int         `json:"timeout"`
	UdpTimeout          int         `json:"udp_timeout"`
	DnsIpv6             bool        `json:"dns_ipv6"`
	ConnectVerboseInfo  int         `json:"connect_verbose_info"`
	Redirect            string      `json:"redirect"`
	FastOpen            bool        `json:"fast_open"`
}

var ssr *Ssr
var sOnce sync.Once

func GetSsrIns() *Ssr {
	sOnce.Do(func() {
		ssr = newSsr()
	})
	return ssr
}

func newSsr() *Ssr {
	s := new(Ssr)
	content, err := ReadTotalFile(cfg.SsrConfigFile)
	if err != nil {
		glog.Error("read ssr config file error,", err)
		return nil
	}
	err = json.Unmarshal(content, &s)
	if err != nil {
		glog.Error("read ssr config file error,", err)
		return nil
	}
	return s
}

func (ssr *Ssr) parseNodeFile() map[string]interface{} {
	content, err := ReadTotalFile(ssrNodeFilePath)
	if err != nil {
		glog.Error("read ssr node file error, ", err)
		return nil
	}
	result := make(map[string]interface{})
	err = yaml.Unmarshal(content, &result)
	if err != nil {
		glog.Error("error: ", err)
		return nil
	}
	return result
}

func (ssr *Ssr) ChangeNode() bool {
	res := ssr.parseNodeFile()
	if res == nil {
		return false
	}
	nodes := res["proxies"].([]interface{})
	flag := false // is current node
	for i := 0; i < len(nodes); i++ {
		node := nodes[i].(map[interface{}]interface{})
		if strings.Contains(node["name"].(string), "VIP") {
			if flag {
				ssr.Server = node["server"].(string)
				ssr.ServerPort = node["port"].(int)
				ssr.Password = node["password"].(string)
				ssr.Method = node["cipher"].(string)
				ssr.Protocol = node["protocol"].(string)
				ssr.ProtocolParam = node["protocol-param"].(string)
				ssr.Obfs = node["obfs"].(string)
				ssr.ObfsParam = node["obfs-param"].(string)
				break
			}
			if node["server"].(string) == ssr.Server && node["port"].(int) == ssr.ServerPort {
				flag = true
			}
		}
		if i == len(nodes)-1 { // start over
			i = 0
		}
	}

	b, err := json.Marshal(ssr)
	if err != nil {
		glog.Error("Error: ", err)
		return false
	}
	if err = WriteCoverFile(cfg.SsrConfigFile, string(b)); err != nil {
		glog.Error(err)
		return false
	}
	_, err = exec.Command("/bin/bash", "-c", "/usr/local/bin/ssr", "restart").Output()
	if err != nil {
		glog.Error("exec command ssr restart error,", err)
		return false
	}
	glog.Info("successful call ssr restart")
	//output := string(out[:])
	//glog.Info(output)
	return true
}

func (ssr *Ssr) ServiceabilityTest() {
	proxyAddr := cfg.cfgFile.Section("tg_bot").Key("proxy_addr").String()
	testUrl := cfg.cfgFile.Section("ssr").Key("test_url").String()
	for {
		cmd := exec.Command("/bin/bash", "-c", "curl  -x "+proxyAddr+" --connect-timeout 2 --retry 3 -I "+testUrl)
		out, err := cmd.Output()
		if err == nil && strings.Contains(string(out), "HTTP/1.1 200 OK") {
			glog.Info("proxy is normal")
			time.Sleep(5 * time.Minute)
			continue
		}
		glog.Info("proxy is abnormal")
		if err != nil {
			glog.Error(err)
		}
		i := 3
		for !ssr.ChangeNode() && i > 0 {
			glog.Info("change proxy failure,retry")
			i--
		}
	}
}

//func (ssr *Ssr) addSsrJobsToCron() {
//	jobs := cfg.Mcron.jobs
//	if _, err := cfg.Mcron.cronEngine.AddFunc(jobs["SSR_ChangeNode"]["schedule"], pingGoogleTest); err != nil {
//		glog.Error("add job SSR_ChangeNode error:", err)
//	}
//}
