package sources

import "log"

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
