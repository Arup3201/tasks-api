package main

import (
	"github.com/Arup3201/gotasks/internal/handlers"
	"github.com/gin-gonic/gin"
)

const PORT = ":8000"

func main() {
	router := gin.Default()

	router.Use(handlers.HandleErrors())

	router.GET("/tasks", handlers.GetAllTasks)
	router.POST("/tasks", handlers.AddTask)
	router.GET("/tasks/:id", handlers.GetTaskWithId)
	router.PUT("/tasks/:id", handlers.EditTask)
	router.PUT("/tasks/:id/mark", handlers.MarkTask)
	router.DELETE("/tasks/:id", handlers.DeleteTask)
	router.GET("/search/tasks", handlers.SearchTask)

	router.Run("localhost:8080")
}
