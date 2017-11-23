// Copyright © 2017 Abhi Yerra <abhi@opszero.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
