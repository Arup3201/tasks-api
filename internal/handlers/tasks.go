package tasks

import (
	"log"
	"net/http"
)

func GetTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Fatalf("%s method not allowed on GetTasks", r.Method)
		return
	}


}

