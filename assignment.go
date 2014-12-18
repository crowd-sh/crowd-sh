package main

import (
	// "fmt"
	// "github.com/gorilla/securecookie" // TODO: Deprecate this
	"time"
)

const (
	ExpireAfterMinutes = 5
)

type Assignment struct {
	Id         string    `json:"id"`
	Assigned   bool      `json:"assigned"`
	StartedAt  time.Time `json:"started_at"`
	Finished   bool      `json:"finished"`
	AssignDone chan bool `json:"-"`
	Value      string    `json:"-"`
	//	Job        *Job      `json:"job"`
	//	InputField *JobField `json:"input_field"`
}

// func (a *Assignment) generateId() string {
// 	return fmt.Sprintf("%x", string(securecookie.GenerateRandomKey(128)))
// }

// func (a *Assignment) Assign() {
// 	a.Assigned = true
// 	a.Id = a.generateId()
// 	a.StartedAt = time.Now()
// }

// func (a *Assignment) TryToAssign() bool {
// 	a.Mutex.Lock()
// 	defer a.Mutex.Unlock()

// 	a.UnassignIfExpired()

// 	if !a.Assigned && !a.Finished {
// 		a.Assign()

// 		return true
// 	}

// 	return false
// }

// func (a *Assignment) UnassignIfExpired() {
// 	duration := time.Since(a.StartedAt) / time.Minute
// 	if duration > ExpireAfterMinutes { // Greater than 5 minutes
// 		a.Assigned = false
// 		a.Id = ""
// 	}
// }

// func (a *Assignment) Finish(value string) {
// 	a.Mutex.Lock()
// 	a.Value = value
// 	a.Finished = true
// 	a.Mutex.Unlock()

// 	a.AssignDone <- true
// }

// func (a *Assignment) Run() {
// 	select {
// 	// case res := <-c1:
// 	// 	fmt.Println(res)
// 	case <-time.After(time.Minute * ExpireAfterMinutes):
// 		a.UnassignIfExpired()
// 	}
// }

// func NewAssignment(assignDone chan bool, j *Job, jf *JobField) {
// 	assignment := Assignment{
// 		Job:        j,
// 		InputField: jf,
// 		AssignDone: assignDone,
// 		Finished:   false,
// 		Assigned:   false,
// 	}

// 	AvailableAssignments = append(AvailableAssignments, assignment)
// }

// type Assignments []Assignment

// func (as Assignments) GetUnfinished() *Assignment {
// 	for i := range as {
// 		a := as[i]

// 		if !a.Finished && !a.Assigned {
// 			fmt.Println(a.Id)
// 			return &a
// 		}
// 	}

// 	return nil
// }

// func (as Assignments) Find(id string) *Assignment {
// 	for _, a := range as {
// 		if a.Id == id {
// 			return &a
// 		}
// 	}

// 	return nil
// }

// func ExpireAssignments() {

// 	// TODO: Iterate over the assignments and expire the ones which are over 5 minutes.
// }
