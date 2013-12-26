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
		StartedAt: time.Now(),
	}

	assignments = append(assignments, assignment)
}

type Assignment struct {
	Assigned     bool
	JobsChan     chan Job
	Job          Job
	AssignmentId string
	StartedAt    time.Time
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

	input_html := ""
	for _, inp := range assign.Job.InputFields {
		input_html += `<div class=form-group>`
		input_html += fmt.Sprintf("<label>%s</label><br>", inp.Description)
		switch inp.Type {
		case ImageType:
			input_html += fmt.Sprintf("<img src=\"%s\" />", inp.Value)
		default:
			input_html += fmt.Sprintf("<p>%s</p>", inp.Value)
		}
		input_html += "</div>"
	}

	output_html := ""
	for _, out := range assign.Job.OutputFields {
		output_html += `<div class=form-group>`
		output_html += fmt.Sprintf("<label>%s</label>", out.Description)
		switch out.Type {
		case CheckBoxType:
			output_html += fmt.Sprintf("<input type=checkbox name=\"%s\" value=\"yes\"/>", out.Id)
		default:
			output_html += fmt.Sprintf("<br><input type=text name=\"%s\" value=\"%s\"/>", out.Id, out.Value)
		}
		output_html += "</div>"
	}

	output_html += fmt.Sprintf("<input type=hidden name=\"assignment_id\" value=\"%s\">", assign.AssignmentId)
	output_html += `<input type=submit value=Submit class="btn btn-default" >`

	template := `<html>
  <head>
    <title>Assignment</title>
    <link href="//netdna.bootstrapcdn.com/bootstrap/3.0.3/css/bootstrap.min.css" rel="stylesheet">
  </head>
  <body>
    <div class="container">
      <form method=post action="/assignment" role="form">
        <div>
          %s
        </div>

        <div>
          %s
        </div>
      </form>
    </div>
  </body>
</html>`

	fmt.Fprintf(w, template, input_html, output_html)
}

func postAssignmentHandler(w http.ResponseWriter, req *http.Request) {
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

	http.Redirect(w, req, "/", 302)
}

func Serve() {
	r := mux.NewRouter()
	r.HandleFunc("/", getAssignmentHandler).Methods("GET")
	r.HandleFunc("/assignment", postAssignmentHandler).Methods("POST")
	http.Handle("/", r)
	http.ListenAndServe(":5000", nil)
}
