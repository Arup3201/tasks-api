package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"testing"

	httpController "github.com/Arup3201/gotasks/internal/controllers/http"
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
func TestAddTaskCheckSuccess(t *testing.T) {
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

// check after adding a new task to the list of 2 tasks, the new task id is 3 and the task has correct content
func TestAddAndViewTaskSuccess(t *testing.T) {
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

// check if I try to access wrong task then I get correct NotFound response body
func TestViewInvalidTaskFail(t *testing.T) {
	// prepare
	prepareDBTasks(2)
	expectedErrorCode := 404
	expectedBody := map[string]any{
		"id":     "NOT_FOUND",
		"title":  "Not found",
		"status": 404,
	}

	// act
	response := makeRequest("GET", "/tasks/3", nil)

	// assert
	assert.Equal(t, expectedErrorCode, response.Code)

	var responseError httpController.HttpError
	if err := json.NewDecoder(response.Body).Decode(&responseError); err != nil {
		log.Fatalf("JSON decoder error: %v", err)
	}

	assert.Equal(t, expectedBody["id"], responseError.Id)
	assert.Equal(t, expectedBody["title"], responseError.Title)
	assert.Equal(t, expectedBody["status"], responseError.Status)
	cleanDB()
}

// check adding a task with missing title and receiving the correct bad request response error body
func TestAddTaskWithMissingTitleFail(t *testing.T) {
	// prepare
	prepareDBTasks(2)
	description := fmt.Sprintf("description - %d", rand.Intn(9999))
	task := httpController.CreateTask{
		Description: &description,
	}
	expectedErrorCode := 400
	expectedBody := map[string]any{
		"id":     "MISSING_BODY_PROPERTY",
		"title":  "Missing body property",
		"status": 400,
		"field":  "title",
	}

	// act
	response := makeRequest("POST", "/tasks", task)

	// assert
	assert.Equal(t, expectedErrorCode, response.Code)

	var responseError httpController.HttpError
	if err := json.NewDecoder(response.Body).Decode(&responseError); err != nil {
		log.Fatalf("JSON decoder error: %v", err)
	}

	assert.Equal(t, expectedBody["id"], responseError.Id)
	assert.Equal(t, expectedBody["title"], responseError.Title)
	assert.Equal(t, expectedBody["status"], responseError.Status)
	assert.Equal(t, expectedBody["field"], responseError.Errors[0].Field)
	cleanDB()
}

// check adding a task with missing description and receiving the correct bad request response error body
func TestAddTaskWithMissingDescriptionFail(t *testing.T) {
	// prepare
	prepareDBTasks(2)
	title := fmt.Sprintf("title - %d", rand.Intn(9999))
	task := httpController.CreateTask{
		Title: &title,
	}
	expectedErrorCode := 400
	expectedBody := map[string]any{
		"id":     "MISSING_BODY_PROPERTY",
		"title":  "Missing body property",
		"status": 400,
		"field":  "description",
	}

	// act
	response := makeRequest("POST", "/tasks", task)

	// assert
	assert.Equal(t, expectedErrorCode, response.Code)

	var responseError httpController.HttpError
	if err := json.NewDecoder(response.Body).Decode(&responseError); err != nil {
		log.Fatalf("JSON decoder error: %v", err)
	}

	assert.Equal(t, expectedBody["id"], responseError.Id)
	assert.Equal(t, expectedBody["title"], responseError.Title)
	assert.Equal(t, expectedBody["status"], responseError.Status)
	assert.Equal(t, expectedBody["field"], responseError.Errors[0].Field)
	cleanDB()
}

// check adding a task with empty title and receiving the correct bad request response error body
func TestAddTaskEmptyTitleFail(t *testing.T) {
	// prepare
	prepareDBTasks(2)
	title := ""
	description := fmt.Sprintf("description - %d", rand.Intn(9999))
	task := httpController.CreateTask{
		Title:       &title,
		Description: &description,
	}
	expectedErrorCode := 400
	expectedBody := map[string]any{
		"id":     "INVALID_BODY_PROPERTY",
		"title":  "Invalid body property value",
		"status": 400,
		"field":  "title",
	}

	// act
	response := makeRequest("POST", "/tasks", task)

	// assert
	assert.Equal(t, expectedErrorCode, response.Code)

	var responseError httpController.HttpError
	if err := json.NewDecoder(response.Body).Decode(&responseError); err != nil {
		log.Fatalf("JSON decoder error: %v", err)
	}

	assert.Equal(t, expectedBody["id"], responseError.Id)
	assert.Equal(t, expectedBody["title"], responseError.Title)
	assert.Equal(t, expectedBody["status"], responseError.Status)
	assert.Equal(t, expectedBody["field"], responseError.Errors[0].Field)
	cleanDB()
}

// check adding a task with empty description and receiving the correct bad request response error body
func TestAddTaskEmptyDescriptionFail(t *testing.T) {
	// prepare
	prepareDBTasks(2)
	title := fmt.Sprintf("title - %d", rand.Intn(9999))
	description := ""
	task := httpController.CreateTask{
		Title:       &title,
		Description: &description,
	}
	expectedErrorCode := 400
	expectedBody := map[string]any{
		"id":     "INVALID_BODY_PROPERTY",
		"title":  "Invalid body property value",
		"status": 400,
		"field":  "description",
	}

	// act
	response := makeRequest("POST", "/tasks", task)

	// assert
	assert.Equal(t, expectedErrorCode, response.Code)

	var responseError httpController.HttpError
	if err := json.NewDecoder(response.Body).Decode(&responseError); err != nil {
		log.Fatalf("JSON decoder error: %v", err)
	}

	assert.Equal(t, expectedBody["id"], responseError.Id)
	assert.Equal(t, expectedBody["title"], responseError.Title)
	assert.Equal(t, expectedBody["status"], responseError.Status)
	assert.Equal(t, expectedBody["field"], responseError.Errors[0].Field)
	cleanDB()
}

// check if non-integer id gives invalid param error response
func TestViewIdParamInvalidFail(t *testing.T) {
	// prepare
	prepareDBTasks(2)
	expectedErrorCode := 400
	expectedBody := map[string]any{
		"id":     "INVALID_PARAMETER_VALUE",
		"title":  "Invalid request parameter value",
		"status": 400,
		"field":  "id",
	}

	// act
	response := makeRequest("GET", "/tasks/asd", nil)

	// assert
	assert.Equal(t, expectedErrorCode, response.Code)

	var responseError httpController.HttpError
	if err := json.NewDecoder(response.Body).Decode(&responseError); err != nil {
		log.Fatalf("JSON decoder error: %v", err)
	}

	assert.Equal(t, expectedBody["id"], responseError.Id)
	assert.Equal(t, expectedBody["title"], responseError.Title)
	assert.Equal(t, expectedBody["status"], responseError.Status)
	assert.Equal(t, expectedBody["field"], responseError.Errors[0].Field)
	cleanDB()
}
