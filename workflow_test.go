package crowdflow

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	//	"time"
)

var ProgramJson = `
{
    "title": "Tag the appropriate images",
    "description": "This is a Desc",
    "price": 1,
    "fields" : [
        {
            "id": "ImageUrl",
            "type": "InputField",
            "description": "Use this image to fill the information below.",
            "field_type": "image"
        },
        {
            "id": "Tags",
            "type": "OutputField",
            "description": "List all the relevent tags separated by a comma for the image. Ex. trees, castle, person"
        },
        {
            "id": "TextInImage",
            "type": "OutputField",
            "description": "Put any caption that appears of the image here. Put one item per line if there are multiple.",
            "field_type" :"long_text"
        },
        {
            "id": "IsCorrectOrientation",
            "type": "OutputField",
            "description": "Is the image in the correct orientation?",
            "field_type": "checkbox"
        },
        {
            "id": "IsLandscape",
            "type": "OutputField",
            "description": "Is the image of a landscape (a non urban setting)?",
            "field_type": "checkbox"
        },
        {
            "id": "IsPattern",
            "type": "OutputField",
            "description": "Is the image of a pattern?",
            "field_type": "checkbox"
        },
        {
            "id": "IsPerson",
            "type": "OutputField",
            "description": "Does the image contain people?",
            "field_type": "checkbox"
        },
        {
            "id": "TraditionalClothing",
            "type": "OutputField",
            "description": "If the image has people are they wearing traditional clothes?",
            "field_type": "checkbox"
        },
        {
            "id": "IsMap",
            "type": "OutputField",
            "description": "Is the image a map?",
            "field_type": "checkbox"
        },
        {
            "id": "IsDiagram",
            "type": "OutputField",
            "description": "Is the image a diagram?",
            "field_type": "checkbox"
        }
    ]
}`

var _ = Describe("NewWorkflow", func() {
	var (
		W Workflow
	)

	BeforeEach(func() {
		W = NewWorkflow(ProgramJson)
	})

	It("Parses the JSON correctly", func() {
		Expect(W.ProgramJson).To(Equal(ProgramJson))
		Expect(W.Program.Title).To(Equal("Tag the appropriate images"))
		Expect(W.Program.Description).To(Equal("This is a Desc"))
		Expect(W.Program.Price).To(Equal(1))
	})
})

var _ = PDescribe("Workflow", func() {

	BeforeEach(func() {
	})

	PDescribe("IsValidTask", func() {

	})

	PDescribe("AddTask", func() {

	})
})
