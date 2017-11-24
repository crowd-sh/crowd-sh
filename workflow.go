package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/mturk"

	"github.com/opszero/workmachine/sources"
)

type Workflow struct {
	Title       string
	Description string
	Tags        string
	Reward      string

	Input  sources.SourceConfig
	Output sources.SourceConfig

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
	if isLive {
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

	w.Input.Init()
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
		AutoApprovalDelayInSeconds:  aws.Int64(86400),
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

func (w *Workflow) newTask(records []map[string]string, i int, t *Task, isRepeat bool) {
	if !isRepeat {
		for _, workflowField := range w.Fields {
			workflowField.Value = records[i][workflowField.Name]
			fmt.Println(records[i])
			fmt.Println(records[i][workflowField.Name])
			t.Fields = append(t.Fields, workflowField)
		}
	}

	resp, err := w.client.CreateHITWithHITType(&mturk.CreateHITWithHITTypeInput{
		HITTypeId:         aws.String(w.MTurk.HitTypeId),
		MaxAssignments:    aws.Int64(1),
		Question:          aws.String(t.Question()),
		LifetimeInSeconds: aws.Int64(86400), // 1 day
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

func (w *Workflow) updateTask(records []map[string]string, i int, t *Task) {
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

	allRejected := true

	for _, ass := range resp.Assignments {
		if *ass.AssignmentStatus != "Rejected" && allRejected {
			allRejected = false
		}
	}

	if allRejected {
		w.newTask(records, i, t, true)
	} else {
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
}

func (w *Workflow) BuildTasks() {
	records := w.Input.Records()

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
			w.newTask(records, i, t, false)
		} else {
			w.updateTask(records, i, t)
		}

		w.Save()
	}
}

func (w *Workflow) SaveOutput() {
	var (
		header []string
		rows   []map[string]string
	)

	header = []string{"ID"}

	for i := range w.Fields {
		header = append(header, w.Fields[i].Name)
	}

	for _, task := range w.Tasks {
		row := make(map[string]string)
		row["ID"] = task.SourceID

		for _, field := range task.Fields {
			row[field.Name] = field.Value
			rows = append(rows, row)
		}
	}

	w.Output.WriteAll(header, rows)

	select {
	case <-time.After(5 * time.Second):
		fmt.Println("Save")
	}
}
