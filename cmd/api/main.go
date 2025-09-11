package main

import (
	"log"

	"github.com/Arup3201/gotasks/internal/handlers"
	"github.com/Arup3201/gotasks/internal/storage"
	"github.com/gin-gonic/gin"
)

const PORT = ":8000"

func main() {
	router := gin.Default()

	router.Use(handlers.HandleErrors())

	db, err := storage.NewPostgres()
	if err != nil {
		log.Fatalf("storage.NewPostgres error: %v", err)
	}
	taskHandler := handlers.NewTaskHandler(db)
	router.GET("/tasks", taskHandler.GetAllTasks)
	router.POST("/tasks", taskHandler.AddTask)
	router.GET("/tasks/:id", taskHandler.GetTaskWithId)
	router.PUT("/tasks/:id", taskHandler.EditTask)
	router.PUT("/tasks/:id/mark", taskHandler.MarkTask)
	router.DELETE("/tasks/:id", taskHandler.DeleteTask)
	router.GET("/search/tasks", taskHandler.SearchTask)
	router.Run("localhost:8080")
}
