package workmachine

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"          // TODO: Deprecate this
	"github.com/gorilla/securecookie" // TODO: Deprecate this
	"net/http"
	"sync"
	"time"
)

type HtmlServe struct{}

func (ss HtmlServe) Execute(jobs chan Job, j Job) {
	assignment := &Assignment{
		Assigned: false,
		JobsChan: jobs,
		Job:      j,
		Finished: false,
	}

	assignments = append(assignments, assignment)
}

type Assignment struct {
	Assigned     bool         `json:"-"`
	JobsChan     chan Job     `json:"-"`
	Job          Job          `json:"job"`
	AssignmentId string       `json:"assignment_id"`
	StartedAt    time.Time    `json:"started_at"`
	Mutex        sync.RWMutex `json:"-"`
	Finished     bool         `json:"-"`
}

func (a *Assignment) generateAssignmentId() string {
	return fmt.Sprintf("%x", string(securecookie.GenerateRandomKey(128)))
}

func (a *Assignment) Assign() {
	a.Assigned = true
	a.AssignmentId = a.generateAssignmentId()
	a.StartedAt = time.Now()
}

const (
	ExpireAfterMinutes = 5
)

func (a *Assignment) UnassignIfExpired() {
	duration := time.Since(a.StartedAt) / time.Minute
	if duration > ExpireAfterMinutes { // Greater than 5 minutes
		a.Assigned = false
		a.AssignmentId = ""
	}
}

type Assignments []*Assignment

var (
	assignments Assignments
)

func (as Assignments) Get() *Assignment {
	for _, a := range as {
		a.Mutex.Lock()

		a.UnassignIfExpired()

		if !a.Assigned && !a.Finished {
			defer a.Mutex.Unlock()
			a.Assign()
			return a
		}
		a.Mutex.Unlock()
	}

	return nil
}

func (as Assignments) Find(assignment_id string) *Assignment {
	for _, a := range as {
		if a.AssignmentId == assignment_id {
			return a
		}
	}

	return nil
}

func getAssignmentHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Num of Assignments", len(assignments))
	assign := assignments.Get()
	if assign == nil {
		w.Write([]byte("Nothing to do."))
		return
	}

	renderJson(w, assign)
}

func renderJson(w http.ResponseWriter, page interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	fmt.Printf("%v", page)

	b, err := json.Marshal(page)
	if err != nil {
		//log.Println("error:", err)
		fmt.Fprintf(w, "")
	}

	w.Write(b)
}

func postAssignmentHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	fmt.Println("Posting Valies", req.Form)
	fmt.Println("Posting", req.FormValue("assignment_id"))

	assign := assignments.Find(req.FormValue("assignment_id"))
	if assign != nil {
		assign.Mutex.Lock()
		for i, out := range assign.Job.OutputFields {
			assign.Job.OutputFields[i].Value = req.FormValue(out.Id)
		}

		assign.Finished = true
		assign.Mutex.Unlock()

		assign.JobsChan <- assign.Job
	}

	renderJson(w, true)
}

func Serve() {
	r := mux.NewRouter()
	r.HandleFunc("/v1/assignment", getAssignmentHandler).Methods("GET")
	r.HandleFunc("/v1/assignment", postAssignmentHandler).Methods("POST")
	http.Handle("/", r)
	http.ListenAndServe(":5000", nil)
}
