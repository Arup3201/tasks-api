package controller

import (
	"net/http"
	"testing"

	controllers "github.com/Arup3201/gotasks/internal/controllers/http"
	"github.com/Arup3201/gotasks/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestUpdateTitleSuccess(t *testing.T) {
	title, description := "Test task 1", "Test description 1"
	newTask := controllers.CreateTask{
		Title:       &title,
		Description: &description,
	}
	makeRequest("POST", "/tasks", newTask)
	title = "Test task 1(edited)"
	editTask := services.UpdateTaskData{
		Title: &title,
	}

	response := makeRequest("PATCH", "/tasks/1", editTask)

	assert.Equal(t, http.StatusOK, response.Code)
}

func TestUpdateDescriptionSuccess(t *testing.T) {
	title, description := "Test task 1", "Test description 1"
	newTask := controllers.CreateTask{
		Title:       &title,
		Description: &description,
	}
	makeRequest("POST", "/tasks", newTask)
	description = "Test description 1(Edited)"
	editTask := services.UpdateTaskData{
		Description: &description,
	}

	response := makeRequest("PATCH", "/tasks/1", editTask)

	assert.Equal(t, http.StatusOK, response.Code)
}

func TestUpdateIsCompletedSuccess(t *testing.T) {
	title, description := "Test task 1", "Test description 1"
	newTask := controllers.CreateTask{
		Title:       &title,
		Description: &description,
	}
	makeRequest("POST", "/tasks", newTask)
	isCompleted := true
	editTask := services.UpdateTaskData{
		IsCompleted: &isCompleted,
	}

	response := makeRequest("PATCH", "/tasks/1", editTask)

	assert.Equal(t, http.StatusOK, response.Code)
}

func TestUpdateEmptyTitleFailed(t *testing.T) {
	title, description := "Test task 1", "Test description 1"
	newTask := controllers.CreateTask{
		Title:       &title,
		Description: &description,
	}
	makeRequest("POST", "/tasks", newTask)
	title = ""
	editTask := services.UpdateTaskData{
		Title: &title,
	}

	response := makeRequest("PATCH", "/tasks/1", editTask)

	assert.Equal(t, http.StatusBadRequest, response.Code)
}
