package main

import (
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	//	"github.com/crowdmob/goamz/aws"
	//	"github.com/crowdmob/goamz/exp/mturk"
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

//	MTurk mturk.MTurk
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

func mturkConnect() {
	//	auth := aws.Auth{AccessKey: "abc", SecretKey: "123"}
	//	MTurk = mturk.New(auth, true)
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
	//	go AssignmentExpirer()

	log.Println("WorkMachine Starting...")

	r := mux.NewRouter()
	r.HandleFunc("/v1/workflows", func(w http.ResponseWriter, r *http.Request) {
		workflow := NewWorkflow(r.Body)
		fmt.Fprintln(w, workflow.Id)
	}).Methods("POST")

	r.HandleFunc("/v1/workflows/{workflow}", func(w http.ResponseWriter, r *http.Request) {

	}).Methods("GET")

	r.HandleFunc("/v1/workflows/{workflow}/tasks", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		workflow := Workflow{}
		Db.First(&workflow, vars["workflow"])
		workflow.Parse()

		t := workflow.AddTask(r.Body)

		fmt.Fprintln(w, t.Id)
	}).Methods("POST")

	r.HandleFunc("/v1/workflows/{workflow}/tasks/{tasks}", func(w http.ResponseWriter, r *http.Request) {

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
