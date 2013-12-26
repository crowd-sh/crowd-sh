package workmachine

import (
	"fmt"
	"reflect"
)

const (
	WorkId   = "work_id"
	WorkDesc = "work_desc"
	WorkType = "work_type"
)

// Interface for serving crowdsourcing jobs.
type Assigner interface {
	// Publish a Batch to be done by workers.
	Execute(jobs chan Job, j Job)
}

type InputField string
type OutputField string

/*
 * Build a Job from a Task.
 */

type JobField struct {
	Id          string
	Value       string
	Description string
	Type        string
}

type Job struct {
	InputFields  []JobField
	OutputFields []JobField
}

/*
 * Task is a way to define the Job that needs to be run.
 */
// work_type: output_image, output_text, input_text
type Task struct {
	Title       string
	Description string
	Price       int // In cents
	Write       func(j *Job)
	Tasks       interface{} // TODO: Name should be renamed Work or something like that.
}

/*
 * Batch is a group of Jobs
 */

type Batch struct {
	Task    Task
	Jobs    []Job
	Results []Job
}

func (b *Batch) Run(ss Assigner) {
	jobs := make(chan Job)

	for _, j := range b.Jobs {
		go ss.Execute(jobs, j)
	}

	for i := 0; i < len(b.Jobs); i++ {
		// TODO: Verify Job is actually done.
		// if not then post it again.
		job_result := <-jobs

		b.Results = append(b.Results, job_result)

		b.Task.Write(&job_result)
	}
}

func NewBatch(task Task) (batch *Batch) {
	// Handle more of the task cases.
	tasks := task.Tasks

	batch = &Batch{Task: task}

	if reflect.TypeOf(tasks).Kind() != reflect.Slice {
		fmt.Println("Wtf kind of shit is this?")
		return nil
	}

	// Iterate over the Tasks
	s := reflect.ValueOf(tasks)
	for i := 0; i < s.Len(); i++ {
		// Figure out the information for one task.
		task := s.Index(i)

		job := Job{}

		// TODO NEED TO ADD TASK: Task: task
		// fmt.Println(task.Type())
		// job.Task = task

		// Iterate over the fields of a struct
		for j := 0; j < task.NumField(); j++ {

			switch task.Type().Field(j).Type.Name() {
			case "InputField":
				job.InputFields = append(job.InputFields, JobField{
					// fmt.Println(task.Type().Field(j).Name)
					Id:          task.Type().Field(j).Tag.Get(WorkId),
					Description: task.Type().Field(j).Tag.Get(WorkDesc),
					Type:        task.Type().Field(j).Tag.Get(WorkType),
					Value:       task.Field(j).String(),
				})
			case "OutputField":
				job.OutputFields = append(job.OutputFields, JobField{
					// fmt.Println(task.Type().Field(j).Name)
					Id:          task.Type().Field(j).Tag.Get(WorkId),
					Description: task.Type().Field(j).Tag.Get(WorkDesc),
					Type:        task.Type().Field(j).Tag.Get(WorkType),
				})

			default:
				// fmt.Println("Unknown")
			}
		}

		batch.Jobs = append(batch.Jobs, job)
	}

	return
}
