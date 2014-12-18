package main

const (
	LongTextType = "long_text"
	ImageType    = "image"
	HiddenType   = "hidden"
	CheckBoxType = "checkbox"
)

type InputField struct {
	Id          string `json:"id"`
	Type        string `json:"type"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

type OutputField struct {
	Id          string `json:"id"`
	Type        string `json:"type"`
	Value       string `json:"value"`
	Description string `json:"description"`
	Validation  string `json:"validation"`
}
