package crowdflow

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux" // TODO: Deprecate this
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type JobField struct {
	Field *Field
	Value string `json:"value"`
}

type Job struct {
	InputFields  []JobField
	OutputFields []JobField
}

type Program struct {
	Title       string
	Description string
	Tags        string
	Url         string

	Fields    Fields
	Jobs      []Job
	CsvConfig map[string]string `json:"csv_config"`

	Writer string // How to write the file
}

var (
	AvailableAssignments Assignments
)

func (p *Program) OutputJobs() (fs []*Field) {
	for i, field := range p.Fields {
		if field.IsOutputField() {
			fs = append(fs, &p.Fields[i])
		}
	}

	return
}

func (p *Program) TaskWriter() {
	switch p.Writer {
	case "csv":
		return
	}
}

func (p *Program) Run() {
	assignDone := make(chan bool)

	// Create the assignments
	log.Println("Creating assignments of jobs", len(p.Jobs))

	for i := range p.Jobs {
		for j := range p.Jobs[i].OutputFields {
			go NewAssignment(assignDone, &p.Jobs[i], &p.Jobs[i].OutputFields[j])
		}
	}

	// Start Serving
	log.Println("Creating http server")
	go p.StartHttpServer()

	log.Println("Pinging index server")
	go p.PingIndexServer()

	// Wait for responses
	log.Println("Created assignments")
	for _ = range assignDone {
		log.Println("Waiting for assignments:", len(AvailableAssignments))
	}
}

func (p *Program) LoadJobs(csvFile io.Reader) {
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
}

func (p *Program) PingIndexServer() {
	v := url.Values{}
	v.Set("title", p.Title)
	v.Set("num_jobs", string(len(AvailableAssignments)))
	v.Set("url", p.Url)

	http.PostForm("http://workmachine.us/v1/tasks", v)
}

func (p *Program) StartHttpServer() {
	log.Println("Serving HTTP Server")

	r := mux.NewRouter()

	r.HandleFunc("/v1/assignment", func(w http.ResponseWriter, req *http.Request) {
		assign := AvailableAssignments.GetUnfinished()
		if assign == nil {
			renderJson(w, false)
			return
		}

		if !assign.TryToAssign() {
			renderJson(w, false)
			return
		}

		renderJson(w, assign)
	}).Methods("GET")

	r.HandleFunc("/v1/assignment", func(w http.ResponseWriter, req *http.Request) {
		log.Println("Posting", req.FormValue("id"))

		assign := AvailableAssignments.Find(req.FormValue("id"))
		if assign != nil {
			assign.Finish(req.FormValue(assign.InputField.Value))
		}

		renderJson(w, true)
	}).Methods("POST")
	http.Handle("/", r)
	http.ListenAndServe(fmt.Sprintf(":%d", 3001), nil)
}

func ParseProgram(program io.Reader) (m Program) {
	log.Println("Parsing Program")
	dec := json.NewDecoder(program)

	err := dec.Decode(&m)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Number of fields:", len(m.Fields))
	log.Println("CSV config:", len(m.CsvConfig))
	log.Println("Finished Parsing Program")

	return
}
