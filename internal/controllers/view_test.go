package controller

import (
	"net/http"
	"testing"

	controllers "github.com/Arup3201/gotasks/internal/controllers/http"
	"github.com/stretchr/testify/assert"
)

func TestViewTaskSuccess(t *testing.T) {
	title, description := "Test task 1", "Test description 1"
	newTask := controllers.CreateTask{
		Title:       &title,
		Description: &description,
	}
	makeRequest("POST", "/tasks", newTask)

	response := makeRequest("GET", "/tasks/1", nil)

	assert.Equal(t, http.StatusOK, response.Code)
}

func TestViewTaskFail(t *testing.T) {

	response := makeRequest("GET", "/tasks/10", nil)

	assert.Equal(t, http.StatusNotFound, response.Code)
}
