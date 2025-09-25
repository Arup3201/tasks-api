package httpController

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Arup3201/gotasks/internal/controllers/http/middlewares"
	entities "github.com/Arup3201/gotasks/internal/entities/task"
	services "github.com/Arup3201/gotasks/internal/services/domain/task"
	"github.com/gin-gonic/gin"
)

func getTestContext(t testing.TB, w http.ResponseWriter, r *http.Request) (*gin.Context, *gin.Engine) {
	t.Helper()

	gin.SetMode(gin.TestMode)

	ctx, engine := gin.CreateTestContext(w)
	ctx.Request = r

	return ctx, engine
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
		routeHandler := GetRouteHandler(serviceHandler)
		request, _ := http.NewRequest("GET", "/tasks", nil)
		response := httptest.NewRecorder()
		ctx, engine := getTestContext(t, response, request)
		engine.GET("/tasks", routeHandler.GetTasks)

		engine.ServeHTTP(response, ctx.Request)

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
		routeHandler := GetRouteHandler(serviceHandler)
		request, _ := http.NewRequest("GET", "/tasks", nil)
		response := httptest.NewRecorder()
		ctx, engine := getTestContext(t, response, request)
		engine.GET("/tasks", routeHandler.GetTasks)

		engine.ServeHTTP(response, ctx.Request)

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
		routeHandler := GetRouteHandler(serviceHandler)
		payload := strings.NewReader(`{
			"title": "Test task", 
			"description": "Test description"
		}`)
		request, _ := http.NewRequest("POST", "/tasks", payload)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		ctx, engine := getTestContext(t, response, request)
		engine.POST("/tasks", routeHandler.AddTask)

		engine.ServeHTTP(response, ctx.Request)

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
		routeHandler := GetRouteHandler(serviceHandler)
		payload := strings.NewReader(`{
			"title": "Test task", 
			"description": "Test description"
		}`)
		request, _ := http.NewRequest("POST", "/tasks", payload)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		ctx, engine := getTestContext(t, response, request)
		engine.POST("/tasks", routeHandler.AddTask)

		engine.ServeHTTP(response, ctx.Request)

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
		routeHandler := GetRouteHandler(serviceHandler)
		payload := strings.NewReader(`{
			"description": "Test description"
		}`)
		request, _ := http.NewRequest("POST", "/tasks", payload)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		ctx, engine := getTestContext(t, response, request)
		engine.Use(middlewares.HttpErrorResponse())
		engine.POST("/tasks", routeHandler.AddTask)

		engine.ServeHTTP(response, ctx.Request)

		want := BadRequest
		if got := response.Result().StatusCode; got != want {
			t.Errorf("expected status code %d but got %d", want, got)
		}
	})
}

func TestGetTask(t *testing.T) {
	t.Run("get task success", func(t *testing.T) {
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
		routeHandler := GetRouteHandler(serviceHandler)
		request, _ := http.NewRequest("GET", "/tasks/2", nil)
		response := httptest.NewRecorder()
		ctx, engine := getTestContext(t, response, request)
		engine.GET("/tasks/:id", routeHandler.GetTask)

		engine.ServeHTTP(response, ctx.Request)

		var got entities.Task
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			log.Fatal("JSON decoding failed")
		}
		want := 2
		if got.Id != want {
			t.Errorf("response is wrong, expected task with ID %d but got %d tasks", want, got.Id)
		}
	})
	t.Run("get task success", func(t *testing.T) {
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
		routeHandler := GetRouteHandler(serviceHandler)
		request, _ := http.NewRequest("GET", "/tasks/3", nil)
		response := httptest.NewRecorder()
		ctx, engine := getTestContext(t, response, request)
		engine.Use(middlewares.HttpErrorResponse())
		engine.GET("/tasks/:id", routeHandler.GetTask)

		engine.ServeHTTP(response, ctx.Request)

		want := NotFound
		if got := response.Result().StatusCode; got != want {
			t.Errorf("should return NotFound error, expected status code %d but got %d", want, got)
		}
	})
}

func TestUpdateTask(t *testing.T) {
	t.Run("update task with new title", func(t *testing.T) {
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
		routeHandler := GetRouteHandler(serviceHandler)
		payload := strings.NewReader(`{
			"title": "Test 2 (edited)"
		}`)
		request, _ := http.NewRequest("PATCH", "/tasks/2", payload)
		response := httptest.NewRecorder()
		ctx, engine := getTestContext(t, response, request)
		engine.PATCH("/tasks/:id", routeHandler.UpdateTask)

		engine.ServeHTTP(response, ctx.Request)

		var got entities.Task
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			log.Fatal("JSON decoding failed")
		}
		want := "Test 2 (edited)"
		if got.Title != want {
			t.Errorf("expected title after update %s, but got %s", want, got.Title)
		}
	})
	t.Run("update task with new description", func(t *testing.T) {
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
		routeHandler := GetRouteHandler(serviceHandler)
		payload := strings.NewReader(`{
			"description": "Test 2 description (edited)"
		}`)
		request, _ := http.NewRequest("PATCH", "/tasks/2", payload)
		response := httptest.NewRecorder()
		ctx, engine := getTestContext(t, response, request)
		engine.PATCH("/tasks/:id", routeHandler.UpdateTask)

		engine.ServeHTTP(response, ctx.Request)

		var result entities.Task
		err := json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			log.Fatal("JSON decoding failed")
		}
		want := "Test 2 description (edited)"
		if got := result.Description; got != want {
			t.Errorf("expected description after update %s, but got %s", want, got)
		}
	})
	t.Run("update task with new is_completed", func(t *testing.T) {
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
		routeHandler := GetRouteHandler(serviceHandler)
		payload := strings.NewReader(`{
			"is_completed": true
		}`)
		request, _ := http.NewRequest("PATCH", "/tasks/1", payload)
		response := httptest.NewRecorder()
		ctx, engine := getTestContext(t, response, request)
		engine.PATCH("/tasks/:id", routeHandler.UpdateTask)

		engine.ServeHTTP(response, ctx.Request)

		var result entities.Task
		err := json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			log.Fatal("JSON decoding failed")
		}
		want := true
		if got := result.IsCompleted; got != want {
			t.Errorf("expected description after update %t, but got %t", want, got)
		}
	})
	t.Run("update task fail not found", func(t *testing.T) {
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
		routeHandler := GetRouteHandler(serviceHandler)
		payload := strings.NewReader(`{
			"title": "Test 3 (edited)"
		}`)
		request, _ := http.NewRequest("PATCH", "/tasks/3", payload)
		response := httptest.NewRecorder()
		ctx, engine := getTestContext(t, response, request)
		engine.Use(middlewares.HttpErrorResponse())
		engine.PATCH("/tasks/:id", routeHandler.UpdateTask)

		engine.ServeHTTP(response, ctx.Request)

		want := NotFound
		if got := response.Result().StatusCode; got != want {
			t.Errorf("expected NotFound error, expected status code %d but got %d", want, got)
		}
	})
	t.Run("update task with no op", func(t *testing.T) {
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
		routeHandler := GetRouteHandler(serviceHandler)
		payload := strings.NewReader(`{}`)
		request, _ := http.NewRequest("PATCH", "/tasks/3", payload)
		response := httptest.NewRecorder()
		ctx, engine := getTestContext(t, response, request)
		engine.Use(middlewares.HttpErrorResponse())
		engine.PATCH("/tasks/:id", routeHandler.UpdateTask)

		engine.ServeHTTP(response, ctx.Request)

		want := NoOp
		if got := response.Result().StatusCode; got != want {
			t.Errorf("expected NoOp response, expected status code %d but got %d", want, got)
		}
	})
}

func TestDeleteTask(t *testing.T) {
	t.Run("delete task success", func(t *testing.T) {
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
		routeHandler := GetRouteHandler(serviceHandler)
		request, _ := http.NewRequest("DELETE", "/tasks/2", nil)
		response := httptest.NewRecorder()
		ctx, engine := getTestContext(t, response, request)
		engine.DELETE("/tasks/:id", routeHandler.DeleteTask)

		engine.ServeHTTP(response, ctx.Request)

		want := 200
		if got := response.Result().StatusCode; got != want {
			t.Errorf("delete failed, expected status code %d but got %d", want, got)
		}
		id := "2"
		if got := response.Body.String(); got != id {
			t.Errorf("delete did return ID, expected %s but got %s", id, got)
		}
	})
	t.Run("delete task fail not found", func(t *testing.T) {
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
		routeHandler := GetRouteHandler(serviceHandler)
		request, _ := http.NewRequest("DELETE", "/tasks/3", nil)
		response := httptest.NewRecorder()
		ctx, engine := getTestContext(t, response, request)
		engine.Use(middlewares.HttpErrorResponse())
		engine.DELETE("/tasks/:id", routeHandler.DeleteTask)

		engine.ServeHTTP(response, ctx.Request)

		want := NotFound
		if got := response.Result().StatusCode; got != want {
			t.Errorf("expected NotFound error %d, but got %d", want, got)
		}
	})
}

func TestSearchTasks(t *testing.T) {
	t.Run("search tasks success", func(t *testing.T) {
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
				{
					Id:          3,
					Title:       "Task 2",
					Description: "Task 2 additional",
					IsCompleted: false,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
		}
		serviceHandler, _ := services.NewTaskService(repo)
		routeHandler := GetRouteHandler(serviceHandler)
		request, _ := http.NewRequest("GET", "/search/tasks?q=2", nil)
		response := httptest.NewRecorder()
		ctx, engine := getTestContext(t, response, request)

		engine.GET("/search/tasks", routeHandler.SearchTasks)
		engine.ServeHTTP(response, ctx.Request)

		var got []entities.Task
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			log.Fatal("JSON decoding failed")
		}
		want := 2
		if len(got) != want {
			t.Errorf("response is wrong, expected %d searched tasks but got %d tasks", want, len(got))
		}
	})
	t.Run("search tasks no task found", func(t *testing.T) {
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
				{
					Id:          3,
					Title:       "Task 2",
					Description: "Task 2 additional",
					IsCompleted: false,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
		}
		serviceHandler, _ := services.NewTaskService(repo)
		routeHandler := GetRouteHandler(serviceHandler)
		request, _ := http.NewRequest("GET", "/search/tasks?q=3", nil)
		response := httptest.NewRecorder()
		ctx, engine := getTestContext(t, response, request)

		engine.GET("/search/tasks", routeHandler.SearchTasks)
		engine.ServeHTTP(response, ctx.Request)

		var got []entities.Task
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			log.Fatal("JSON decoding failed")
		}
		want := 0
		if len(got) != want {
			t.Errorf("response is wrong, expected %d searched tasks but got %d tasks", want, len(got))
		}
	})
}
