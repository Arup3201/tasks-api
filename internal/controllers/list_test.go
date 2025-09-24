package controller

import (
	"net/http"
	"testing"

	controllers "github.com/Arup3201/gotasks/internal/controllers/http"
	"github.com/stretchr/testify/assert"
)

func TestListTasks(t *testing.T) {
	testCases := []struct {
		title       string
		description string
	}{
		{
			title:       "Learn Golang",
			description: "Cover concurrency topic",
		},
		{
			title:       "Learn Python",
			description: "Cover asynchronous topic",
		},
		{
			title:       "Learn C#",
			description: "Cover .NET based REST API development",
		},
	}
	for _, tc := range testCases {
		newTask := controllers.CreateTask{
			Title:       &tc.title,
			Description: &tc.description,
		}
		makeRequest("POST", "/tasks", newTask)
	}

	response := makeRequest("GET", "/tasks", nil)

	assert.Equal(t, http.StatusOK, response.Code)
}
