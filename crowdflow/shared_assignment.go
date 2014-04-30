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

type SharedAssignment struct {
	Assigned  bool         `json:"-"`
	Id        string       `json:"id"`
	StartedAt time.Time    `json:"started_at"`
	Mutex     sync.RWMutex `json:"-"`
	Finished  bool         `json:"-"`
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
