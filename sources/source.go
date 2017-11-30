package sources

import (
	"log"
)

type SourceType string

const (
	CSVSourceType  SourceType = "csv"
	ExecSourceType            = "exec"
)

type Source interface {
	Init()
	Headers() []string
	Records() []map[string]string
	WriteAll(headers []string, rows []map[string]string)
}

type SourceConfig struct {
	Config map[string]string
	Type   SourceType

	source Source
}

func (s *SourceConfig) setSourceType() {
	switch s.Type {
	case CSVSourceType:
		log.Println("Config CSV")
		s.source = &CSVSource{Config: s.Config}
	case ExecSourceType:
		log.Println("Config Exec")
		s.source = &ExecSource{Config: s.Config}
	}
}

func (s *SourceConfig) Init() {
	s.setSourceType()
	s.source.Init()
}

func (s *SourceConfig) Headers() []string {
	return s.source.Headers()
}

func (s *SourceConfig) Records() []map[string]string {
	return s.source.Records()
}

func (s *SourceConfig) WriteAll(headers []string, rows []map[string]string) {
	s.setSourceType()

	log.Println(rows)
	s.source.WriteAll(headers, rows)
}
