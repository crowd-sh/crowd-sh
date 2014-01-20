package crowdflow

import (
	"fmt"
	"github.com/gorilla/mux" // TODO: Deprecate this
	"net/http"
)

func getSplitAssignmentHandler(w http.ResponseWriter, req *http.Request) {
	assign := bulk_assignments.Get()
	if assign == nil {
		renderJson(w, false)
		return
	}

	renderJson(w, assign)
}

func postSplitAssignmentHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Posting", req.FormValue("id"))

	assign := split_assignments.Find(req.FormValue("id"))
	if assign != nil {
		assign.Finish(req.FormValue(assign.Job.OutputField.Id))
	}

	renderJson(w, true)
}

func HtmlServer() {
	r := mux.NewRouter()
	r.HandleFunc("/v1/assignment/split", getSplitAssignmentHandler).Methods("GET")
	r.HandleFunc("/v1/assignment/split", postSplitAssignmentHandler).Methods("POST")
	http.Handle("/", r)
	http.ListenAndServe(":5000", nil)
}
