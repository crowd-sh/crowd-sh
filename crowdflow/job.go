package crowdflow

import (
	"fmt"
)

// Interface for serving crowdsourcing jobs.
type Assigner interface {
	// Publish a Batch to be done by workers.
	Execute(jobs chan Job, j Job)
}

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

type Job struct {
	Parent      *MetaJob  `json:"info"`
	OutputField *JobField `json:"output"`
}

type MetaJob struct {
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	InputFields  []JobField `json:"input_fields"`
	OutputFields []JobField `json:"-"`
	Jobs         []Job      `json:"-"`
}

func (mj *MetaJob) StartJobs(assigner Assigner, meta_jobs chan MetaJob, parent MetaJob) {
	jobs := make(chan Job)

	for i, _ := range parent.OutputFields {
		go assigner.Execute(jobs, Job{
			Parent:      &parent,
			OutputField: &parent.OutputFields[i],
		})
	}

	for _ = range parent.OutputFields {
		// Do a review on these jobs.

		_ = <-jobs
	}

	fmt.Printf("\n\nJob returned: %v\n", parent.OutputFields)

	meta_jobs <- parent
}
