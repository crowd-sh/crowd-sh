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
	"bytes"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/mturk"

	csvmap "github.com/recursionpharma/go-csv-map"
)

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

func (w *Workflow) Config() {
	file, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}

	err = json.Unmarshal(file, w)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	endpoint := SandboxEndpoint
	if len(os.Args) > 2 && os.Args[2] == "live" {
		endpoint = LiveEndpoint
	}

	fmt.Println("Endpoint:", endpoint)

	sess := session.Must(session.NewSession())
	w.client = mturk.New(sess, &aws.Config{
		Credentials: credentials.NewSharedCredentials("/Users/abhi/.aws/credentials", "opszero_mturk"),
		Endpoint:    aws.String(endpoint),
		Region:      aws.String("us-east-1"),
	})

	resp, err := w.client.GetAccountBalance(nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
}

func (w *Workflow) Save() {
	b, err := json.MarshalIndent(w, "", "    ")
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile(os.Args[1], b, os.ModePerm)
}

func (w *Workflow) BuildHit() {
	if w.MTurk.HitTypeId != "" {
		// Update HIT here
		return
	}

	resp, err := w.client.CreateHITType(&mturk.CreateHITTypeInput{
		AssignmentDurationInSeconds: aws.Int64(3000),
		AutoApprovalDelayInSeconds:  aws.Int64(3000),
		Title:       aws.String(w.Title),
		Description: aws.String(w.Description),
		Keywords:    aws.String(w.Tags),
		Reward:      aws.String(w.Reward),
	})

	if err != nil {
		fmt.Println(err)
	}

	w.MTurk.HitTypeId = *resp.HITTypeId
}

func (w *Workflow) newTask(records []map[string]string, i int, t *Task) {
	for _, workflowField := range w.Fields {
		workflowField.Value = records[i][workflowField.Name]
		fmt.Println(records[i])
		fmt.Println(records[i][workflowField.Name])
		t.Fields = append(t.Fields, workflowField)
	}

	resp, err := w.client.CreateHITWithHITType(&mturk.CreateHITWithHITTypeInput{
		HITTypeId:         aws.String(w.MTurk.HitTypeId),
		MaxAssignments:    aws.Int64(1),
		Question:          aws.String(t.Question()),
		LifetimeInSeconds: aws.Int64(86400 * 5),
	})

	if err == nil {
		t.HitID = *resp.HIT.HITId
		w.Tasks = append(w.Tasks, *t)
	} else {
		fmt.Println(err)
		if r := recover(); r != nil {
			fmt.Println("Recovered", r)
		}
	}
}

func (w *Workflow) BuildTasks() {
	file, err := ioutil.ReadFile(w.InputFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	reader := csvmap.NewReader(bytes.NewReader(file))
	reader.Columns, err = reader.ReadHeader()
	if err != nil {
		fmt.Println(err)
	}

	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	for i := range records {
		newTask := true

		t := &Task{Title: w.Title, Description: w.Description}
		t.SourceID = records[i]["ID"]

		for et := range w.Tasks {
			if w.Tasks[et].SourceID == t.SourceID {
				t = &w.Tasks[et]
				newTask = false
			}
		}

		if newTask {
			w.newTask(records, i, t)
		} else {
			// UpdateHITTypeOfHIT
			for field := range t.Fields {
				f := &t.Fields[field]
				f.Value = records[i][f.Name]
			}

			resp, err := w.client.ListAssignmentsForHIT(&mturk.ListAssignmentsForHITInput{
				HITId: aws.String(t.HitID),
			})

			fmt.Println(err)
			fmt.Println(resp)

			t.MTurk.Assignments = resp.Assignments

			if len(resp.Assignments) > 0 {
				xml.Unmarshal([]byte(*resp.Assignments[0].Answer), &t.MTurk.QuestionFormAnswers)

				for field := range t.Fields {
					f := &t.Fields[field]

					for _, answer := range t.MTurk.QuestionFormAnswers.Answer {
						if f.Name == answer.QuestionIdentifier {
							f.Value = strings.TrimSpace(answer.FreeText)
						}
					}
				}

			}
		}

		w.Save()
	}
}

func (w *Workflow) SaveOutput() {
	file, err := os.Create(w.OutputFile)
	fmt.Println(err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	var header []string
	header = []string{"ID"}

	for i := range w.Fields {
		header = append(header, w.Fields[i].Name)
	}
	writer.Write(header)

	for _, task := range w.Tasks {
		fields := []string{task.SourceID}

		for _, field := range task.Fields {
			fields = append(fields, field.Value)
		}

		writer.Write(fields)
	}
}
