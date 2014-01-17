package crowdflow

import (
	"fmt"
)

type MetaJob struct {
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	InputFields  []JobField `json:"input_fields"`
	OutputFields []JobField `json:"-"`
	Jobs         []Job      `json:"-"`
}

func (mj *MetaJob) StartSplitJob(assigner SplitAssigner, meta_jobs chan MetaJob, parent MetaJob) {
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

func (mj *MetaJob) StartJob(assigner Assigner, meta_jobs chan MetaJob, parent MetaJob) {
	job := make(chan MetaJob)

	go assigner.Execute(job, parent)

	_ = <-job

	meta_jobs <- parent
}
