package main

import (
	flag "github.com/spf13/pflag"
)

var (
	isLive     bool
	isReview   bool
	isExpire   bool
	output     bool
	configFile string
)

const (
	SandboxEndpoint = "https://mturk-requester-sandbox.us-east-1.amazonaws.com"
	LiveEndpoint    = "https://mturk-requester.us-east-1.amazonaws.com"
)

func init() {
	flag.BoolVar(&isLive, "live", false, "Live on Mechanical Turk.")
	flag.BoolVar(&isExpire, "expire", false, "Expire Work.")
	flag.BoolVar(&output, "output", false, "Output Work.")
	flag.Parse()
}

func main() {
	configFile = flag.Arg(0)

	w := &Workflow{}
	w.Config()

	if isExpire {
		//w.ExpireTasks()
	} else {
		w.BuildHitType()
		w.Save()
		w.Sync()
		w.Save()
	}
}
