package crowdflow

type FieldType string
type InputType string

const (
	InputFieldType  FieldType = "input"
	OutputFieldType FieldType = "output"
)

const (
	LongTextType InputType = "long_text"
	ImageType    InputType = "image"
	HiddenType   InputType = "hidden"
	CheckBoxType InputType = "checkbox"
)

type Field struct {
	Id          string
	Type        FieldType
	InputType   InputType
	Value       string
	Description string
}
