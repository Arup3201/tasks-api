package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Arup3201/gotasks/internal/services"
	"github.com/gin-gonic/gin"
)

type CreateTask struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

type routeHandler struct {
	serviceHandler services.ServiceHandler
}

func getRouteHandler(handler services.ServiceHandler) *routeHandler {
	return &routeHandler{
		serviceHandler: handler,
	}
}

func (handler *routeHandler) GetTasks(c *gin.Context) {
	tasks, err := handler.serviceHandler.GetAllTasks()
	if err != nil {
		c.Error(err)
		return
	}
	c.IndentedJSON(http.StatusOK, tasks)
}

func (handler *routeHandler) AddTask(c *gin.Context) {
	var payload CreateTask

	if err := c.BindJSON(&payload); err != nil {
		c.Error(fmt.Errorf("c.BindJSON failed with error %v", err))
		return
	}

	if payload.Title == nil {
		c.Error(fmt.Errorf("Payload missing 'title'"))
		return
	}
	if payload.Description == nil {
		c.Error(fmt.Errorf("Payload missing 'description'"))
		return
	}

	newTask, err := handler.serviceHandler.CreateTask(*payload.Title, *payload.Description)
	if err != nil {
		c.Error(err)
		return
	}

	c.IndentedJSON(http.StatusCreated, newTask)
}

func (handler *routeHandler) GetTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(err)
		return
	}

	task, err := handler.serviceHandler.GetTask(id)
	if err != nil {
		c.Error(err)
		return
	}

	c.IndentedJSON(http.StatusOK, task)
}

func (handler *routeHandler) UpdateTask(c *gin.Context) {
	var id, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(err)
		return
	}
	var payload services.UpdateTaskData
	if err := c.BindJSON(&payload); err != nil {
		c.Error(fmt.Errorf("c.BindJSON failed with error %v", err))
		return
	}

	if payload.Title == nil && payload.Description == nil && payload.IsCompleted == nil {
		c.Error(fmt.Errorf("Update payload is empty"))
		return
	}

	editedTask, err := handler.serviceHandler.UpdateTask(id, payload)
	if err != nil {
		c.Error(err)
		return
	}

	c.IndentedJSON(http.StatusOK, editedTask)
}

func (handler *routeHandler) SearchTasks(c *gin.Context) {
	var query string = c.Query("q")
	tasks, err := handler.serviceHandler.SearchTasks(query)
	if err != nil {
		c.Error(err)
		return
	}

	c.IndentedJSON(http.StatusOK, tasks)
}
