package handlers

import (
	"fmt"
	"net/http"

	"github.com/Arup3201/gotasks/internal/models"
	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	store models.TaskStore
}

func NewTaskHandler(store models.TaskStore) TaskHandler {
	return TaskHandler{
		store: store,
	}
}

func (handler TaskHandler) GetAllTasks(c *gin.Context) {
	model := models.NewTaskModel(handler.store)
	tasks, _ := model.AllTasks()
	c.IndentedJSON(http.StatusOK, tasks)
}

func (handler TaskHandler) GetTaskWithId(c *gin.Context) {
	id := c.Param("id")

	model := models.NewTaskModel(handler.store)
	task, err := model.GetTaskByID(id)
	if err != nil {
		c.Error(err)
		return
	}

	c.IndentedJSON(http.StatusOK, task)
}

func (handler TaskHandler) AddTask(c *gin.Context) {
	var payload models.CreateTask

	if err := c.BindJSON(&payload); err != nil {
		c.Error(fmt.Errorf("c.BindJSON failed with error %v", err))
		return
	}

	model := models.NewTaskModel(handler.store)
	newTask, err := model.AddTask(payload)
	if err != nil {
		c.Error(err)
		return
	}

	c.IndentedJSON(http.StatusCreated, newTask)
}

func (handler TaskHandler) EditTask(c *gin.Context) {
	var id = c.Param("id")
	var payload models.EditTask

	if err := c.BindJSON(&payload); err != nil {
		c.Error(fmt.Errorf("c.BindJSON failed with error %v", err))
		return
	}

	model := models.NewTaskModel(handler.store)
	editedTask, err := model.EditTask(id, payload)
	if err != nil {
		c.Error(err)
		return
	}

	c.IndentedJSON(http.StatusOK, editedTask)
}

func (handler TaskHandler) MarkTask(c *gin.Context) {
	var id = c.Param("id")
	var payload models.MarkTask

	if err := c.BindJSON(&payload); err != nil {
		c.Error(fmt.Errorf("c.BindJSON failed with error %v", err))
		return
	}

	model := models.NewTaskModel(handler.store)
	editedTask, err := model.MarkTask(id, payload)
	if err != nil {
		c.Error(err)
		return
	}

	c.IndentedJSON(http.StatusOK, editedTask)
}

func (handler TaskHandler) DeleteTask(c *gin.Context) {
	var id = c.Param("id")
	model := models.NewTaskModel(handler.store)
	taskId, err := model.DeleteTask(id)
	if err != nil {
		c.Error(err)
		return
	}

	c.IndentedJSON(http.StatusOK, taskId)
}

func (handler TaskHandler) SearchTask(c *gin.Context) {
	var query string = c.Query("q")
	model := models.NewTaskModel(handler.store)
	tasks, err := model.SearchTask(query)
	if err != nil {
		c.Error(err)
		return
	}

	c.IndentedJSON(http.StatusOK, tasks)
}
