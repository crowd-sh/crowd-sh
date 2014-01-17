package crowdflow

import (
	"fmt"
	"reflect"
)

const (
	WorkId   = "work_id"
	WorkDesc = "work_desc"
	WorkType = "work_type"
)

/*
 * TaskDesc is a way to define the Job that needs to be run.
 */
// work_type: output_image, output_text, input_text
type TaskDesc struct {
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
	TaskDesc TaskDesc
	MetaJobs []MetaJob
}

func (b *Batch) Run(assigner Assigner, splitJob bool) {
	meta_jobs := make(chan MetaJob)

	for _, j := range b.MetaJobs {
		go j.StartJob(assigner, meta_jobs, j)
	}

	for i := 0; i < len(b.MetaJobs); i++ {
		// TODO: Verify Job is actually done.
		// if not then post it again.
		job_result := <-meta_jobs

		b.TaskDesc.Write(&job_result)
	}
}

func (b *Batch) RunSplit(assigner SplitAssigner, splitJob bool) {
	meta_jobs := make(chan MetaJob)

	for _, j := range b.MetaJobs {
		go j.StartSplitJob(assigner, meta_jobs, j)
	}

	for i := 0; i < len(b.MetaJobs); i++ {
		// TODO: Verify Job is actually done.
		// if not then post it again.
		job_result := <-meta_jobs

		b.TaskDesc.Write(&job_result)
	}
}

func NewBatch(t TaskDesc) (batch *Batch) {
	// TODO: Handle more of the task cases.

	tasks := t.Tasks

	batch = &Batch{TaskDesc: t}

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
