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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/mturk"
)

type TaskStatus string

const (
	FinishedStatus TaskStatus = "finished"
)

type Field struct {
	Name        string
	Type        string
	Description string
	Value       string
}

type Task struct {
	HitID    string
	SourceID string
	Fields   []Field
	Status   TaskStatus
}

type Workflow struct {
	Title       string
	Description string
	Tags        string
	Reward      string

	InputFile  string
	OutputFile string

	Fields []Field

	Tasks []Task
	MTurk struct {
		HitTypeId string
	}

	client *mturk.MTurk
}

const (
	SandboxEndpoint = "https://mturk-requester-sandbox.us-east-1.amazonaws.com"
)

func (w *Workflow) Config() {
	file, e := ioutil.ReadFile(os.Args[1])
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	fmt.Printf("%s\n", string(file))

	json.Unmarshal(file, w)

	sess := session.Must(session.NewSession())
	w.client = mturk.New(sess, &aws.Config{
		Credentials: credentials.NewSharedCredentials("/home/abhi/.aws/credentials", "opszero_mturk"),
		Endpoint:    aws.String(SandboxEndpoint),
		Region:      aws.String("us-east-1"),
	})

	x, err := w.client.GetAccountBalance(nil)
	fmt.Println(err)
	fmt.Println(x)
}

func (w *Workflow) Save() {
	b, err := json.MarshalIndent(w, "", "    ")
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile(os.Args[1], b, os.ModePerm)
}

func (w *Workflow) BuildHit() {
	a := aws.String(w.Reward)
	fmt.Println(*a)

	resp, err := w.client.CreateHITType(&mturk.CreateHITTypeInput{
		AssignmentDurationInSeconds: aws.Int64(300),
		AutoApprovalDelayInSeconds:  aws.Int64(3000),
		Title:       aws.String(w.Title),
		Description: aws.String(w.Description),
		Keywords:    aws.String(w.Tags),
		Reward:      aws.String(w.Reward),
	})

	fmt.Println(err)
	fmt.Println(resp)

	w.MTurk.HitTypeId = *resp.HITTypeId
}

func main() {
	w := &Workflow{}
	w.Config()
	w.BuildHit()
	w.Save()

}
