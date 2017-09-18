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
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/recursionpharma/go-csv-map"

	"encoding/csv"
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

func (t *Field) HTML() string {
	return fmt.Sprintf(`
  <h2>%s</h2>
  <p>%s</p>
  <p><textarea name='%s' cols='80' rows='3'>%s</textarea></p>`, t.Name, t.Description, t.Name, t.Value)
}

type Task struct {
	// Copied from Workflow
	Title       string
	Description string

	HitID    string
	SourceID string
	Fields   []Field
	Status   TaskStatus

	MTurk struct {
		QuestionFormAnswers struct {
			Answer []struct {
				QuestionIdentifier string
				FreeText           string
			}
		}

		Assignments []*mturk.Assignment
	}
}

func (t *Task) Question() string {
	var fieldsHTML string
	for i := range t.Fields {
		fieldsHTML += t.Fields[i].HTML()
	}

	return fmt.Sprintf(`
<HTMLQuestion xmlns="http://mechanicalturk.amazonaws.com/AWSMechanicalTurkDataSchemas/2011-11-11/HTMLQuestion.xsd">
  <HTMLContent><![CDATA[
<!DOCTYPE html>
<html>
 <head>
  <meta http-equiv='Content-Type' content='text/html; charset=UTF-8'/>
  <script type='text/javascript' src='https://s3.amazonaws.com/mturk-public/externalHIT_v1.js'></script>
 </head>
 <body>
  <form name='mturk_form' method='post' id='mturk_form' action='https://www.mturk.com/mturk/externalSubmit'>
<h1>%s</h1>

<p>
%s
</p>

%s

  <p>
    <input type='hidden' value='' name='assignmentId' id='assignmentId'/>
    <input type='submit' id='submitButton' value='Submit' /> 
  </p>
  </form>

  <script language='Javascript'>turkSetAssignmentID();</script>
 </body>
</html>
]]>
  </HTMLContent>
  <FrameHeight>450</FrameHeight>
</HTMLQuestion>
`, t.Title, t.Description, fieldsHTML)
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

	if w.MTurk.HitTypeId != "" {
		// Update Lead here
		return
	}

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

func (w *Workflow) BuildTasks() {
	file, e := ioutil.ReadFile(w.InputFile)
	fmt.Println(e)
	fmt.Println(file)

	var err error

	reader := csvmap.NewReader(bytes.NewReader(file))
	reader.Columns, err = reader.ReadHeader()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(reader.Columns)

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
				LifetimeInSeconds: aws.Int64(86400),
			})

			fmt.Println(err)

			t.HitID = *resp.HIT.HITId

			w.Tasks = append(w.Tasks, *t)
		} else {
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
							f.Value = answer.FreeText
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
	for i := range w.Fields {
		header = append(header, w.Fields[i].Name)
	}
	writer.Write(header)

	for i := range w.Tasks {
		var fields []string

		for _, field := range w.Tasks[i].Fields {
			fields = append(fields, field.Value)
		}

		writer.Write(fields)
	}
}

func main() {
	w := &Workflow{}
	w.Config()
	w.BuildHit()
	w.Save()
	w.BuildTasks()
	w.SaveOutput()
}
