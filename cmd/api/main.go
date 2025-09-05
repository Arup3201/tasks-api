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
	http.HandleFunc("/tasks", middlewares.Chain(methods.Map([]methods.MethodHandler{
		{
			Handler: handlers.GetTasks,
			Method:  "GET",
		},
		{
			Handler: handlers.AddTask,
			Method:  "POST",
		},
	}), logging.HttpLogger()))

	log.Printf("Starting the server at %s\n", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}
