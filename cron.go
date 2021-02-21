package main

import (
	"encoding/json"
	"github.com/golang/glog"
	"github.com/robfig/cron/v3"
	"os/exec"
	"sync"
	"time"
)

type ProcInfo struct {
	Nickname string   `json:"nickname"`
	Name     string   `json:"name"`
	Schedule string   `json:"schedule"`
	Run      string   `json:"run"`
	Args     []string `json:"args"`
	RunFun   func()
	EntryID  cron.EntryID
}

type Mcron struct {
	cmds       map[string]*ProcInfo // All the proc from file
	jobs       map[string]map[string]string
	cronEngine *cron.Cron
}

var mc *Mcron
var mcOnce sync.Once

func GetMcronIns() *Mcron {
	mcOnce.Do(func() {
		mc = newMcron()
	})
	return mc
}

func newMcron() *Mcron {
	cronCmdFileContent, err := ReadTotalFile(cronCmdFilePath)
	if err != nil {
		glog.Error("load cron config error,", err)
	}
	var procInfos []ProcInfo
	if err = json.Unmarshal(cronCmdFileContent, &procInfos); err != nil {
		glog.Error("convert cron json array to struct error", err)
	}
	u := map[string]*ProcInfo{}
	for _, procInfo := range procInfos {
		p := procInfo
		u[procInfo.Name] = &p
	}
	cronJobFileContent, err := ReadTotalFile(cronJobFilePath)
	if err != nil {
		glog.Error("load cron config error,", err)
	}
	var jobs map[string]map[string]string
	if err = json.Unmarshal(cronJobFileContent, &jobs); err != nil {
		glog.Error(err)
	}
	mc := new(Mcron)
	mc.cronEngine = cron.New()
	mc.cmds = u
	mc.jobs = jobs
	return mc
}

func (mc *Mcron) Run() {
	mc.AddStageProcToMcron()
	mc.AddJobToMcron()
	mc.cronEngine.Start()
}

func (mc *Mcron) AddJobToMcron() (err error) {
	return nil
}

func (mc *Mcron) AddStageProcToMcron() (err error) {
	for name, proc := range mc.cmds {
		p := proc
		var cb func()
		if p.Name == "wyy" || p.Name == "smzdm" || p.Name == "mistepupdate" {
			cb = func() {
				time := time.Now().Format("20060102")
				logPath := "/home/" + p.Name + "/" + p.Name + "." + time + ".log"
				cmd := exec.Command("/bin/bash", "-c", "tail -5 "+logPath)
				out, err := cmd.Output()
				if err != nil {
					glog.Error(err)
				}
				GetTgBotIns().SendToMeStr(p.Nickname + "\n" + string(out))
			}
		}
		p.RunFun = func() { mc.StartProc(p.Name, cb) }
		// add proc to cron.
		glog.Info("adding ", name, " to mcron... ")
		if err := mc.addCmdToCron(p); err != nil {
			glog.Error("adding ", name, " to mcron error:", err)
			continue
		}
		glog.Info(name, " has been added to mcron.")
	}
	return
}

func (mc *Mcron) RefreshMcron() {
	glog.Info("now clearing mcron")
	//stop cron clock
	mc.cronEngine.Stop()
	glog.Info(" mcron stopped! ")
	mc.cronEngine = cron.New()
	mc.cronEngine.Start()
	glog.Info(" mcron started! ")
	time.Sleep(2 * time.Second)
	mc.cmds = make(map[string]*ProcInfo)
}

func (mc *Mcron) StartProc(name string, c func()) (err error) {
	Cmd := exec.Command(mc.cmds[name].Run, mc.cmds[name].Args...)
	glog.Info("now starting ", name)
	if err = Cmd.Start(); err != nil {
		glog.Error("run ", name, " err:", err)
		return err
	}
	Cmd.Wait()
	glog.Info(name, " has stopped and has runned. ")
	if c != nil {
		c()
	}
	return nil
}

func (mc *Mcron) removeFromCron(name string) {
	mc.cronEngine.Remove(mc.cmds[name].EntryID)
	delete(mc.cmds, name)
	glog.Info(name, " has been removed from cmds and cron engine")
}

func (mc *Mcron) addCmdToCron(proc *ProcInfo) (err error) {
	entryId, err := mc.cronEngine.AddFunc(proc.Schedule, proc.RunFun)
	if err != nil {
		return err
	}
	proc.EntryID = entryId
	return nil
}
