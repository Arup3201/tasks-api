package task

import (
	"fmt"
	"strings"

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

func (ts TaskService) CreateTask(title, description string) (*task.Task, error) {
	if strings.TrimSpace(title) == "" {
		return nil, errors.InputValidationError("Invalid task 'title'", "Task property 'title' can't be empty")
	}

	if strings.TrimSpace(description) == "" {
		return nil, errors.InputValidationError("Invalid task 'description'", "Task property 'description' can't be empty")
	}

	task, err := ts.taskRepository.Insert(title, description)

	if err != nil {
		return nil, err
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

func (ts TaskService) UpdateTask(taskId string, data UpdateTaskData) (*task.Task, error) {
	update := map[string]any{}

	if data.title != nil {
		update["Title"] = *data.title
	}

	if data.description != nil {
		update["Description"] = *data.description
	}

	if data.isCompleted != nil {
		update["IsCompleted"] = *data.isCompleted
	}

	task := ts.taskRepository.Update(taskId, update)
	if task == nil {
		return nil, errors.NotFoundError(fmt.Sprintf("Task with ID '%s' not found", taskId))
	}

	return task, nil
}

func (ts TaskService) SearchTasks(query string) []task.Task {
	allTasks := ts.taskRepository.List()

	matches := []task.Task{}
	contains := false
	for _, task := range allTasks {
		contains = true
		for part := range strings.SplitSeq(query, " ") {
			if !strings.Contains(strings.ToLower(task.Title), strings.ToLower(part)) {
				contains = false
			}
		}
		if contains {
			matches = append(matches, task)
		}
	}

	return matches
}
