package main

import (
	"time"

	flag "github.com/spf13/pflag"
)

var (
	isLive bool
)

const (
	SandboxEndpoint = "https://mturk-requester-sandbox.us-east-1.amazonaws.com"
	LiveEndpoint    = "https://mturk-requester.us-east-1.amazonaws.com"
)

func init() {
	flag.BoolVar(&isLive, "live", false, "Live on Mechanical Turk.")
	flag.Parse()
}

func main() {
	w := &Workflow{}
	w.Config()
	w.BuildHit()
	w.Save()
	go w.SaveOutput()
	w.BuildTasks()

	time.Sleep(8 * time.Second)
}
