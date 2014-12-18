package main

import (
	"encoding/csv"
	//	"encoding/json"
	//	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Workflow struct {
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	Tags         string        `json:"tags"`
	Url          string        `json:"url"`
	InputFields  []InputField  `json:"input_fields"`
	OutputFields []OutputField `json:"output_fields"`
}

type Work struct {
	Id          int64
	WorkflowRaw string `sql:"workflow"`
	DataRaw     string `sql:"data"`
	Notes       string
}

func (w Work) CreateAssignments() {
	// Find all input fields.
	// For each input field create an assignment
}

func (w Work) IsFinished() {
	var assignments []Assignment
	Db.Model(Assignment{}).Where("work_id = ?", w.Id).Find(&assignments)
	// See if all the assignments are done
}

// func NewWorkflow(program string) (m Workflow) {
// 	m = Workflow{ProgramJson: program}
// 	m.Parse()

// 	log.Println("Number of fields:", len(m.Program.Fields))
// 	log.Println("Finished Parsing Program")

// 	return
// }

// func (w *Workflow) Parse() {
// 	log.Println("Parsing Program")
// 	dec := json.NewDecoder(strings.NewReader(w.ProgramJson))

// 	err := dec.Decode(&w.Program)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// func (w *Workflow) IsValidTask(t Task) {

// }

// func (w *Workflow) AddTask(json io.Reader) {
// 	jobData := csv.NewReader(csvFile)

// 	records, err := jobData.ReadAll()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	for _, row := range records {
// 		newJob := Job{}

// 		// Assign the input fields
// 		for rowNum, rowFieldId := range p.CsvConfig {
// 			field := p.Fields.Get(rowFieldId)
// 			if field == nil {
// 				log.Fatal("Can't load field", rowFieldId)
// 			}

// 			i, _ := strconv.Atoi(rowNum)

// 			log.Println(row[i])
// 			newJob.InputFields = append(newJob.InputFields, JobField{
// 				Field: field,
// 				Value: row[i],
// 			})

// 		}

// 		// Assign the OutputFields
// 		for _, f := range p.Fields.OutputFields() {
// 			newJob.OutputFields = append(newJob.OutputFields, JobField{
// 				Field: f,
// 			})
// 		}

// 		p.Jobs = append(p.Jobs, newJob)
// 	}
// }

// // // func (p *Program) OutputJobs() (fs []*Field) {
// // // 	for i, field := range p.Fields {
// // // 		if field.IsOutputField() {
// // // 			fs = append(fs, &p.Fields[i])
// // // 		}
// // // 	}

// // // 	return
// // // }

func NewFlow(workflow string, csv *csv.Reader) {
	// Read the CSV file

	rows, err := csv.Read()
	if err != nil {
		log.Println(err)
		return
	}

	for i := range rows {
		w := Work{
			WorkflowRaw: workflow,
			DataRaw:     rows[i],
		}
		Db.Create(&w)

		go w.CreateAssignments()
	}

}

func NewWorkflowsHandler(w http.ResponseWriter, r *http.Request) {
	//	vars := mux.Vars(r)

	// for key, _ := range r.Form {
	// 	log.Println(key)
	// 	//LOG: {"test": "that"}
	// 	err := json.Unmarshal([]byte(key), &t)
	// 	if err != nil {
	// 		log.Println(err.Error())
	// 	}
	// }

	// NewFlow()

	// workflow := NewWorkflow("Json")
	// Db.Create(&workflow)
}

func ShowWorkflowsHandler(w http.ResponseWriter, r *http.Request) {

}

// func newTaskHandler(w http.ResponseWriter, req *http.Request) {
// 	for _, name := range []string{"name", "num_jobs", "url"} {
// 		if req.FormValue(name) == "" {
// 			renderJson(w, fmt.Sprintf("error: Need value %s", name))
// 			return
// 		}
// 	}

// 	task := Task{
// 		Name:      req.FormValue("name"),
// 		NumJobs:   req.FormValue("num_jobs"),
// 		Url:       req.FormValue("url"),
// 		CreatedAt: time.Now(),
// 	}
// 	task.GenerateId()

// 	tasks = append(tasks, task)

// 	log.Println("New Task", task.Id, req.FormValue("name"), req.FormValue("num_jobs"), req.FormValue("url"))

// 	fmt.Fprintln(w, task.Id)
// }
