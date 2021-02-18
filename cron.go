package main

import (
	"github.com/golang/glog"
	"github.com/robfig/cron/v3"
	"os/exec"
	"time"
)

type ProcInfo struct {
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

func NewMcron(procInfos map[string]*ProcInfo, jobs map[string]map[string]string) Mcron {
	var mc Mcron
	mc.cronEngine = cron.New()
	mc.cronEngine.Start()
	mc.cmds = procInfos
	mc.jobs = jobs
	mc.AddStageProcToMcron()
	mc.AddJobToMcron(jobs)
	return mc
}

func (mc *Mcron) AddJobToMcron(jobs map[string]map[string]string) (err error) {
	return nil
}

func (mc *Mcron) AddStageProcToMcron() (err error) {
	for name, proc := range mc.cmds {
		p := proc
		var cb func()
		if p.Name == "xxx" {
			cb = func() {
				time := time.Now().Format("20060102")
				logPath := "/home/" + p.Name + "/" + p.Name + "." + time + ".log"
				cmd := exec.Command("/bin/bash", "-c", "tail -5 "+logPath)
				out, err := cmd.Output()
				if err != nil {
					glog.Error(err)
				}
				cfg.TgBot.SendToMeStr(string(out))
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
