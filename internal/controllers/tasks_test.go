package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"testing"

	entities "github.com/Arup3201/gotasks/internal/entities/task"
	"github.com/stretchr/testify/assert"
)

// Check that we have 10 tasks initially
func TestViewAllTasksSuccess(t *testing.T) {
	// prepare
	prepareDBTasks(2)
	expectedTasksNum := 2

	// act
	response := makeRequest("GET", "/tasks", nil)

	// assert
	assert.Equal(t, http.StatusOK, response.Code)

	tasks := []entities.Task{}
	if err := json.NewDecoder(response.Body).Decode(&tasks); err != nil {
		log.Fatalf("JSON decoder error: %v", err)
	}
	assert.Equal(t, expectedTasksNum, len(tasks))
	cleanDB()
}

// check tasks increase to 11 after adding a new task
func TestAddTaskCheck(t *testing.T) {
	// prepare
	prepareDBTasks(2)
	task := entities.Task{
		Title:       fmt.Sprintf("title - %d", rand.Intn(9999)),
		Description: fmt.Sprintf("description - %d", rand.Intn(9999)),
	}
	expectedTasksNum := 3
	expectedTaskId := 3

	// act
	response1 := makeRequest("POST", "/tasks", task)
	response2 := makeRequest("GET", "/tasks", nil)

	// assert
	assert.Equal(t, http.StatusCreated, response1.Code)

	var responseTask entities.Task
	if err := json.NewDecoder(response1.Body).Decode(&responseTask); err != nil {
		log.Fatalf("JSON decoder error: %v", err)
	}

	assert.Equal(t, expectedTaskId, responseTask.Id)

	assert.Equal(t, http.StatusOK, response2.Code)

	var responseTasks []entities.Task
	if err := json.NewDecoder(response2.Body).Decode(&responseTasks); err != nil {
		log.Fatalf("JSON decoder error: %v", err)
	}
	assert.Equal(t, expectedTasksNum, len(responseTasks))

	cleanDB()
}

func TestAddAndViewTask(t *testing.T) {
	// prepare
	prepareDBTasks(2)
	task := entities.Task{
		Title:       fmt.Sprintf("title - %d", rand.Intn(9999)),
		Description: fmt.Sprintf("description - %d", rand.Intn(9999)),
	}
	expected := entities.Task{
		Id:          3,
		Title:       task.Title,
		Description: task.Description,
	}

	// act
	response1 := makeRequest("POST", "/tasks", task)
	response2 := makeRequest("GET", "/tasks/3", nil)

	// assert
	assert.Equal(t, http.StatusCreated, response1.Code)

	var responseTask entities.Task
	if err := json.NewDecoder(response1.Body).Decode(&responseTask); err != nil {
		log.Fatalf("JSON decoder error: %v", err)
	}
	assert.Equal(t, expected.Id, responseTask.Id)
	assert.Equal(t, expected.Title, responseTask.Title)
	assert.Equal(t, expected.Description, responseTask.Description)

	assert.Equal(t, http.StatusOK, response2.Code)
	if err := json.NewDecoder(response2.Body).Decode(&responseTask); err != nil {
		log.Fatalf("JSON decoder error: %v", err)
	}

	assert.Equal(t, responseTask.IsCompleted, false)
	assert.NotZero(t, responseTask.CreatedAt)
	assert.NotZero(t, responseTask.UpdatedAt)

	cleanDB()
}
