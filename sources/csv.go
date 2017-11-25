package sources

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	csvmap "github.com/recursionpharma/go-csv-map"
)

type CSVSource struct {
	Config map[string]string

	columns []string
	records []map[string]string
}

func (w *CSVSource) Init() {
	file, err := ioutil.ReadFile(w.Config["File"])
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

func (w *CSVSource) Headers() []string {
	return w.columns
}

func (w *CSVSource) Records() []map[string]string {
	return w.records
}

func (s *CSVSource) WriteAll(headers []string, rows []map[string]string) {
	log.Println("Writing to ", s.Config["File"])
	file, err := os.Create(s.Config["File"])
	fmt.Println(err)
	defer file.Close()

	var r [][]string

	writer := csv.NewWriter(file)
	defer writer.Flush()

	r = append(r, headers)

	for _, row := range rows {
		var r2 []string
		for _, header := range headers {
			r2 = append(r2, row[header])
		}

		r = append(r, r2)
	}

	fmt.Println(r)
	writer.WriteAll(r)
}
