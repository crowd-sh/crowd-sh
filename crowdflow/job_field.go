package crowdflow

type InputField string
type OutputField string

type JobFieldType string

const (
	LongTextType JobFieldType = "long_text"
	ImageType    JobFieldType = "image"
	HiddenType   JobFieldType = "hidden"
	CheckBoxType JobFieldType = "checkbox"
)

type JobField struct {
	Id          string
	Value       string
	Description string
	Type        JobFieldType
}
