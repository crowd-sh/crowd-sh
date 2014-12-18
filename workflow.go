package main

// import (
// 	// "encoding/csv"
// 	"encoding/json"
// 	"log"
// 	"strings"
// 	// "strconv"
// )

// type Workflow struct {
// 	Id          int64
// 	UserId      int64
// 	ProgramJson string
// 	Tasks       []Task

// 	Program struct {
// 		Title       string
// 		Description string
// 		Price       int
// 		Tags        string
// 		Url         string

// 		Fields Fields
// 	} `sql:"-"`
// }

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
