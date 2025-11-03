package controller

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"testing"

	httpController "github.com/Arup3201/gotasks/internal/controllers/http"
	httperrors "github.com/Arup3201/gotasks/internal/controllers/http/errors"
	entities "github.com/Arup3201/gotasks/internal/entities/task"
	"github.com/Arup3201/gotasks/internal/services"
	"github.com/stretchr/testify/assert"
)

func assertTask(t testing.TB, want, got entities.Task) {
	t.Helper()

	assert.Equal(t, want.Id, got.Id)
	assert.Equal(t, want.Title, got.Title)
	assert.Equal(t, want.Description, got.Description)
	assert.Equal(t, want.IsCompleted, got.IsCompleted)
	assert.Equal(t, want.CreatedAt, got.CreatedAt)
}

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
		t.Fail()
		t.Logf("JSON decoder error: %v", err)
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

	// act
	response1 := makeRequest("POST", "/tasks", task)
	response2 := makeRequest("GET", "/tasks", nil)

	// assert
	assert.Equal(t, http.StatusCreated, response1.Code)

	var responseTask entities.Task
	if err := json.NewDecoder(response1.Body).Decode(&responseTask); err != nil {
		t.Fail()
		t.Logf("JSON decoder error: %v", err)
	}

	assert.Equal(t, task.Title, responseTask.Title)
	assert.Equal(t, task.Description, responseTask.Description)

	assert.Equal(t, http.StatusOK, response2.Code)

	var responseTasks []entities.Task
	if err := json.NewDecoder(response2.Body).Decode(&responseTasks); err != nil {
		t.Fail()
		t.Logf("JSON decoder error: %v", err)
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

	// act
	response1 := makeRequest("POST", "/tasks", task)
	response2 := makeRequest("GET", "/tasks", nil)

	// assert
	assert.Equal(t, http.StatusCreated, response1.Code)

	var responseTask entities.Task
	if err := json.NewDecoder(response1.Body).Decode(&responseTask); err != nil {
		t.Fail()
		t.Logf("JSON decoder error: %v", err)
	}
	assert.Equal(t, task.Title, responseTask.Title)
	assert.Equal(t, task.Description, responseTask.Description)

	assert.Equal(t, http.StatusOK, response2.Code)
	var allTasks []entities.Task
	if err := json.NewDecoder(response2.Body).Decode(&allTasks); err != nil {
		t.Fail()
		t.Logf("JSON decoder error: %v", err)
	}

	assert.Equal(t, 3, len(allTasks))
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
	response := makeRequest("GET", "/tasks/abcd010", nil)

	// assert
	assert.Equal(t, expectedErrorCode, response.Code)

	var responseError httperrors.HttpError
	if err := json.NewDecoder(response.Body).Decode(&responseError); err != nil {
		t.Fail()
		t.Logf("JSON decoder error: %v", err)
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

	var responseError httperrors.HttpError
	if err := json.NewDecoder(response.Body).Decode(&responseError); err != nil {
		t.Fail()
		t.Logf("JSON decoder error: %v", err)
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

	var responseError httperrors.HttpError
	if err := json.NewDecoder(response.Body).Decode(&responseError); err != nil {
		t.Fail()
		t.Logf("JSON decoder error: %v", err)
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

	var responseError httperrors.HttpError
	if err := json.NewDecoder(response.Body).Decode(&responseError); err != nil {
		t.Fail()
		t.Logf("JSON decoder error: %v", err)
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

	var responseError httperrors.HttpError
	if err := json.NewDecoder(response.Body).Decode(&responseError); err != nil {
		t.Fail()
		t.Logf("JSON decoder error: %v", err)
	}

	assert.Equal(t, expectedBody["id"], responseError.Id)
	assert.Equal(t, expectedBody["title"], responseError.Title)
	assert.Equal(t, expectedBody["status"], responseError.Status)
	assert.Equal(t, expectedBody["field"], responseError.Errors[0].Field)
	cleanDB()
}

// success update task title
func TestUpdateTitleSuccess(t *testing.T) {
	// prepare
	tasks := prepareDBTasks(2)
	prepTask := makeRequest("GET", fmt.Sprintf("/tasks/%s", tasks[0].Id), nil)
	var expectedBody entities.Task
	if err := json.NewDecoder(prepTask.Body).Decode(&expectedBody); err != nil {
		t.Fail()
		t.Logf("JSON decoder error: %v", err)
	}
	expectedCode := 200
	title := "Task title (updated)"
	expectedBody.Title = title
	updatePayload := services.UpdateTaskData{
		Title: &title,
	}

	// act
	response := makeRequest("PATCH", fmt.Sprintf("/tasks/%s", tasks[0].Id), updatePayload)

	// assert
	assert.Equal(t, expectedCode, response.Code)

	var responseBody entities.Task
	if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
		t.Fail()
		t.Logf("JSON decode error: %v", err)
	}

	assertTask(t, expectedBody, responseBody)
	cleanDB()
}

// success update task description
func TestUpdateDescriptionSuccess(t *testing.T) {
	// prepare
	tasks := prepareDBTasks(2)
	prepTask := makeRequest("GET", fmt.Sprintf("/tasks/%s", tasks[0].Id), nil)
	var expectedBody entities.Task
	if err := json.NewDecoder(prepTask.Body).Decode(&expectedBody); err != nil {
		t.Fail()
		t.Logf("JSON decoder error: %v", err)
	}
	expectedCode := 200
	description := "Task description (updated)"
	expectedBody.Description = description
	updatePayload := services.UpdateTaskData{
		Description: &description,
	}

	// act
	response := makeRequest("PATCH", fmt.Sprintf("/tasks/%s", tasks[0].Id), updatePayload)

	// assert
	assert.Equal(t, expectedCode, response.Code)

	var responseBody entities.Task
	if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
		t.Fail()
		t.Logf("JSON decode error: %v", err)
	}

	assertTask(t, expectedBody, responseBody)
	cleanDB()
}

// success update task is_completed
func TestUpdateIsCompletedSuccess(t *testing.T) {
	// prepare
	tasks := prepareDBTasks(2)
	prepTask := makeRequest("GET", fmt.Sprintf("/tasks/%s", tasks[0].Id), nil)
	var expectedBody entities.Task
	if err := json.NewDecoder(prepTask.Body).Decode(&expectedBody); err != nil {
		t.Fail()
		t.Logf("JSON decode error: %v", err)
	}
	expectedCode := 200
	isCompleted := true
	expectedBody.IsCompleted = isCompleted
	updatePayload := services.UpdateTaskData{
		IsCompleted: &isCompleted,
	}

	// act
	response := makeRequest("PATCH", fmt.Sprintf("/tasks/%s", tasks[0].Id), updatePayload)

	// assert
	assert.Equal(t, expectedCode, response.Code)

	var responseBody entities.Task
	if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
		t.Fail()
		t.Logf("JSON decode error: %v", err)
	}

	assertTask(t, expectedBody, responseBody)
	cleanDB()
}

// success update task title, description and is_completed
func TestUpdateAll3Success(t *testing.T) {
	// prepare
	tasks := prepareDBTasks(2)
	prepTask := makeRequest("GET", fmt.Sprintf("/tasks/%s", tasks[0].Id), nil)
	var expectedBody entities.Task
	if err := json.NewDecoder(prepTask.Body).Decode(&expectedBody); err != nil {
		t.Fail()
		t.Logf("JSON decode error: %v", err)
	}
	expectedCode := 200
	title := "Task title (updated)"
	description := "Task description (updated)"
	isCompleted := true
	expectedBody.Title = title
	expectedBody.Description = description
	expectedBody.IsCompleted = isCompleted
	updatePayload := services.UpdateTaskData{
		Title:       &title,
		Description: &description,
		IsCompleted: &isCompleted,
	}

	// act
	response := makeRequest("PATCH", fmt.Sprintf("/tasks/%s", tasks[0].Id), updatePayload)

	// assert
	assert.Equal(t, expectedCode, response.Code)

	var responseBody entities.Task
	if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
		t.Fail()
		t.Logf("JSON decode error: %v", err)
	}

	assertTask(t, expectedBody, responseBody)
	cleanDB()
}

// empty payload return no-op response
func TestUpdateNoOpResponse(t *testing.T) {
	// prepare
	tasks := prepareDBTasks(2)
	expectedCode := 204
	expectedBody := map[string]any{
		"id":     "NO_MODIFICATION",
		"title":  "Not modified",
		"status": 204,
	}
	updatePayload := services.UpdateTaskData{}

	// act
	response := makeRequest("PATCH", fmt.Sprintf("/tasks/%s", tasks[0].Id), updatePayload)

	// assert
	assert.Equal(t, expectedCode, response.Code)

	var responseBody httperrors.HttpError
	if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
		t.Fail()
		t.Logf("JSON decode error: %v", err)
	}

	assert.Equal(t, expectedBody["id"], responseBody.Id)
	assert.Equal(t, expectedBody["title"], responseBody.Title)
	assert.Equal(t, expectedBody["status"], responseBody.Status)
	cleanDB()
}

// invalid task id return not found error
func TestUpdateInvalidTaskFail(t *testing.T) {
	// prepare
	prepareDBTasks(2)
	expectedCode := 404
	expectedBody := map[string]any{
		"id":     "NOT_FOUND",
		"title":  "Not found",
		"status": 404,
	}
	title := "Task title (updated)"
	updatePayload := services.UpdateTaskData{
		Title: &title,
	}

	// act
	response := makeRequest("PATCH", "/tasks/abcd1001", updatePayload)

	// assert
	assert.Equal(t, expectedCode, response.Code)

	var responseError httperrors.HttpError
	if err := json.NewDecoder(response.Body).Decode(&responseError); err != nil {
		t.Fail()
		t.Logf("JSON decode error: %v", err)
	}

	assert.Equal(t, expectedBody["id"], responseError.Id)
	assert.Equal(t, expectedBody["title"], responseError.Title)
	assert.Equal(t, expectedBody["status"], responseError.Status)
	cleanDB()
}

// invalid title value fail return bad request error
func TestUpdateInvalidTaskTitleFail(t *testing.T) {
	// prepare
	tasks := prepareDBTasks(2)
	expectedCode := 400
	expectedBody := map[string]any{
		"id":     "INVALID_BODY_PROPERTY",
		"title":  "Invalid body property value",
		"status": 400,
		"field":  "title",
	}
	title := ""
	updatePayload := services.UpdateTaskData{
		Title: &title,
	}

	// act
	response := makeRequest("PATCH", fmt.Sprintf("/tasks/%s", tasks[0].Id), updatePayload)

	// assert
	assert.Equal(t, expectedCode, response.Code)

	var responseError httperrors.HttpError
	if err := json.NewDecoder(response.Body).Decode(&responseError); err != nil {
		t.Fail()
		t.Logf("JSON decode error: %v", err)
	}

	assert.Equal(t, expectedBody["id"], responseError.Id)
	assert.Equal(t, expectedBody["title"], responseError.Title)
	assert.Equal(t, expectedBody["status"], responseError.Status)
	assert.Equal(t, expectedBody["field"], responseError.Errors[0].Field)
	cleanDB()
}

// invalid description value fail return bad request error
func TestUpdateInvalidTaskDescriptionFail(t *testing.T) {
	// prepare
	tasks := prepareDBTasks(2)
	expectedCode := 400
	expectedBody := map[string]any{
		"id":     "INVALID_BODY_PROPERTY",
		"title":  "Invalid body property value",
		"status": 400,
		"field":  "description",
	}
	description := ""
	updatePayload := services.UpdateTaskData{
		Description: &description,
	}

	// act
	response := makeRequest("PATCH", fmt.Sprintf("/tasks/%s", tasks[0].Id), updatePayload)

	// assert
	assert.Equal(t, expectedCode, response.Code)

	var responseError httperrors.HttpError
	if err := json.NewDecoder(response.Body).Decode(&responseError); err != nil {
		t.Fail()
		t.Logf("JSON decode error: %v", err)
	}

	assert.Equal(t, expectedBody["id"], responseError.Id)
	assert.Equal(t, expectedBody["title"], responseError.Title)
	assert.Equal(t, expectedBody["status"], responseError.Status)
	assert.Equal(t, expectedBody["field"], responseError.Errors[0].Field)
	cleanDB()
}

// search tasks return matched single word
func TestSearchSingleWord(t *testing.T) {
	// prepare
	preparedTasks := []struct {
		title       string
		description string
	}{
		{
			title:       "Learn Python",
			description: "Lists comprehension and problem solving",
		},
		{
			title:       "Learn Go",
			description: "Concurrency model",
		},
		{
			title:       "Read story book",
			description: "Feludar sampta kando",
		},
	}
	for _, t := range preparedTasks {
		payload := httpController.CreateTask{
			Title:       &t.title,
			Description: &t.description,
		}
		makeRequest("POST", "/tasks", payload)
	}
	expectedCode := 200
	expectedMatches := []string{
		"Learn Python", "Learn Go",
	}
	term := "learn"

	// act
	url := fmt.Sprintf("/search/tasks?q=%s", term)
	response := makeRequest("GET", url, nil)

	// assert
	assert.Equal(t, expectedCode, response.Code)

	var responseTasks []entities.Task
	if err := json.NewDecoder(response.Body).Decode(&responseTasks); err != nil {
		t.Fail()
		t.Logf("JSON decode error: %v", err)
	}
	assert.Equal(t, len(expectedMatches), len(responseTasks))

	var responseTaskTitles []string
	for _, task := range responseTasks {
		responseTaskTitles = append(responseTaskTitles, task.Title)
	}
	assert.Equal(t, expectedMatches, responseTaskTitles)
	cleanDB()
}

// search tasks return matched multi-word
func TestSearchMultiWord(t *testing.T) {
	// prepare
	preparedTasks := []struct {
		title       string
		description string
	}{
		{
			title:       "Learn DSA with Python",
			description: "Lists comprehension and problem solving",
		},
		{
			title:       "Learn Concurrency with Go",
			description: "Concurrency model",
		},
		{
			title:       "Read story book at afternoon",
			description: "Feludar sampta kando",
		},
		{
			title:       "Read networking book at evening",
			description: "Feludar sampta kando",
		},
	}
	for _, t := range preparedTasks {
		payload := httpController.CreateTask{
			Title:       &t.title,
			Description: &t.description,
		}
		makeRequest("POST", "/tasks", payload)
	}
	expectedCode := 200
	expectedMatches := []string{
		"Read story book at afternoon", "Read networking book at evening",
	}
	term := "read book"

	// act
	url := fmt.Sprintf("/search/tasks?q=%s", term)
	response := makeRequest("GET", url, nil)

	// assert
	assert.Equal(t, expectedCode, response.Code)

	var responseTasks []entities.Task
	if err := json.NewDecoder(response.Body).Decode(&responseTasks); err != nil {
		t.Fail()
		t.Logf("JSON decode error: %v", err)
	}
	assert.Equal(t, len(expectedMatches), len(responseTasks))

	var responseTaskTitles []string
	for _, task := range responseTasks {
		responseTaskTitles = append(responseTaskTitles, task.Title)
	}
	assert.Equal(t, expectedMatches, responseTaskTitles)
	cleanDB()
}

// search returns empty list
func TestSearchReturnEmpty(t *testing.T) {
	// prepare
	preparedTasks := []struct {
		title       string
		description string
	}{
		{
			title:       "Learn DSA with Python",
			description: "Lists comprehension and problem solving",
		},
		{
			title:       "Learn Concurrency with Go",
			description: "Concurrency model",
		},
		{
			title:       "Read story book at afternoon",
			description: "Feludar sampta kando",
		},
		{
			title:       "Read networking book at evening",
			description: "Feludar sampta kando",
		},
	}
	for _, t := range preparedTasks {
		payload := httpController.CreateTask{
			Title:       &t.title,
			Description: &t.description,
		}
		makeRequest("POST", "/tasks", payload)
	}
	expectedCode := 200
	expectedMatches := []string{}
	term := "black magic"

	// act
	url := fmt.Sprintf("/search/tasks?q=%s", term)
	response := makeRequest("GET", url, nil)

	// assert
	assert.Equal(t, expectedCode, response.Code)

	var responseTasks []entities.Task
	if err := json.NewDecoder(response.Body).Decode(&responseTasks); err != nil {
		t.Fail()
		t.Logf("JSON decode error: %v", err)
	}
	assert.Equal(t, len(expectedMatches), len(responseTasks))

	responseTaskTitles := []string{}
	for _, task := range responseTasks {
		responseTaskTitles = append(responseTaskTitles, task.Title)
	}
	assert.Equal(t, expectedMatches, responseTaskTitles)
	cleanDB()
}

// delete success
func TestDeleteTaskSuccess(t *testing.T) {
	// prepare
	tasks := prepareDBTasks(2)
	expectedCode := 200
	expectedResponse := tasks[0].Id

	// act
	url := fmt.Sprintf("/tasks/%s", tasks[0].Id)
	response := makeRequest("DELETE", url, nil)

	// assert
	assert.Equal(t, expectedCode, response.Code)
	assert.Equal(t, expectedResponse, strings.Trim(response.Body.String(), "\""))
	cleanDB()
}

// delete invalid tasks returns not found
func TestDeleteInvalidTaskFail(t *testing.T) {
	// prepare
	prepareDBTasks(2)
	expectedCode := 404
	expectedBody := map[string]any{
		"id":     "NOT_FOUND",
		"title":  "Not found",
		"status": 404,
	}

	// act
	response := makeRequest("DELETE", "/tasks/asdf", nil)

	// assert
	assert.Equal(t, expectedCode, response.Code)

	var responseError httperrors.HttpError
	if err := json.NewDecoder(response.Body).Decode(&responseError); err != nil {
		t.Fail()
		t.Logf("JSON decode error: %v", err)
	}

	assert.Equal(t, expectedBody["id"], responseError.Id)
	assert.Equal(t, expectedBody["title"], responseError.Title)
	assert.Equal(t, expectedBody["status"], responseError.Status)
	cleanDB()
}
