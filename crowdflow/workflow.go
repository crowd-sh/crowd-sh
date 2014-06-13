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
 * TaskConfig is a way to define the Job that needs to be run.
 */
// work_type: output_image, output_text, input_text
type TaskConfig struct {
	Title       string
	Description string
	Tags        string
	Price       uint // In cents
	Write       func(j *MetaJob)
	Tasks       interface{} // TODO: Name should be renamed Work or something like that.
}

/*
 * Batch is a group of Jobs
 */

type Batch struct {
	TaskConfig TaskConfig
	MetaJobs   []MetaJob
}

type MetaJob struct {
	TaskConfig   *TaskConfig `json:"task_desc"`
	InputFields  []JobField  `json:"input_fields"`
	OutputFields []JobField  `json:"-"`
}

func (b *Batch) Run() {
	assignDone := make(chan bool)

	for _, j := range b.MetaJobs {
		go NewAssignment(assignDone, b, &j)
	}

	for i := 0; i < len(b.MetaJobs); i++ {
		<-assignDone
	}
}

func NewBatch(t TaskConfig) (batch *Batch) {
	// TODO: Handle more of the task cases.
	tasks := t.Tasks

	batch = &Batch{TaskConfig: t}

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
			TaskConfig: &t,
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
