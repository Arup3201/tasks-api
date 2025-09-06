package storage

import (
	"strings"
	"time"

	"github.com/Arup3201/gotasks/internal/models"
)

var tasks = []models.Task{
	{
		Id:          "1",
		Title:       "Learn Golang",
		Description: "Learn web dev with golang gin",
		Completed:   false,
		CreatedAt: time.Date(
			2025, 9, 6, 11, 54, 58, 651387237, time.UTC),
		UpdatedAt: time.Date(
			2025, 9, 6, 11, 54, 58, 651387237, time.UTC),
	},
	{
		Id:          "2",
		Title:       "Learn DSA",
		Description: "Solve 2 sorting problems",
		Completed:   false,
		CreatedAt: time.Date(
			2025, 9, 6, 11, 54, 58, 651387237, time.UTC),
		UpdatedAt: time.Date(
			2025, 9, 6, 11, 54, 58, 651387237, time.UTC),
	},
	{
		Id:          "3",
		Title:       "Learn REST API",
		Description: "Build tasks API with 4 API endpoints",
		Completed:   false,
		CreatedAt: time.Date(
			2025, 9, 6, 11, 54, 58, 651387237, time.UTC),
		UpdatedAt: time.Date(
			2025, 9, 6, 11, 54, 58, 651387237, time.UTC),
	},
}

func GetAllTasks() []models.Task {
	return tasks
}

func GetTaskWithId(id string) (*models.Task, bool) {
	for _, task := range tasks {
		if task.Id == id {
			return &task, true
		}
	}

	return nil, false
}

func AddTask(task models.Task) *models.Task {
	tasks = append(tasks, task)

	return &task
}

func UpdateTask(id string, edit models.UpdateTask) (*models.Task, bool) {
	updated := false
	for i, task := range tasks {
		if task.Id == id {
			if edit.Title != nil && *edit.Title != task.Title {
				task.Title = *edit.Title
				updated = true
			}
			if edit.Description != nil && *edit.Description != task.Description {
				task.Description = *edit.Description
				updated = true
			}
			if edit.Completed != nil && *edit.Completed != task.Completed {
				task.Completed = *edit.Completed
				updated = true
			}

			if updated {
				tasks[i] = task
			}

			return &task, true
		}
	}

	return nil, false
}

func DeleteTask(id string) (*models.Task, bool) {
	for i, task := range tasks {
		if task.Id == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return &task, true
		}
	}

	return nil, false
}

func SearchTasks(query string) []models.Task {
	matches := []models.Task{}
	for _, task := range tasks {
		if strings.Contains(strings.ToLower(task.Title), strings.ToLower(query)) {
			matches = append(matches, task)
		}
	}

	return matches
}
