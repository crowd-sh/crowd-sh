package crowdflow

import (
	"encoding/json"
	"io"
	"log"
)

type Program struct {
	Title       string
	Description string
	Fields      []Field
	Tags        string
	Writer      string // How to write the file
}

func (p *Program) OutputJobs() (fs []*Field) {
	for i := range p.Fields {
		if p.Fields[i] == OutputFieldType {
			fd = append(fs, &p.Fields[i])
		}
	}
}

func (p *Program) TaskWriter() {
	switch p.Writer {
	case "csv":
		return
	}
}

func (p *Program) Run() {
	assignDone := make(chan bool, len(p.MetaJobs))

	for j := range p.MetaJobs {
		go NewAssignment(assignDone, p, &p.MetaJobs[j])
	}

	for i := 0; i < len(p.MetaJobs); i++ {
		<-assignDone
	}
}

func ParseProgram(program io.Reader) (m Program) {
	log.Println("Parsing Program")
	dec := json.NewDecoder(program)

	err := dec.Decode(&m)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Finished Parsing Program")

	return
}
