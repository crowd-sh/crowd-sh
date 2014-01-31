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

var (
	AvailableAssignments Assignments
)

type SharedAssignment struct {
	Assigned  bool         `json:"-"`
	Id        string       `json:"id"`
	StartedAt time.Time    `json:"started_at"`
	Mutex     sync.RWMutex `json:"-"`
	Finished  bool         `json:"-"`
}

type Assignment struct {
	SharedAssignment

	AssignDone chan bool `json:"-"`
	Job        *MetaJob  `json:"job"`
}

type Assignments []Assignment

func NewAssignment(assignDone chan bool, b *Batch, j *MetaJob) {
	assignment := Assignment{
		AssignDone: assignDone,
		Job:        j,
	}
	go assignment.Execute(b, assignDone)

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

func (a *Assignment) Execute(b *Batch, assignDone chan bool) {
	go NewMTurkAssigner(assignDone, a)

	<-assignDone

	b.TaskDesc.Write(a.Job)
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

func (a Assignment) Finish(val string) {

}
