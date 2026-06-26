package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Arup3201/gotasks/internal/models"
)

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
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsCompleted bool      `json:"is_completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		IsCompleted: task.IsCompleted,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	})
}
