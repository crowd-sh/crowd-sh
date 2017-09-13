package workflow

import (
	"io"
	"io/ioutil"
	"log"
)

const (
	CSVSource = "CSV"
)

type Source struct {
	Type     string
	Filename string
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

func (w *Workflow) Upload() {
	// Find all input fields.
	// For each input field create an assignment
}

func (w *Workflow) Gather() {
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
