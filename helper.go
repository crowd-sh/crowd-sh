package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func renderJson(w http.ResponseWriter, page interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	b, err := json.Marshal(page)
	if err != nil {
		log.Println("error:", err)
		fmt.Fprintf(w, "")
	}

	log.Println("Rendered Page")

	w.Write(b)
}
