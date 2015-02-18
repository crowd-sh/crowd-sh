package main

import (
	"github.com/coreos/go-etcd/etcd"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

const (
	DatabaseUrlKey = "/workmachine/database_url"
)

var (
	Db gorm.DB
)

func dbConnect(databaseUrl string) {
	log.Println("Connecting to database:", databaseUrl)
	var err error
	Db, err = gorm.Open("postgres", databaseUrl)
	if err != nil {
		log.Println(err)
	}
	Db.LogMode(true)
}

func init() {
	etcdHosts := os.Getenv("ETCD_HOSTS")
	if etcdHosts == "" {
		etcdHosts = "http://127.0.0.1:4001"
	}

	etcdClient := etcd.NewClient([]string{etcdHosts})

	resp, err := etcdClient.Get(DatabaseUrlKey, false, false)
	if err != nil {
		panic(err)
	}

	databaseUrl := resp.Node.Value
	dbConnect(databaseUrl)
}

func main() {
	go AssignmentExpirer()

	log.Println("WorkMachine Starting...")

	r := mux.NewRouter()
	r.HandleFunc("/v1/workflow", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		for key, _ := range r.Form {
			log.Println(key)
			err := json.Unmarshal([]byte(key), &t)
			if err != nil {
				log.Println(err.Error())
			}
		}

		workflow := NewWorkflow("Json")
	}).Methods("POST")

	r.HandleFunc("/v1/workflow/{workflow}", func(w http.ResponseWriter, r *http.Request) {

	}).Methods("GET")

	r.HandleFunc("/v1/workflow/{workflow}/tasks", func(w http.ResponseWriter, r *http.Request) {
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

	}).Methods("POST")

	r.HandleFunc("/v1/workflow/{workflow}/tasks/{tasks}", func(w http.ResponseWriter, r *http.Request) {

	}).Methods("GET")

	//r.HandleFunc("/v1/workflow/{workflow}/tests", TaskHandler).Methods("PUT")

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
	http.ListenAndServe(":3002", nil)
}
