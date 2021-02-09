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
	EntryID  cron.EntryID
}

type Mcron struct {
	stage      map[string]*ProcInfo // All the running proc
	cronEngine *cron.Cron
}

func NewMcron(procInfos map[string]*ProcInfo) Mcron {
	var mc Mcron
	mc.cronEngine = cron.New()
	mc.cronEngine.Start()
	mc.stage = procInfos
	mc.AddStageProcToMcron()
	return mc
}

func (mc *Mcron) AddStageProcToMcron() (err error) {
	for name, proc := range mc.stage {
		// add proc to cron.
		glog.Info("adding ", name, " to mcron... ")
		if err := mc.addToCron(proc); err != nil {
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
	mc.stage = make(map[string]*ProcInfo)
}

func (mc *Mcron) StartProc(name string) (err error) {
	Cmd := exec.Command(mc.stage[name].Run, mc.stage[name].Args...)
	glog.Info("now starting ", name)
	if err = Cmd.Start(); err != nil {
		glog.Error(err)
		return err
	}
	Cmd.Wait()
	glog.Info(name, " stoped and has runned. ")
	return nil
}

func (mc *Mcron) removeFromCron(name string) {
	mc.cronEngine.Remove(mc.stage[name].EntryID)
	delete(mc.stage, name)
	glog.Info(name, " has been removed from stage and cron engine")
}

func (mc *Mcron) addToCron(proc *ProcInfo) (err error) {
	entryId, err := mc.cronEngine.AddFunc(proc.Schedule, func() { mc.StartProc(proc.Name) })
	if err != nil {
		return err
	}
	proc.EntryID = entryId
	return nil
}
