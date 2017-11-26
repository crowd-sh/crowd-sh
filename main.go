package main

import (
	flag "github.com/spf13/pflag"
)

var (
	isLive     bool
	isReview   bool
	configFile string
)

const (
	SandboxEndpoint = "https://mturk-requester-sandbox.us-east-1.amazonaws.com"
	LiveEndpoint    = "https://mturk-requester.us-east-1.amazonaws.com"
)

func init() {
	flag.BoolVar(&isLive, "live", false, "Live on Mechanical Turk.")
	flag.BoolVar(&isReview, "review", false, "Review Work.")
	flag.Parse()
}

func main() {
	configFile = flag.Arg(0)

	w := &Workflow{}
	w.Config()

	if isReview {
		w.ReviewTasks()
	} else {
		w.BuildHit()
		w.Save()
		w.BuildTasks()
		w.SaveOutput()
	}
}
