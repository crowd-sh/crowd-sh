package main

import (
	"time"
)

const (
	Working  = "Working"
	Finished = "Finished"
)

type Task struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Status    string
	CreatedAt time.Time `json:"-"`

	TaskFields []TaskField
}

type TaskField struct {
	TaskId int64

	Type  string
	Key   string
	Value string
}

// func (t *Task) CreateAssignments() {

// }

// func (t *Task) Status() {
// 	// All the assignments are done
// }

// func NewTask(w *Workflow) {
// 	assignDone := make(chan bool)

// 	// Create the assignments
// 	log.Println("Creating assignments of jobs", len(p.Jobs))

// 	for i := range p.Jobs {
// 		for j := range p.Jobs[i].OutputFields {
// 			go NewAssignment(assignDone, &p.Jobs[i], &p.Jobs[i].OutputFields[j])
// 		}
// 	}

// 	// Wait for responses
// 	log.Println("Created assignments")
// 	for _ = range assignDone {
// 		log.Println("Waiting for assignments:", len(AvailableAssignments))
// 	}
// }
