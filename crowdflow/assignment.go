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
	Assigned  bool         `json:"-"`
	Id        string       `json:"id"`
	StartedAt time.Time    `json:"started_at"`
	Mutex     sync.RWMutex `json:"-"`
	Finished  bool         `json:"-"`
}

func (a *Assignment) generateId() string {
	return fmt.Sprintf("%x", string(securecookie.GenerateRandomKey(128)))
}

func (a *Assignment) Assign() {
	a.Assigned = true
	a.Id = a.generateId()
	a.StartedAt = time.Now()
}

func (a *Assignment) UnassignIfExpired() {
	duration := time.Since(a.StartedAt) / time.Minute
	if duration > ExpireAfterMinutes { // Greater than 5 minutes
		a.Assigned = false
		a.Id = ""
	}
}
