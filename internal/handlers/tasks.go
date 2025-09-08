package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Arup3201/gotasks/internal/handlers/clienterror"
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
		c.Error(clienterror.NewNotFoundError(nil))
		return
	}

	c.IndentedJSON(http.StatusOK, task)
}

func AddTask(c *gin.Context) {
	var payload models.CreateTask

	if err := c.BindJSON(&payload); err != nil {
		c.Error(fmt.Errorf("c.BindJSON failed with error %v", err))
		return
	}

	if payload.Title == nil {
		c.Error(clienterror.NewMissingBodyProperyError(nil, []clienterror.ErrorDetail{
			{
				Detail:  "The body property 'title' is required",
				Pointer: "#/title",
			},
		}))
		return
	}

	if payload.Description == nil {
		c.Error(clienterror.NewMissingBodyProperyError(nil, []clienterror.ErrorDetail{
			{
				Detail:  "The body property 'description' is required",
				Pointer: "#/description",
			},
		}))
		return
	}

	if *payload.Title == "" {
		c.Error(clienterror.NewInvalidBodyValueError(nil, []clienterror.ErrorDetail{
			{
				Detail:  "Body property 'title' can't be empty",
				Pointer: "#/title",
			},
		}))
		return
	}

	if *payload.Description == "" {
		c.Error(clienterror.NewInvalidBodyValueError(nil, []clienterror.ErrorDetail{
			{
				Detail:  "Body property 'description' can't be empty",
				Pointer: "#/description",
			},
		}))
		return
	}

	newTask := models.Task{
		Id:          utils.GenerateID("TASK_"),
		Title:       *payload.Title,
		Description: *payload.Description,
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
		c.Error(fmt.Errorf("c.BindJSON failed with error %v", err))
		return
	}

	if payload.Title == nil && payload.Description == nil {
		c.Error(clienterror.NewMissingBodyProperyError(nil, []clienterror.ErrorDetail{
			{
				Detail:  "Atleast one body property 'title' or 'description' is required",
				Pointer: "#/title, #/description",
			},
		}))
		return
	}

	if payload.Title != nil && *payload.Title == "" {
		c.Error(clienterror.NewInvalidBodyValueError(nil, []clienterror.ErrorDetail{
			{
				Detail:  "Body property 'title' can't be empty",
				Pointer: "#/title",
			},
		}))
		return
	}

	if payload.Description != nil && *payload.Description == "" {
		c.Error(clienterror.NewInvalidBodyValueError(nil, []clienterror.ErrorDetail{
			{
				Detail:  "Body property 'description' can't be empty",
				Pointer: "#/description",
			},
		}))
		return
	}

	task, ok := storage.UpdateTask(id, models.UpdateTask{
		Title:       payload.Title,
		Description: payload.Description,
		Completed:   nil,
	})
	if !ok {
		c.Error(clienterror.NewNotFoundError(nil))
		return
	}

	c.IndentedJSON(http.StatusOK, task)
}

func MarkTask(c *gin.Context) {
	var id = c.Param("id")
	var payload models.MarkTask

	if err := c.BindJSON(&payload); err != nil {
		c.Error(fmt.Errorf("c.BindJSON failed with error %v", err))
		return
	}

	if payload.Completed == nil {
		c.Error(clienterror.NewMissingBodyProperyError(nil, []clienterror.ErrorDetail{
			{
				Detail:  "Body property 'completed' is required",
				Pointer: "#/completed",
			},
		}))
		return
	}

	task, ok := storage.UpdateTask(id, models.UpdateTask{
		Title:       nil,
		Description: nil,
		Completed:   payload.Completed,
	})
	if !ok {
		c.Error(clienterror.NewNotFoundError(nil))
		return
	}

	c.IndentedJSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	var id = c.Param("id")
	task, ok := storage.DeleteTask(id)
	if !ok {
		c.Error(clienterror.NewNotFoundError(nil))
		return
	}

	c.IndentedJSON(http.StatusOK, task)
}

func SearchTask(c *gin.Context) {
	var query string = c.Query("q")

	if query == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "'q' can't be empty while searching"})
		return
	}

	tasks := storage.SearchTasks(query)

	c.IndentedJSON(http.StatusOK, tasks)
}
