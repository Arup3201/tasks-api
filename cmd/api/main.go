package main

import (
	"log"
	"net/http"

	"github.com/Arup3201/gotasks/internal/handlers"
)

const PORT = ":8000"

func main() {
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.GetTasks(w, r)
			return
		}

		if r.Method == http.MethodPost {
			handlers.AddTask(w, r)
			return
		}

		if r.Method == http.MethodPut {
			handlers.EditTask(w, r)
			return
		}

		http.Error(w, "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed)
	})

	log.Printf("Starting the server at %s\n", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}
