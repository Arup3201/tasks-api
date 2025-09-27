package services

import "github.com/Arup3201/gotasks/internal/entities/task"

type UpdateTaskData struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	IsCompleted *bool   `json:"is_completed"`
}

type ServiceHandler interface {
	GetAllTasks() ([]task.Task, error)
	CreateTask(title, description string) (*task.Task, error)
	GetTask(taskId int) (*task.Task, error)
	UpdateTask(taskId int, data UpdateTaskData) (*task.Task, error)
	DeleteTask(taskId int) (*int, error)
	SearchTasks(query string) ([]task.Task, error)
	UpdateLastInsertedId(lastInsertedId int) // setter funtion
}
