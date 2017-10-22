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
	"html"
	"io/ioutil"
	"os"

	"github.com/recursionpharma/go-csv-map"

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
<div class="row">
<div class="col-md-12">
  <h2>%s</h2>
  <p>%s</p>
  <p>
    <textarea name='%s' cols='80' rows='3'>
    %s
    </textarea>
  </p>
</div>
</div>`, html.EscapeString(t.Name), html.EscapeString(t.Description), t.Name, t.Value)
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

<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">
<script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js" integrity="sha384-Tc5IQib027qvyjSMfHjOMaLkfuWVxZxUPnCJA7l2mCWNIpG9mGCD8wGNIcPD7Txa" crossorigin="anonymous"></script>

 </head>
 <body>
  <div class="container">
    <form name='mturk_form' method='post' id='mturk_form' action='https://www.mturk.com/mturk/externalSubmit'>
    <h1>%s</h1>

    <p>
    %s
    </p>

    %s

    <p>
	<input type='hidden' value='' name='assignmentId' id='assignmentId'/>
	<input type='submit' id='submitButton' value='Submit' class='btn btn-success' /> 
    </p>
    </form>
  </div>

  <script language='Javascript'>turkSetAssignmentID();</script>
 </body>
</html>
]]>
  </HTMLContent>
  <FrameHeight>1000</FrameHeight>
</HTMLQuestion>
`, html.EscapeString(t.Title), html.EscapeString(t.Description), fieldsHTML)
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
	LiveEndpoint    = "https://mturk-requester.us-east-1.amazonaws.com"
)

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
		Credentials: credentials.NewSharedCredentials("/home/abhi/.aws/credentials", "opszero_mturk"),
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

func main() {
	w := &Workflow{}
	w.Config()
	w.BuildHit()
	w.Save()
	w.BuildTasks()
	w.SaveOutput()
}
