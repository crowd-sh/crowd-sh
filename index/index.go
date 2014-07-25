package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	uuid "github.com/nu7hatch/gouuid"
	redis "github.com/xuyu/goredis"
	"log"
	"net/http"
)

const (
	ExpirationSeconds = 300 // 5 minutes
)

var (
	client *redis.Redis

	port       string
	redisDeets string
)

func init() {
	var err error

	client, err = redis.Dial(&redis.DialConfig{
		Address: redisDeets,
	})

	if err != nil {
		panic(err)
	}

	flag.StringVar(&port, "port", "3000", "Port")
	flag.StringVar(&redisDeets, "redis", "127.0.0.1:6379", "Redis Host:Port")
	flag.Parse()
}

func renderJson(w http.ResponseWriter, page interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	b, err := json.Marshal(page)
	if err != nil {
		log.Println("error:", err)
		fmt.Fprintf(w, "")
	}

	w.Write(b)
}

func indexTaskHandler(w http.ResponseWriter, req *http.Request) {
	taskIds, err := client.Keys("WM_TASK:*")
	if err != nil {
		panic(err)
	}

	fmt.Println(taskIds)

	type Task struct {
		Id      string `json:"id"`
		Title   string `json:"title"`
		NumJobs string `json:"num_jobs"`
		Url     string `json:"url"`
	}

	var tasks []Task

	for _, taskId := range taskIds {
		task, _ := client.HGetAll(taskId)

		tasks = append(tasks, Task{
			Id:      taskId,
			Title:   task["title"],
			NumJobs: task["num_jobs"],
			Url:     task["url"],
		})
	}

	renderJson(w, tasks)
}

func newTaskHandler(w http.ResponseWriter, req *http.Request) {
	u, _ := uuid.NewV4()
	taskId := fmt.Sprintf("WM_TASK:%s", u.String())

	for _, name := range []string{"name", "num_jobs", "url"} {
		if req.FormValue(name) == "" {
			renderJson(w, fmt.Sprintf("error: Need value %s", name))
			return
		}
	}

	client.HSet(taskId, "name", req.FormValue("name"))
	client.HSet(taskId, "num_jobs", req.FormValue("num_jobs"))
	client.HSet(taskId, "url", req.FormValue("url"))
	client.Expire(taskId, ExpirationSeconds)

	log.Println("New Task", taskId, req.FormValue("name"), req.FormValue("num_jobs"), req.FormValue("url"))

	fmt.Fprintln(w, taskId)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static", 302)
	})
	r.HandleFunc("/v1/tasks", indexTaskHandler).Methods("GET")
	r.HandleFunc("/v1/tasks", newTaskHandler).Methods("POST")
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path[1:])
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	http.Handle("/", r)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
