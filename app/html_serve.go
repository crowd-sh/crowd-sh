package workmachine

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"net/http"
	"sync"
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
	fmt.Println(j)

	<-assignment.JobsChan
}

type Assignment struct {
	Assigned     bool
	JobsChan     chan Job
	Job          Job
	AssignmentId string
	ExpiresIn    int

	Mutex    sync.RWMutex
	Finished bool
}

var (
	assignments []*Assignment
)

func getAssignment() *Assignment {
	for _, a := range assignments {
		a.Mutex.Lock()
		if !a.Assigned && !a.Finished {
			a.Assigned = true
			a.AssignmentId = string(securecookie.GenerateRandomKey(128))

			defer a.Mutex.Unlock()

			return a
		}
		a.Mutex.Unlock()
	}

	return nil
}

func findAssignment(assignment_id string) *Assignment {
	for _, a := range assignments {
		if fmt.Sprintf("%x", a.AssignmentId) == assignment_id {

			fmt.Printf("hello %s from %s", assignment_id, fmt.Sprintf("%x", a.AssignmentId))
			return a
		}
	}

	return nil
}

func getAssignmentHandler(w http.ResponseWriter, req *http.Request) {
	assign := getAssignment()
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

	output_html += fmt.Sprintf("<input type=text name=\"assignment_id\" value=\"%x\">", assign.AssignmentId)
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
	a := findAssignment(req.FormValue("assignment_id"))
	// a.Mutex.Lock()
	// defer a.Mutex.Unlock()

	fmt.Fprintf(w, "hi %v", a)
}

func Serve() {
	r := mux.NewRouter()
	r.HandleFunc("/assignment", getAssignmentHandler).Methods("GET")
	r.HandleFunc("/assignment", postAssignmentHandler).Methods("POST")
	http.Handle("/", r)
	http.ListenAndServe(":5000", nil)
}
