package http

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	entities "github.com/Arup3201/gotasks/internal/entities/task"
	services "github.com/Arup3201/gotasks/internal/services/domain/task"
	"github.com/gin-gonic/gin"
)

func getTestContext(t testing.TB, w http.ResponseWriter, r *http.Request) *gin.Context {
	t.Helper()

	gin.SetMode(gin.TestMode)

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = r

	return ctx
}

func TestGetTasks(t *testing.T) {
	t.Run("get all tasks", func(t *testing.T) {
		repo := &MockRepository{
			tasks: []entities.Task{
				{
					Id:          1,
					Title:       "Task 1",
					Description: "No description",
					IsCompleted: false,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
				{
					Id:          2,
					Title:       "Task 2",
					Description: "No description",
					IsCompleted: true,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
		}
		serviceHandler, _ := services.NewTaskService(repo)
		routeHandler := getRouteHandler(serviceHandler)
		request, _ := http.NewRequest("GET", "/tasks", nil)
		response := httptest.NewRecorder()
		ctx := getTestContext(t, response, request)

		routeHandler.GetTasks(ctx)

		var got []entities.Task
		want := 2
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			log.Fatal("JSON decoding failed")
		}
		if len(got) != want {
			t.Errorf("response is wrong, expected %d tasks but got %d tasks", want, len(got))
		}
	})
	t.Run("get no tasks", func(t *testing.T) {
		repo := &MockRepository{
			tasks: []entities.Task{},
		}
		serviceHandler, _ := services.NewTaskService(repo)
		routeHandler := getRouteHandler(serviceHandler)
		request, _ := http.NewRequest("GET", "/tasks", nil)
		response := httptest.NewRecorder()
		ctx := getTestContext(t, response, request)

		routeHandler.GetTasks(ctx)

		var got []entities.Task
		want := 0
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			log.Fatal("JSON decoding failed")
		}
		if len(got) != want {
			t.Errorf("response is wrong, expected %d tasks but got %d tasks", want, len(got))
		}
	})
}

func TestAddTask(t *testing.T) {
	t.Run("add a task success", func(t *testing.T) {
		repo := &MockRepository{
			tasks: []entities.Task{},
		}
		serviceHandler, _ := services.NewTaskService(repo)
		routeHandler := getRouteHandler(serviceHandler)
		payload := strings.NewReader(`{
			"title": "Test task", 
			"description": "Test description"
		}`)
		request, _ := http.NewRequest("POST", "/tasks", payload)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		ctx := getTestContext(t, response, request)

		routeHandler.AddTask(ctx)

		var got entities.Task
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			log.Fatal("JSON decoding failed")
		}
		want := 1
		if got.Id != want {
			t.Errorf("task ID expected %d but got %d", want, got.Id)
		}
		title := "Test task"
		if got.Title != title {
			t.Errorf("task title expected %s but got %s", title, got.Title)
		}
		description := "Test description"
		if got.Description != description {
			t.Errorf("task description expected %s but got %s", description, got.Description)
		}
		if got.CreatedAt.IsZero() || got.UpdatedAt.IsZero() {
			t.Errorf("task should have created_at or updated_at")
		}
	})
	t.Run("add task increases ID", func(t *testing.T) {
		repo := &MockRepository{
			tasks: []entities.Task{
				{
					Id:          1,
					Title:       "Task 1",
					Description: "No description",
					IsCompleted: false,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
		}
		serviceHandler, _ := services.NewTaskService(repo)
		routeHandler := getRouteHandler(serviceHandler)
		payload := strings.NewReader(`{
			"title": "Test task", 
			"description": "Test description"
		}`)
		request, _ := http.NewRequest("POST", "/tasks", payload)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		ctx := getTestContext(t, response, request)

		routeHandler.AddTask(ctx)

		var got entities.Task
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			log.Fatal("JSON decoding failed")
		}
		want := 2
		if got.Id != want {
			t.Errorf("expected ID to be %d, but got %d", want, got.Id)
		}
	})
	t.Run("add task fail for missing title", func(t *testing.T) {
		repo := &MockRepository{
			tasks: []entities.Task{},
		}
		serviceHandler, _ := services.NewTaskService(repo)
		routeHandler := getRouteHandler(serviceHandler)
		payload := strings.NewReader(`{
			"description": "Test description"
		}`)
		request, _ := http.NewRequest("POST", "/tasks", payload)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		ctx := getTestContext(t, response, request)

		routeHandler.AddTask(ctx)

		want := BadRequest
		if got := response.Result().StatusCode; got != want {
			t.Errorf("expected status code %d but got %d", want, got)
		}
	})
}
