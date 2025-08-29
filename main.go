package main

import (
	"log"
	"net/http"

	tasks "github.com/Arup3201/gotasks/internal/handlers"
)

const PORT = ":8000"

func main() {
	http.HandleFunc("/tasks", func (w http.ResponseWriter, r *http.Request) {
		if r.Method==http.MethodGet {
			tasks.GetTasks(w, r)
		}
	})


	log.Printf("Starting the server at %s\n", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}
