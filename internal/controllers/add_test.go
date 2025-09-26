package controller

import (
	"net/http"
	"testing"

	controllers "github.com/Arup3201/gotasks/internal/controllers/http"
	"github.com/stretchr/testify/assert"
)

func TestAddSuccess(t *testing.T) {
	title, description := "Test task 1", "Test description 1"
	newTask := controllers.CreateTask{
		Title:       &title,
		Description: &description,
	}

	response := makeRequest("POST", "/tasks", newTask)
	assert.Equal(t, http.StatusCreated, response.Code)
}

func TestAddFailMissingTitle(t *testing.T) {
	description := "Test description 2"
	newTask := controllers.CreateTask{
		Description: &description,
	}

	response := makeRequest("POST", "/tasks", newTask)
	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestAddFailMissingDescription(t *testing.T) {
	title := "Test task 3"
	newTask := controllers.CreateTask{
		Title: &title,
	}

	response := makeRequest("POST", "/tasks", newTask)
	assert.Equal(t, http.StatusBadRequest, response.Code)
}
