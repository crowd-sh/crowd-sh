package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
)

type Backender interface {
	Publish(batch *Batch)
	Execute()
	Aggregate()
}

type TermServe struct {
	Batch *Batch
}

func (ss TermServe) Publish(batch *Batch) {
	ss.Batch = batch
}

func (ss TermServe) Execute() {
	for _, j := range ss.Batch.Jobs {
		for _, inp := range j.InputFields {
			fmt.Println(inp.Description)
			fmt.Println(inp.Value)
		}

		bio := bufio.NewReader(os.Stdin)
		for i, out := range j.OutputFields {
			fmt.Println(out.Description)
			fmt.Println(out.Value)
			line, _, _ := bio.ReadLine()
			j.OutputFields[i].Value = string(line)
		}
	}
}

func (ss TermServe) Aggregate() {
	fmt.Println(ss.Batch)

	// for _, j := range ss.Batch.Jobs {
	// 	// for _, out := range j.OutputFields {
	// 	// 	fmt.Printf("%s\t%s", out.Id, out.Description)
	// 	// }
	// 	fmt.Println(j)
	// }
}

/*
 * Task is a way to define the Job that needs to be run.
 */
// crowd_type: output_image, output_text, input_text
type Task struct {
	Title       string
	Description string
	Price       int // In cents
	Tasks       interface{}
}

type InputField string
type OutputField string

type BusinessCard struct {
	ImageUrl    InputField  `crowd_desc:"Use this information for the data below." crowd_id:"image_url" crowd_type:"image"`
	Name        OutputField `crowd_desc:"Find the name from the business card" crowd_id:"name"`
	Company     OutputField `crowd_desc:"Find the company from the business card" crowd_id:"company"`
	Email       OutputField `crowd_desc:"Find the email from the business card" crowd_id:"email"`
	PhoneNumber OutputField `crowd_desc:"Find the phone number from the business card" crowd_id:"phone_number"`
}

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
	Task         Task
	InputFields  []JobField
	OutputFields []JobField
}

func (j Job) Execute() {
	// fmt.Println("<div>")
	// for _, inp := range j.InputFields {
	// 	fmt.Println(inp.Html())
	// }
	// fmt.Println("</div>")

	// fmt.Println("<div>")
	// for _, out := range j.OutputFields {
	// 	fmt.Println(out.Html())
	// }
	// fmt.Println("</div>")
}

/*
 * Batch is a group of Jobs
 */

type Batch struct {
	Jobs []Job
}

func (b *Batch) Run(ss Backender) {
	ss.Publish(b)
	ss.Execute()

	ss.Aggregate()
}

func NewBatch(task ...Task) (batch *Batch) {
	// Handle more of the task cases.
	tasks := task[0].Tasks

	batch = &Batch{}

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

		// Iterate over the fields of a struct
		for j := 0; j < task.NumField(); j++ {

			switch task.Type().Field(j).Type.Name() {
			case "InputField":
				job.InputFields = append(job.InputFields, JobField{
					// fmt.Println(task.Type().Field(j).Name)
					Id:          task.Type().Field(j).Tag.Get("crowd_id"),
					Description: task.Type().Field(j).Tag.Get("crowd_desc"),
					Type:        task.Type().Field(j).Tag.Get("crowd_type"),
					Value:       task.Field(j).String(),
				})
			case "OutputField":
				job.OutputFields = append(job.OutputFields, JobField{
					// fmt.Println(task.Type().Field(j).Name)
					Id:          task.Type().Field(j).Tag.Get("crowd_id"),
					Description: task.Type().Field(j).Tag.Get("crowd_desc"),
					Type:        task.Type().Field(j).Tag.Get("crowd_type"),
				})

			default:
				// fmt.Println("Unknown")
			}
		}

		batch.Jobs = append(batch.Jobs, job)
	}

	return
}

func main() {
	business_cards := Task{
		Title:       "Business Card Fields",
		Description: "Enter the fields.",
		Price:       1,
		Tasks: []BusinessCard{
			BusinessCard{
				ImageUrl: "google.com",
			},
			// BusinessCard{
			// 	ImageUrl: "yahoo.com",
			// },
		},
	}

	serve := TermServe{}
	var backend Backender = serve
	NewBatch(business_cards).Run(backend)

}
