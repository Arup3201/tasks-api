package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Arup3201/gotasks/internal/models"
	"github.com/Arup3201/gotasks/internal/storage"
	"github.com/Arup3201/gotasks/internal/utils"
	"github.com/gin-gonic/gin"
)

func GetAllTasks(c *gin.Context) {
	var tasks []models.Task = storage.GetAllTasks()
	c.IndentedJSON(http.StatusOK, tasks)
}

func GetTaskWithId(c *gin.Context) {
	id := c.Param("id")

	task, ok := storage.GetTaskWithId(id)
	if !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Task with ID '%s' not found", id)})
		return
	}

	c.IndentedJSON(http.StatusOK, task)
}

func AddTask(c *gin.Context) {
	var payload models.CreateTask

	if err := c.BindJSON(&payload); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error in unpacking the payload"})
		return
	}

	newTask := models.Task{
		Id:          utils.GenerateID("TASK_"),
		Title:       payload.Title,
		Description: payload.Description,
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	storage.AddTask(newTask)

	c.IndentedJSON(http.StatusCreated, newTask)
}

func EditTask(c *gin.Context) {
	var id = c.Param("id")
	var payload models.EditTask

	if err := c.BindJSON(&payload); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error in unpacking the payload"})
		return
	}

	task, ok := storage.EditTask(id, payload)
	if !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Task with ID '%s' not found", id)})
		return
	}

	c.IndentedJSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	var id = c.Param("id")
	task, ok := storage.DeleteTask(id)
	if !ok {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Task with ID '%s' not found", id)})
		return
	}

	c.IndentedJSON(http.StatusOK, task)
}
