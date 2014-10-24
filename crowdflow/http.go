package crowdflow

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	//	"time"
)

// TODO:
// - [ ] Check to make sure api key exists
// - [ ] Create a new entry in the Database
func newWorkflowHandler(w http.ResponseWriter, r *http.Request) {
	workflow := NewWorkflow("Json")
	Db.Create(&workflow)

	// args := flag.Args()
	// if len(args) < 2 {
	// 	log.Println("Crowdflow program file and data csv are missing.")
	// 	os.Exit(1)
	// }

	// csvFile, err := os.Open(args[1])
	// if err != nil {
	// 	log.Fatal(err)
	// }

}

func showWorkflowHandler(w http.ResponseWriter, r *http.Request) {

}

// func newTaskHandler(w http.ResponseWriter, req *http.Request) {
// 	for _, name := range []string{"name", "num_jobs", "url"} {
// 		if req.FormValue(name) == "" {
// 			renderJson(w, fmt.Sprintf("error: Need value %s", name))
// 			return
// 		}
// 	}

// 	task := Task{
// 		Name:      req.FormValue("name"),
// 		NumJobs:   req.FormValue("num_jobs"),
// 		Url:       req.FormValue("url"),
// 		CreatedAt: time.Now(),
// 	}
// 	task.GenerateId()

// 	tasks = append(tasks, task)

// 	log.Println("New Task", task.Id, req.FormValue("name"), req.FormValue("num_jobs"), req.FormValue("url"))

// 	fmt.Fprintln(w, task.Id)
// }

func StartServer() {
	log.Println("CrowdFlow Starting...")

	r := mux.NewRouter()
	// r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	http.Redirect(w, r, "/static", 302)
	// })

	r.HandleFunc("/v1/workflows", newWorkflowHandler).Methods("PUT")
	r.HandleFunc("/v1/workflow/{workflow}", showWorkflowHandler).Methods("GET")

	//r.HandleFunc("/v1/workflow/{workflow}/tasks", TaskHandler).Methods("PUT")
	//r.HandleFunc("/v1/workflow/{workflow}/tests", TaskHandler).Methods("PUT")
	//r.HandleFunc("/v1/tasks", newTaskHandler).Methods("POST")
	//r.HandleFunc("/v1/assignments", newTaskHandler).Methods("POST")

	// r.HandleFunc("/v1/assignments", func(w http.ResponseWriter, req *http.Request) {
	// 	assign := AvailableAssignments.GetUnfinished()
	// 	if assign == nil {
	// 		renderJson(w, false)
	// 		return
	// 	}

	// 	if !assign.TryToAssign() {
	// 		renderJson(w, false)
	// 		return
	// 	}

	// 	renderJson(w, assign)
	// }).Methods("GET")

	// r.HandleFunc("/v1/assignments", func(w http.ResponseWriter, req *http.Request) {
	// 	log.Println("Posting", req.FormValue("id"))

	// 	assign := AvailableAssignments.Find(req.FormValue("id"))
	// 	if assign != nil {
	// 		assign.Finish(req.FormValue(assign.InputField.Value))
	// 	}

	// 	renderJson(w, true)
	// }).Methods("POST")

	r.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path[1:])
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	http.Handle("/", r)
	http.ListenAndServe(fmt.Sprintf(":%s", Port), nil)
}
