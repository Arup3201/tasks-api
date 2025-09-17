package task

import (
	"fmt"
	"strings"
	"time"

	"github.com/Arup3201/gotasks/internal/entities/task"
	"github.com/Arup3201/gotasks/internal/services/errors"
	"github.com/Arup3201/gotasks/internal/storages"
)

type UpdateTaskData struct {
	title       *string
	description *string
	isCompleted *bool
}

type TaskService struct {
	taskRepository storages.TaskRepository
}

func NewTaskService(repo storages.TaskRepository) *TaskService {
	return &TaskService{
		taskRepository: repo,
	}
}

func (ts TaskService) CreateTask(title, description string) (task.Task, error) {
	if strings.TrimSpace(title) == "" {
		return task.Task{}, errors.InputValidationError("Invalid task 'title'", "Task property 'title' can't be empty")
	}

	if strings.TrimSpace(description) == "" {
		return task.Task{}, errors.InputValidationError("Invalid task 'description'", "Task property 'description' can't be empty")
	}

	task, err := ts.taskRepository.Insert(title, description)

	if err != nil {
		return task, err
	}

	return task, nil
}

func (ts TaskService) GetTask(taskId string) (*task.Task, error) {
	task := ts.taskRepository.Get(taskId)
	if task == nil {
		return nil, errors.NotFoundError(fmt.Sprintf("Task with ID '%s' not found", taskId))
	}

	return task, nil
}

func (ts TaskService) UpdateTask(taskId string, data UpdateTaskData) (task.Task, error) {
	var task task.Task
	if data.title != nil {
		task, err := ts.taskRepository.Update(*data.title)
	}

	if data.description != nil {
		task, err := ts.taskRepository.Update(*data.title)

	}

	if data.isCompleted != nil {
		task.IsCompleted = *data.isCompleted
		task.UpdatedAt = time.Now()
	}

	return task, nil
}
