package sources

import (
	"bytes"
	"fmt"
	"io/ioutil"

	csvmap "github.com/recursionpharma/go-csv-map"
)

type CSVSource struct {
	InputFile string

	columns []string
	records []map[string]string
}

func (w *CSVSource) Init() {
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

	w.columns = reader.Columns
	w.records = records
}

func (w *CSVSource) Columns() []string {
	return w.columns
}

func (w *CSVSource) Records() []map[string]string {
	return w.records
}
