package sources

type SourceType string

const (
	CSVSourceType SourceType = "csv"
)

type Source interface {
	Init()
	Columns() []string
	Records() []map[string]string
}
