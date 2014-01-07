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

/*
 * Task is a way to define the Job that needs to be run.
 */
// work_type: output_image, output_text, input_text
type Task struct {
	Title       string
	Description string
	Price       uint // In cents
	Write       func(j *MetaJob)
	Tasks       interface{} // TODO: Name should be renamed Work or something like that.
}

/*
 * Batch is a group of Jobs
 */

type Batch struct {
	Task     Task
	MetaJobs []MetaJob
}

func (b *Batch) Run(assigner Assigner) {
	meta_jobs := make(chan MetaJob)

	for _, j := range b.MetaJobs {
		go j.StartJobs(assigner, meta_jobs, j)
	}

	for i := 0; i < len(b.MetaJobs); i++ {
		// TODO: Verify Job is actually done.
		// if not then post it again.
		job_result := <-meta_jobs

		b.Task.Write(&job_result)
	}
}

func NewBatch(t Task) (batch *Batch) {
	// Handle more of the task cases.
	tasks := t.Tasks

	batch = &Batch{Task: t}

	if reflect.TypeOf(tasks).Kind() != reflect.Slice {
		fmt.Println("Wtf kind of shit is this?")
		return nil
	}

	// Iterate over the Tasks
	s := reflect.ValueOf(tasks)
	for i := 0; i < s.Len(); i++ {
		// Figure out the information for one task.
		task := s.Index(i)

		job := MetaJob{
			Title:       t.Title,
			Description: t.Description,
		}

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
					Type:        JobFieldType(task.Type().Field(j).Tag.Get(WorkType)),
					Value:       task.Field(j).String(),
				})
			case "OutputField":
				job.OutputFields = append(job.OutputFields, JobField{
					// fmt.Println(task.Type().Field(j).Name)
					Id:          task.Type().Field(j).Tag.Get(WorkId),
					Description: task.Type().Field(j).Tag.Get(WorkDesc),
					Type:        JobFieldType(task.Type().Field(j).Tag.Get(WorkType)),
				})

			default:
				// fmt.Println("Unknown")
			}
		}

		batch.MetaJobs = append(batch.MetaJobs, job)
	}

	return
}
