package main

import (
	//	"encoding/json"
	"log"
)

type Workflow struct {
	Id           int64
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	Tags         string        `json:"tags"`
	Url          string        `json:"url"`
	InputFields  []InputField  `json:"input_fields"`
	OutputFields []OutputField `json:"output_fields"`
	rawData      string        `json:"-"`
}

func (w Workflow) AddTask(taskJson string) {
	// Find all input fields.
	// For each input field create an assignment

	jobData := csv.NewReader(csvFile)

	records, err := jobData.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range records {
		newJob := Job{}

		// Assign the input fields
		for rowNum, rowFieldId := range p.CsvConfig {
			field := p.Fields.Get(rowFieldId)
			if field == nil {
				log.Fatal("Can't load field", rowFieldId)
			}

			i, _ := strconv.Atoi(rowNum)

			log.Println(row[i])
			newJob.InputFields = append(newJob.InputFields, JobField{
				Field: field,
				Value: row[i],
			})

		}

		// Assign the OutputFields
		for _, f := range p.Fields.OutputFields() {
			newJob.OutputFields = append(newJob.OutputFields, JobField{
				Field: f,
			})
		}

		p.Jobs = append(p.Jobs, newJob)
	}

	task.PublishToMTurk()
}

func (w Workflow) PollTasks() {
	// Find all input fields.
	// For each input field create an assignment
}

func (w Workflow) IsFinished() {
	var assignments []Assignment
	Db.Model(Assignment{}).Where("work_id = ?", w.Id).Find(&assignments)
	// See if all the assignments are done
}

func (w Workflow) Parse() {
	log.Println("Parsing Program")
	dec := json.NewDecoder(strings.NewReader(w.ProgramJson))

	err := dec.Decode(&w.Program)
	if err != nil {
		log.Fatal(err)
	}
}

func NewWorkflow(workflow string) (w *Workflow) {
	log.Println("Creating a new workflow")

	w.rawData = workflow
	w.Parse()

	Db.Create(&workflow)

	log.Println("Done creating workflow")

	return
}
