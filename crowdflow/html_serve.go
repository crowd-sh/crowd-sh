package crowdflow

import (
	"fmt"
	"github.com/gorilla/mux" // TODO: Deprecate this
	"log"
	"net/http"
	"net/url"
)

type IndexConfig struct {
	Name    string
	NumJobs int
	Url     string
	Port    string
}

var (
	HttpServerEnabled = false
	IndexConf         IndexConfig
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

func pingIndexServer() {
	v := url.Values{}
	v.Set("name", "Ava")
	v.Set("num_jobs", "Ava")
	v.Set("url", "Ava")

	http.PostForm("http://workmachine.us/v1/tasks", v)
}

func startHttpServer() {
	log.Println("Serving HTTP Server")

	r := mux.NewRouter()
	r.HandleFunc("/v1/assignment/split", getSplitAssignmentHandler).Methods("GET")
	r.HandleFunc("/v1/assignment/split", postSplitAssignmentHandler).Methods("POST")
	r.HandleFunc("/v1/assignment/form", getSplitAssignmentHandler).Methods("GET")
	r.HandleFunc("/v1/assignment/form", postSplitAssignmentHandler).Methods("POST")
	http.Handle("/", r)
	http.ListenAndServe(fmt.Sprintf(":%d", 3001), nil)
}

func EnableHtmlServer(indexCfg IndexConfig) {
	HttpServerEnabled = true

	go pingIndexServer()
	go startHttpServer()
}
