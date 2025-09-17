package task

import (
	"strings"
	"time"

	"github.com/Arup3201/gotasks/internal/entities/task"
	"github.com/Arup3201/gotasks/internal/services/errors"
	"github.com/google/uuid"
)

func CreateTask(title, description string) (task.Task, error) {
	if strings.TrimSpace(title) == "" {
		return task.Task{}, errors.NewInputValidationError("Invalid property value", "task property 'title' can't be empty")
	}

	if strings.TrimSpace(description) == "" {
		return task.Task{}, errors.NewInputValidationError("Invalid property value", "task property 'description' can't be empty")
	}

	return task.Task{
		Id:          uuid.New().String(),
		Title:       title,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}
