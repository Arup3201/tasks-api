package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Arup3201/gotasks/internal/models"
)

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsCompleted bool      `json:"is_completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TaskController struct {
	taskService *models.TaskService
}

func NewTaskController(taskService *models.TaskService) *TaskController {
	return &TaskController{taskService}
}

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateTaskResponse struct {
	Task Task `json:"task"`
}

func (tc *TaskController) CreateTask(w http.ResponseWriter, r *http.Request) {
	var data CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {

		http.Error(w,
			"json parse error",
			http.StatusBadRequest)
		return
	}

	userID, err := getUserID(r)
	if err != nil {
		http.Error(w,
			"not authenticated",
			http.StatusUnauthorized)
		return
	}

	task, err := tc.taskService.CreateTask(r.Context(), userID, data.Title, data.Description)
	if errors.Is(err, models.ErrInvalidTask) {
		http.Error(w,
			"invalid task title or description provided",
			http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w,
			"server error",
			http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreateTaskResponse{
		Task: Task{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			IsCompleted: task.IsCompleted,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		},
	})
}

type TaskUpdateRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	IsCompleted *bool   `json:"is_completed"`
}

type UpdateTaskResponse struct {
	Task Task `json:"task"`
}

func (tc *TaskController) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w,
			"empty task ID",
			http.StatusBadRequest)
		return
	}

	var data TaskUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {

		http.Error(w,
			"json parse error",
			http.StatusBadRequest)
		return
	}

	userID, err := getUserID(r)
	if err != nil {
		http.Error(w,
			"not authenticated",
			http.StatusUnauthorized)
		return
	}

	task, err := tc.taskService.UpdateTask(r.Context(), id, userID, data.Title, data.Description, data.IsCompleted)
	switch {
	case errors.Is(err, models.ErrTaskNotFound):
		http.Error(w,
			"task not found",
			http.StatusNotFound)
	case errors.Is(err, models.ErrInvalidTask):
		http.Error(w,
			"invalid task title or description provided",
			http.StatusBadRequest)
	case err != nil:
		http.Error(w,
			"server error",
			http.StatusInternalServerError)
	}
	if err != nil {
		return
	}

	json.NewEncoder(w).Encode(UpdateTaskResponse{
		Task: Task{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			IsCompleted: task.IsCompleted,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		},
	})
}

type ListTasksResponse struct {
	Tasks []Task `json:"tasks"`
}

func (tc *TaskController) ListTasks(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		http.Error(w,
			"not authenticated",
			http.StatusUnauthorized)
		return
	}

	tasks, err := tc.taskService.ListTasks(r.Context(), userID)
	if err != nil {
		http.Error(w,
			"server error",
			http.StatusInternalServerError)
		return
	}

	listedTasks := []Task{}
	for _, t := range tasks {
		listedTasks = append(listedTasks, Task{
			ID:          t.ID,
			Title:       t.Title,
			Description: t.Description,
			IsCompleted: t.IsCompleted,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		})
	}

	json.NewEncoder(w).Encode(ListTasksResponse{
		Tasks: listedTasks,
	})
}
