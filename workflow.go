package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
)

const (
	LongTextType = "long_text"
	ImageType    = "image"
	HiddenType   = "hidden"
	CheckBoxType = "checkbox"
)

type InputField struct {
	Id          string `json:"id"`
	Type        string `json:"field_type"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

type OutputField struct {
	Id          string `json:"id"`
	Type        string `json:"field_type"`
	Value       string `json:"value"`
	Description string `json:"description"`
	Validation  string `json:"validation"`
}

type Workflow struct {
	Id           int64         `json:"id"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	Tags         string        `json:"tags"`
	Url          string        `json:"url"`
	Price        int           `json:"price"`
	InputFields  []InputField  `json:"input_fields"`
	OutputFields []OutputField `json:"output_fields"`
	RawData      string        `json:"-"`
}

func (w *Workflow) AddTask(task io.ReadCloser) (t Task) {
	body, err := ioutil.ReadAll(task)
	if err != nil {
		log.Println(err)
	}

	t.RawData = string(body)
	t.Parse()
	if !t.VerifyWithWorkflow(w) {
		return Task{}
	}
	t.WorkflowId = w.Id

	Db.Create(&t)
	t.PublishToMTurk()

	return
}

func (w Workflow) PollTasks() {
	// Find all input fields.
	// For each input field create an assignment
}

// func (w Workflow) IsFinished() {
// 	var assignments []Assignment
// 	Db.Model(Assignment{}).Where("work_id = ?", w.Id).Find(&assignments)
// 	// See if all the assignments are done
// }

func (w *Workflow) Parse() {
	err := json.Unmarshal([]byte(w.RawData), w)
	if err != nil {
		log.Println(err)
	}
}

func (w *Workflow) IsValid() bool {
	// Make sure inputs aren't more than 10 fields

	return true
}

func NewWorkflow(workflow io.ReadCloser) (w Workflow) {
	body, err := ioutil.ReadAll(workflow)
	if err != nil {
		log.Println(err)
	}

	w.RawData = string(body)
	w.Parse()
	if !w.IsValid() {
		return Workflow{}
	}

	log.Println("Creating a new workflow", w.Title)
	log.Println(string(body))
	log.Println()

	Db.Create(&w)

	log.Println("Done creating workflow")

	return
}
