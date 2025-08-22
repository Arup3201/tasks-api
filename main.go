package main

import (
	"fmt"
	"net/http"
	"os"
)

const PORT = ":8000"

func createTask(w http.ResponseWriter, r *http.Request) {
	if r.Method!=http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	value := r.FormValue("task")
	fmt.Println(value)

	fmt.Fprintf(w, "[TODO]")
}

func main() {
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/tasks", createTask)

	fmt.Printf("http server started at %s\n", PORT)
	err := http.ListenAndServe(PORT, nil)
	if err!=nil {
		fmt.Fprintf(os.Stderr, "http listen and server failed - %s", err)
	}
}
