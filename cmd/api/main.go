package main

import (
	"log"
	"net/http"

	"github.com/Arup3201/gotasks/internal/errors"
	"github.com/Arup3201/gotasks/internal/handlers"
	"github.com/gin-gonic/gin"
)

func HandleErrors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			clientError, ok := err.(errors.ClientError)
			if !ok {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "An unexpected error occured"})
				return
			}

			body, err := clientError.ResponseBody()
			if err != nil {
				log.Printf("ClientError.ResponseBody error: %v", err)
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "An unexpected error occured"})
				return
			}

			status, headers := clientError.ResponseHeader()
			for k, v := range headers {
				c.Writer.Header().Set(k, v)
			}
			c.Writer.WriteHeader(status)
			c.Writer.Write(body)
		}
	}
}

const PORT = ":8000"

func main() {
	router := gin.Default()

	router.Use(HandleErrors())

	router.GET("/tasks", handlers.GetAllTasks)
	router.POST("/tasks", handlers.AddTask)
	router.GET("/tasks/:id", handlers.GetTaskWithId)
	router.PUT("/tasks/:id", handlers.EditTask)
	router.PUT("/tasks/:id/mark", handlers.MarkTask)
	router.DELETE("/tasks/:id", handlers.DeleteTask)
	router.GET("/search/tasks", handlers.SearchTask)

	router.Run("localhost:8080")
}
