package crowdflow

import (
	"fmt"
	"github.com/gorilla/mux" // TODO: Deprecate this
	"log"
	"net/http"
)

func getFormAssignmentHandler(w http.ResponseWriter, req *http.Request) {
	assign := GetFormAssignment()
	if assign == nil {
		renderJson(w, false)
		return
	}

	renderJson(w, assign)
}

func postFormAssignmentHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Posting", req.FormValue("id"))

	assign := FindFormAssignment(req.FormValue("id"))
	if assign != nil {
		assign.Finish(req.FormValue(assign.Job.OutputFields[0].Id))
	}

	renderJson(w, true)
}

func getSplitAssignmentHandler(w http.ResponseWriter, req *http.Request) {
	assign := GetSplitAssignment()
	if assign == nil {
		renderJson(w, false)
		return
	}

	renderJson(w, assign)
}

func postSplitAssignmentHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Posting", req.FormValue("id"))

	assign := FindSplitAssignment(req.FormValue("id"))
	if assign != nil {
		assign.Finish(req.FormValue(assign.OutputField.Id))
	}

	renderJson(w, true)
}

func EnableHtmlServer(port uint) {
	go func() {
		log.Println("Serving HTTP Server")
		r := mux.NewRouter()
		r.HandleFunc("/v1/assignment/split", getSplitAssignmentHandler).Methods("GET")
		r.HandleFunc("/v1/assignment/split", postSplitAssignmentHandler).Methods("POST")
		r.HandleFunc("/v1/assignment/form", getSplitAssignmentHandler).Methods("GET")
		r.HandleFunc("/v1/assignment/form", postSplitAssignmentHandler).Methods("POST")
		http.Handle("/", r)
		http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}()
}
