package workmachine

import (
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
		Assigned:  false,
		JobsChan:  jobs,
		Job:       j,
		ExpiresIn: 1000,
	}
	assignments = append(assignments, assignment)
}

type Assignment struct {
	Assigned     bool
	JobsChan     chan Job
	Job          Job
	AssignmentId string
	ExpiresIn    time.Duration
	Mutex        sync.RWMutex
	Finished     bool
}

func (a *Assignment) generateAssignmentId() string {
	return fmt.Sprintf("%x", string(securecookie.GenerateRandomKey(128)))
}

func (a *Assignment) Assign() {
	a.Assigned = true
	a.AssignmentId = a.generateAssignmentId()
}

func (a *Assignment) UnassignExpired() {
	// TODO
	// a.Assigned = false
	// a.AssignmentId = ""
}

type Assignments []*Assignment

var (
	assignments Assignments
)

func (as Assignments) Get() *Assignment {
	for _, a := range as {
		a.Mutex.Lock()

		a.UnassignExpired()

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
	assign := assignments.Get()
	if assign == nil {
		w.Write([]byte("Nothing to do."))
		return
	}

	input_html := ""
	for _, inp := range assign.Job.InputFields {
		input_html += "<div>"
		input_html += fmt.Sprintf("<label>%s</label>", inp.Description)
		input_html += fmt.Sprintf("<input type=text value=\"%s\"/>", inp.Value)
		input_html += "</div>"
	}

	output_html := ""
	for _, out := range assign.Job.OutputFields {
		output_html += "<div>"
		output_html += fmt.Sprintf("<label>%s</label><br>", out.Description)
		output_html += fmt.Sprintf("<input type=text name=\"%s\" value=\"%s\"/>", out.Id, out.Value)
		output_html += "</div>"
	}

	output_html += fmt.Sprintf("<input type=text name=\"assignment_id\" value=\"%s\">", assign.AssignmentId)
	output_html += "<input type=submit>"

	template := `<html>
<head><title>Assignment</title></head>
<body>
<form method=post>
  <div>
    %s
  </div>

  <div>
    %s
  </div>
</form>
</body>
</html>`

	fmt.Fprintf(w, template, input_html, output_html)
}

func postAssignmentHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Posting", req.FormValue("assignment_id"))
	assign := assignments.Find(req.FormValue("assignment_id"))

	assign.Mutex.Lock()
	for i, out := range assign.Job.OutputFields {
		assign.Job.OutputFields[i].Value = req.FormValue(out.Id)
	}

	assign.Finished = true
	assign.Mutex.Unlock()

	assign.JobsChan <- assign.Job

	fmt.Fprintf(w, "Submitted")
}

func Serve() {
	r := mux.NewRouter()
	r.HandleFunc("/assignment", getAssignmentHandler).Methods("GET")
	r.HandleFunc("/assignment", postAssignmentHandler).Methods("POST")
	http.Handle("/", r)
	http.ListenAndServe(":5000", nil)
}
