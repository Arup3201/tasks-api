package task

import (
	"fmt"
	"strings"

	"github.com/Arup3201/gotasks/internal/entities/task"
	"github.com/Arup3201/gotasks/internal/errors"
	"github.com/Arup3201/gotasks/internal/services"
	"github.com/Arup3201/gotasks/internal/storages"
)

type TaskService struct {
	taskRepository storages.TaskRepository
	lastTaskId     int
}

func NewTaskService(repo storages.TaskRepository) *TaskService {
	return &TaskService{
		taskRepository: repo,
		lastTaskId:     0,
	}
}

func (ts *TaskService) CreateTask(title, description string) (*task.Task, error) {
	if strings.TrimSpace(title) == "" {
		return nil, errors.InputValidationError("Invalid task 'title'", "Task property 'title' can't be empty")
	}

	if strings.TrimSpace(description) == "" {
		return nil, errors.InputValidationError("Invalid task 'description'", "Task property 'description' can't be empty")
	}

	task, err := ts.taskRepository.Insert(ts.lastTaskId+1, title, description)
	if err != nil {
		return nil, err
	}
	ts.lastTaskId += 1

	return task, nil
}

func (ts *TaskService) GetTask(taskId int) (*task.Task, error) {
	task, err := ts.taskRepository.Get(taskId)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, err
	}

	return task, nil
}

func (ts *TaskService) UpdateTask(taskId int, data services.UpdateTaskData) (*task.Task, error) {
	update := map[string]any{}

	if data.Title != nil {
		update["Title"] = *data.Title
	}

	if data.Description != nil {
		update["Description"] = *data.Description
	}

	if data.IsCompleted != nil {
		update["IsCompleted"] = *data.IsCompleted
	}

	task, err := ts.taskRepository.Update(taskId, update)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, errors.NotFoundError(fmt.Sprintf("Task with ID '%d' not found", taskId))
	}

	return task, nil
}

func (ts *TaskService) SearchTasks(query string) ([]task.Task, error) {
	allTasks, err := ts.taskRepository.List()
	if err != nil {
		return nil, err
	}

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

	return matches, nil
}
