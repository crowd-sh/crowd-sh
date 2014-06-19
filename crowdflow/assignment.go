package crowdflow

import (
	"fmt"
	"github.com/gorilla/securecookie" // TODO: Deprecate this
	"sync"
	"time"
)

const (
	ExpireAfterMinutes = 5
)

type Assignment struct {
	Assigned   bool         `json:"-"`
	Id         string       `json:"id"`
	StartedAt  time.Time    `json:"started_at"`
	Mutex      sync.RWMutex `json:"-"`
	Finished   bool         `json:"-"`
	AssignDone chan bool    `json:"-"`
	// Job        *MetaJob  `json:"job"`
}

func (a *SharedAssignment) generateId() string {
	return fmt.Sprintf("%x", string(securecookie.GenerateRandomKey(128)))
}

func (a *SharedAssignment) Assign() {
	a.Assigned = true
	a.Id = a.generateId()
	a.StartedAt = time.Now()
}

func (a *SharedAssignment) TryToAssign() bool {
	a.Mutex.Lock()
	defer a.Mutex.Unlock()

	a.UnassignIfExpired()

	if !a.Assigned && !a.Finished {
		a.Assign()

		return true
	}

	return false
}

func (a *SharedAssignment) UnassignIfExpired() {
	duration := time.Since(a.StartedAt) / time.Minute
	if duration > ExpireAfterMinutes { // Greater than 5 minutes
		a.Assigned = false
		a.Id = ""
	}
}

type Assignments []Assignment

var (
	AvailableAssignments Assignments
)

func NewAssignment(assignDone chan bool, b *Program, j *MetaJob) {
	assignment := Assignment{
		AssignDone: assignDone,
		Job:        j,
	}

	AvailableAssignments = append(AvailableAssignments, assignment)
}

func (as Assignments) GetUnfinished() *Assignment {
	for _, a := range as {
		if !a.SharedAssignment.Finished && !a.SharedAssignment.Assigned {
			return &a
		}
	}

	return nil
}

func (a Assignment) Finish(val string) {

	// b.TaskConfig.Write(a.Job)
}
