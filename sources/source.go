package sources

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

type SourceType string

const (
	CSVSourceType SourceType = "csv"
)

type Source interface {
	Init()
	Headers() []string
	Records() []map[string]string
}

type SourceConfig struct {
	Config map[string]string
	Type   SourceType

	source Source
}

func (s *SourceConfig) Init() {
	switch s.Type {
	case CSVSourceType:
		log.Println("Config CSV Input")
		s.source = &CSVSource{Config: s.Config}
	}

	s.source.Init()
}

func (s *SourceConfig) Headers() []string {
	return s.source.Headers()
}

func (s *SourceConfig) Records() []map[string]string {
	return s.source.Records()
}

func (s *SourceConfig) WriteAll(headers []string, rows []map[string]string) {
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
