package main

import (
	"log"
	"net/http"

	"github.com/Arup3201/gotasks/internal/handlers"
	"github.com/Arup3201/gotasks/internal/middlewares"
	"github.com/Arup3201/gotasks/internal/middlewares/logging"
	"github.com/Arup3201/gotasks/internal/middlewares/methods"
)

const PORT = ":8000"

func main() {
	http.HandleFunc("/tasks", middlewares.Chain(handlers.GetTasks, logging.HttpLogger(), methods.Methods([]string{"GET"})))
	http.HandleFunc("/tasks", middlewares.Chain(handlers.AddTask, logging.HttpLogger(), methods.Methods([]string{"POST"})))
	http.HandleFunc("/tasks", middlewares.Chain(handlers.EditTask, logging.HttpLogger(), methods.Methods([]string{"PUT"})))

	log.Printf("Starting the server at %s\n", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}
