package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Arup3201/gotasks/internal/handlers/apperr"
	"github.com/Arup3201/gotasks/internal/models"
	"github.com/gin-gonic/gin"
)

type MockStorage struct {
	tasks map[string]models.Task
}

func (m *MockStorage) Get(id string) (models.Task, error) {
	task, ok := m.tasks[id]
	if !ok {
		return task, apperr.NotFoundError()
	}

	return task, nil
}
func (m *MockStorage) Insert(title, description string) (string, error) {
	return "ERR", nil
}
func (m *MockStorage) Update(id string, title, description *string, completed *bool) (models.Task, error) {
	return models.Task{}, nil
}
func (m *MockStorage) Delete(id string) (string, error) {
	return "ERR", nil
}
func (m *MockStorage) List() ([]models.Task, error) {
	tasks := []models.Task{}
	for _, task := range m.tasks {
		tasks = append(tasks, task)
	}

	return tasks, nil
}
func (m *MockStorage) Search(by models.FieldName, query string) ([]models.Task, error) {
	return nil, nil
}

func getTestContext(t testing.TB, w http.ResponseWriter, r *http.Request) *gin.Context {
	t.Helper()

	gin.SetMode(gin.TestMode)

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = r

	return ctx
}

func TestGETAllTasks(t *testing.T) {
	t.Run("returns 2 tasks", func(t *testing.T) {
		mock := &MockStorage{
			tasks: map[string]models.Task{
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
			},
		}
		tHandler := NewTaskHandler(mock)
		request, _ := http.NewRequest("GET", "/tasks", nil)
		response := httptest.NewRecorder()
		ctx := getTestContext(t, response, request)

		tHandler.GetAllTasks(ctx)

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
		mock := &MockStorage{}
		tHandler := NewTaskHandler(mock)
		request, _ := http.NewRequest("GET", "/tasks", nil)
		response := httptest.NewRecorder()
		ctx := getTestContext(t, response, request)

		tHandler.GetAllTasks(ctx)

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
