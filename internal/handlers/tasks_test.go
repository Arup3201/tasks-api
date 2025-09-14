package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/Arup3201/gotasks/internal/handlers/apperr"
	"github.com/Arup3201/gotasks/internal/models"
	"github.com/gin-gonic/gin"
)

type mockStorage struct {
	lastId int
	tasks  map[string]models.Task
}

func NewMockStorage(tasks map[string]models.Task) *mockStorage {
	if len(tasks) > 0 {
		return &mockStorage{
			tasks:  tasks,
			lastId: len(tasks),
		}
	} else {
		return &mockStorage{tasks: make(map[string]models.Task), lastId: 1}
	}
}

func (m *mockStorage) Get(id string) (models.Task, error) {
	task, ok := m.tasks[id]
	if !ok {
		return task, apperr.NotFoundError()
	}

	return task, nil
}
func (m *mockStorage) Insert(title, description string) (string, error) {
	id := strconv.Itoa(m.lastId + 1)
	m.tasks[id] = models.Task{
		Id:          id,
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	m.lastId += 1
	return id, nil
}
func (m *mockStorage) Update(id string, title, description *string, completed *bool) (models.Task, error) {
	task, ok := m.tasks[id]
	if !ok {
		return task, apperr.NotFoundError()
	}

	if title != nil {
		task.Title = *title
	}
	if description != nil {
		task.Description = *description
	}
	if completed != nil {
		task.Completed = *completed
	}

	return task, nil
}
func (m *mockStorage) Delete(id string) (string, error) {
	delete(m.tasks, id)

	return id, nil
}
func (m *mockStorage) List() ([]models.Task, error) {
	tasks := []models.Task{}
	for _, task := range m.tasks {
		tasks = append(tasks, task)
	}

	return tasks, nil
}
func (m *mockStorage) Search(by models.FieldName, query string) ([]models.Task, error) {
	return nil, nil
}

func getTestRouter(t testing.TB) *gin.Engine {
	t.Helper()

	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(HandleErrors())

	return router
}

func TestGETAllTasks(t *testing.T) {
	t.Run("returns 2 tasks", func(t *testing.T) {
		mock := NewMockStorage(map[string]models.Task{
			"1": {
				Id:          "1",
				Title:       "Task 1",
				Description: "No description",
				Completed:   false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			"2": {
				Id:          "2",
				Title:       "Task 2",
				Description: "No description",
				Completed:   true,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		})
		tHandler := NewTaskHandler(mock)
		router := getTestRouter(t)
		router.GET("/tasks", tHandler.GetAllTasks)
		request, _ := http.NewRequest("GET", "/tasks", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		var got []models.Task
		want := 2
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			log.Fatal("JSON decoding failed")
		}
		if len(got) != want {
			t.Errorf("response is wrong, expected %d tasks but got %d tasks", want, len(got))
		}
	})
	t.Run("returns no tasks", func(t *testing.T) {
		mock := NewMockStorage(map[string]models.Task{})
		tHandler := NewTaskHandler(mock)
		router := getTestRouter(t)
		router.GET("/tasks", tHandler.GetAllTasks)
		request, _ := http.NewRequest("GET", "/tasks", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		var got []models.Task
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

func TestGetTaskWithId(t *testing.T) {
	t.Run("returns task with id 1", func(t *testing.T) {
		mock := NewMockStorage(map[string]models.Task{
			"1": {
				Id:          "1",
				Title:       "Task 1",
				Description: "No description",
				Completed:   false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			"2": {
				Id:          "2",
				Title:       "Task 2",
				Description: "No description",
				Completed:   true,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		})
		tHandler := NewTaskHandler(mock)
		router := getTestRouter(t)
		router.GET("/tasks/:id", tHandler.GetTaskWithId)
		request, _ := http.NewRequest("GET", "/tasks/1", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		want := 200
		if got := response.Result().StatusCode; got != want {
			t.Errorf("response status is wrong, expected %d but got %d", want, got)
		}
	})
	t.Run("returns task not found", func(t *testing.T) {
		mock := NewMockStorage(map[string]models.Task{})
		tHandler := NewTaskHandler(mock)
		router := getTestRouter(t)
		router.GET("/tasks/:id", tHandler.GetTaskWithId)
		request, _ := http.NewRequest("GET", "/tasks/1", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		want := 404
		if got := response.Result().StatusCode; got != want {
			t.Errorf("response status is wrong, expected %d but got %d", want, got)
		}
	})
}

func TestAddTask(t *testing.T) {
	t.Run("add task with valid payload", func(t *testing.T) {
		mock := NewMockStorage(map[string]models.Task{})
		tHandler := NewTaskHandler(mock)
		router := getTestRouter(t)
		router.POST("/tasks", tHandler.AddTask)
		title, description := "Task adding", "Added desc"
		task := models.CreateTask{
			Title:       &title,
			Description: &description,
		}
		marshalled, err := json.Marshal(task)
		if err != nil {
			t.Fatalf("marshalling failed with error: %v", err)
		}
		request, _ := http.NewRequest("POST", "/tasks", bytes.NewReader(marshalled))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		want := 201
		if got := response.Result().StatusCode; got != want {
			t.Log(response.Body.String())
			t.Errorf("response status is wrong, expected %d but got %d", want, got)
		}
	})
	t.Run("add task with missing title", func(t *testing.T) {
		mock := NewMockStorage(map[string]models.Task{})
		tHandler := NewTaskHandler(mock)
		router := getTestRouter(t)
		router.POST("/tasks", tHandler.AddTask)
		description := "Added desc"
		task := models.CreateTask{
			Description: &description,
		}
		marshalled, err := json.Marshal(task)
		if err != nil {
			t.Fatalf("marshalling failed with error: %v", err)
		}
		request, _ := http.NewRequest("POST", "/tasks", bytes.NewReader(marshalled))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		want := 400
		if got := response.Result().StatusCode; got != want {
			t.Errorf("response status is wrong, expected %d but got %d", want, got)
		}
	})
	t.Run("add task with missing description", func(t *testing.T) {
		mock := NewMockStorage(map[string]models.Task{})
		tHandler := NewTaskHandler(mock)
		router := getTestRouter(t)
		router.POST("/tasks", tHandler.AddTask)
		title := "Task adding"
		task := models.CreateTask{
			Title: &title,
		}
		marshalled, err := json.Marshal(task)
		if err != nil {
			t.Fatalf("marshalling failed with error: %v", err)
		}
		request, _ := http.NewRequest("POST", "/tasks", bytes.NewReader(marshalled))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		want := 400
		if got := response.Result().StatusCode; got != want {
			t.Errorf("response status is wrong, expected %d but got %d", want, got)
		}
	})
	t.Run("add task with empty title", func(t *testing.T) {
		mock := NewMockStorage(map[string]models.Task{})
		tHandler := NewTaskHandler(mock)
		router := getTestRouter(t)
		router.POST("/tasks", tHandler.AddTask)
		title, description := "", "Added desc"
		task := models.CreateTask{
			Title:       &title,
			Description: &description,
		}
		marshalled, err := json.Marshal(task)
		if err != nil {
			t.Fatalf("marshalling failed with error: %v", err)
		}
		request, _ := http.NewRequest("POST", "/tasks", bytes.NewReader(marshalled))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		want := 400
		if got := response.Result().StatusCode; got != want {
			t.Errorf("response status is wrong, expected %d but got %d", want, got)
		}
	})
	t.Run("add task with empty description", func(t *testing.T) {
		mock := NewMockStorage(map[string]models.Task{})
		tHandler := NewTaskHandler(mock)
		router := getTestRouter(t)
		router.POST("/tasks", tHandler.AddTask)
		title, description := "Task adding", ""
		task := models.CreateTask{
			Title:       &title,
			Description: &description,
		}
		marshalled, err := json.Marshal(task)
		if err != nil {
			t.Fatalf("marshalling failed with error: %v", err)
		}
		request, _ := http.NewRequest("POST", "/tasks", bytes.NewReader(marshalled))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		want := 400
		if got := response.Result().StatusCode; got != want {
			t.Errorf("response status is wrong, expected %d but got %d", want, got)
		}
	})
}

func TestEditTask(t *testing.T) {
	t.Run("edit task with valid payload", func(t *testing.T) {
		mock := NewMockStorage(map[string]models.Task{
			"1": {
				Id:          "1",
				Title:       "Task 1",
				Description: "No description",
				Completed:   false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		})
		tHandler := NewTaskHandler(mock)
		router := getTestRouter(t)
		router.PUT("/tasks/:id", tHandler.EditTask)
		title := "Task 1 (edited)"
		task := models.EditTask{
			Title: &title,
		}
		marshalled, err := json.Marshal(task)
		if err != nil {
			t.Fatalf("marshalling failed with error: %v", err)
		}
		request, _ := http.NewRequest("PUT", "/tasks/1", bytes.NewReader(marshalled))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		want := 200
		if got := response.Result().StatusCode; got != want {
			t.Log(response.Body.String())
			t.Errorf("response status is wrong, expected %d but got %d", want, got)
		}
	})
	t.Run("edit task with invalid payload", func(t *testing.T) {
		mock := NewMockStorage(map[string]models.Task{
			"1": {
				Id:          "1",
				Title:       "Task 1",
				Description: "No description",
				Completed:   false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		})
		tHandler := NewTaskHandler(mock)
		router := getTestRouter(t)
		router.PUT("/tasks/:id", tHandler.EditTask)
		task := models.EditTask{}
		marshalled, err := json.Marshal(task)
		if err != nil {
			t.Fatalf("marshalling failed with error: %v", err)
		}
		request, _ := http.NewRequest("PUT", "/tasks/1", bytes.NewReader(marshalled))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		want := 400
		if got := response.Result().StatusCode; got != want {
			t.Errorf("response status is wrong, expected %d but got %d", want, got)
		}
	})
	t.Run("edit task with invalid title", func(t *testing.T) {
		mock := NewMockStorage(map[string]models.Task{
			"1": {
				Id:          "1",
				Title:       "Task 1",
				Description: "No description",
				Completed:   false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		})
		tHandler := NewTaskHandler(mock)
		router := getTestRouter(t)
		router.PUT("/tasks/:id", tHandler.EditTask)
		title := ""
		task := models.EditTask{
			Title: &title,
		}
		marshalled, err := json.Marshal(task)
		if err != nil {
			t.Fatalf("marshalling failed with error: %v", err)
		}
		request, _ := http.NewRequest("PUT", "/tasks/1", bytes.NewReader(marshalled))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		want := 400
		if got := response.Result().StatusCode; got != want {
			t.Errorf("response status is wrong, expected %d but got %d", want, got)
		}
	})
	t.Run("edit task with invalid description", func(t *testing.T) {
		mock := NewMockStorage(map[string]models.Task{
			"1": {
				Id:          "1",
				Title:       "Task 1",
				Description: "No description",
				Completed:   false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		})
		tHandler := NewTaskHandler(mock)
		router := getTestRouter(t)
		router.PUT("/tasks/:id", tHandler.EditTask)
		description := ""
		task := models.EditTask{
			Description: &description,
		}
		marshalled, err := json.Marshal(task)
		if err != nil {
			t.Fatalf("marshalling failed with error: %v", err)
		}
		request, _ := http.NewRequest("PUT", "/tasks/1", bytes.NewReader(marshalled))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		want := 400
		if got := response.Result().StatusCode; got != want {
			t.Errorf("response status is wrong, expected %d but got %d", want, got)
		}
	})
	t.Run("edit task with invalid task id", func(t *testing.T) {
		mock := NewMockStorage(map[string]models.Task{})
		tHandler := NewTaskHandler(mock)
		router := getTestRouter(t)
		router.PUT("/tasks/:id", tHandler.EditTask)
		title := "Task 1 (edited)"
		task := models.EditTask{
			Title: &title,
		}
		marshalled, err := json.Marshal(task)
		if err != nil {
			t.Fatalf("marshalling failed with error: %v", err)
		}
		request, _ := http.NewRequest("PUT", "/tasks/1", bytes.NewReader(marshalled))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		want := 404
		if got := response.Result().StatusCode; got != want {
			t.Errorf("response status is wrong, expected %d but got %d", want, got)
		}
	})
}

func TestMarkTask(t *testing.T) {
	t.Run("mark task as done with valid payload", func(t *testing.T) {
		mock := NewMockStorage(map[string]models.Task{
			"1": {
				Id:          "1",
				Title:       "Task 1",
				Description: "No description",
				Completed:   false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		})
		tHandler := NewTaskHandler(mock)
		router := getTestRouter(t)
		router.PUT("/tasks/:id", tHandler.MarkTask)
		completed := true
		payload := models.MarkTask{
			Completed: &completed,
		}
		marshalled, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("marshalling failed with error: %v", err)
		}
		request, _ := http.NewRequest("PUT", "/tasks/1", bytes.NewReader(marshalled))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		want := 200
		if got := response.Result().StatusCode; got != want {
			t.Log(response.Body.String())
			t.Errorf("response status is wrong, expected %d but got %d", want, got)
		}
	})
	t.Run("mark task as not done with valid payload", func(t *testing.T) {
		mock := NewMockStorage(map[string]models.Task{
			"1": {
				Id:          "1",
				Title:       "Task 1",
				Description: "No description",
				Completed:   true,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		})
		tHandler := NewTaskHandler(mock)
		router := getTestRouter(t)
		router.PUT("/tasks/:id", tHandler.MarkTask)
		completed := false
		payload := models.MarkTask{
			Completed: &completed,
		}
		marshalled, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("marshalling failed with error: %v", err)
		}
		request, _ := http.NewRequest("PUT", "/tasks/1", bytes.NewReader(marshalled))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		want := 200
		if got := response.Result().StatusCode; got != want {
			t.Log(response.Body.String())
			t.Errorf("response status is wrong, expected %d but got %d", want, got)
		}
	})
	t.Run("mark task with invalid payload", func(t *testing.T) {
		mock := NewMockStorage(map[string]models.Task{
			"1": {
				Id:          "1",
				Title:       "Task 1",
				Description: "No description",
				Completed:   true,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		})
		tHandler := NewTaskHandler(mock)
		router := getTestRouter(t)
		router.PUT("/tasks/:id", tHandler.MarkTask)
		payload := models.MarkTask{}
		marshalled, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("marshalling failed with error: %v", err)
		}
		request, _ := http.NewRequest("PUT", "/tasks/1", bytes.NewReader(marshalled))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		want := 400
		if got := response.Result().StatusCode; got != want {
			t.Errorf("response status is wrong, expected %d but got %d", want, got)
		}
	})
	t.Run("mark task with invalid id", func(t *testing.T) {
		mock := NewMockStorage(map[string]models.Task{
			"1": {
				Id:          "1",
				Title:       "Task 1",
				Description: "No description",
				Completed:   true,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		})
		tHandler := NewTaskHandler(mock)
		router := getTestRouter(t)
		router.PUT("/tasks/:id", tHandler.MarkTask)
		completed := false
		payload := models.MarkTask{
			Completed: &completed,
		}
		marshalled, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("marshalling failed with error: %v", err)
		}
		request, _ := http.NewRequest("PUT", "/tasks/2", bytes.NewReader(marshalled))
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		want := 404
		if got := response.Result().StatusCode; got != want {
			t.Errorf("response status is wrong, expected %d but got %d", want, got)
		}
	})
}

func TestSearchTask(t *testing.T) {
	t.Run("search task with title golang", func(t *testing.T) {
		mock := NewMockStorage(map[string]models.Task{
			"1": {
				Id:          "1",
				Title:       "Learning Golang",
				Description: "No description",
				Completed:   true,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		})
		tHandler := NewTaskHandler(mock)
		router := getTestRouter(t)
		router.GET("/search/tasks", tHandler.SearchTask)
		request, _ := http.NewRequest("GET", "/search/tasks?q=golang", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		want := 200
		if got := response.Result().StatusCode; got != want {
			t.Log(response.Body.String())
			t.Errorf("response status is wrong, expected %d but got %d", want, got)
		}
	})
}
