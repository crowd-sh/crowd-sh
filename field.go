package main

type FieldType string
type InputType string

const (
	InputFieldType  FieldType = "InputField"
	OutputFieldType FieldType = "OutputField"
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

func (f *Field) IsOutputField() bool {
	return f.Type == OutputFieldType
}

func (f *Field) IsInputField() bool {
	return f.Type == InputFieldType
}

type Fields []Field

func (fs Fields) Get(id string) *Field {
	for i := range fs {
		if fs[i].Id == id {
			return &fs[i]
		}
	}

	return nil
}

func (fs Fields) OutputFields() (ofs []*Field) {
	for i, f := range fs {
		if f.IsOutputField() {
			ofs = append(ofs, &fs[i])
		}
	}

	return
}
