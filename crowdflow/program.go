package crowdflow

import (
	"encoding/json"
	"io"
	"log"
)

type Field struct {
	Id          string
	Type        string
	Description string
	FieldType   string
}

func (f *Field) GetInputFields() (f []Field) {

}

func (f *Field) GetOutputFields() (f []Field) {

}

type Program struct {
	Title       string
	Description string
	Fields      []Field
	Tags        string
	Writer      string // How to write the file
}

func (b *Program) TaskWriter() {
	switch b.Writer {
	case "csv":
		return
	}
}

func (b *Program) Run() {
	assignDone := make(chan bool, len(b.MetaJobs))

	for j := range b.MetaJobs {
		go NewAssignment(assignDone, b, &b.MetaJobs[j])
	}

	for i := 0; i < len(b.MetaJobs); i++ {
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
