package task

import (
	"testing"

	"github.com/Arup3201/gotasks/internal/services/errors"
)

func TestAddTask(t *testing.T) {
	t.Run("Create a task - 1", func(t *testing.T) {
		title := "Test task"
		description := "Test task description"
		ts := NewTaskService(NewMockTaskRepository())

		got, _ := ts.CreateTask(title, description)

		if got.Title != title {
			t.Errorf("expected title %s but got %s", title, got.Title)
		}
		if got.Description != description {
			t.Errorf("expected description %s but got %s", description, got.Description)
		}
	})
	t.Run("Create a task - 2", func(t *testing.T) {
		title := "Test task 1"
		description := "Test task 1 description"
		ts := NewTaskService(NewMockTaskRepository())

		got, _ := ts.CreateTask(title, description)

		if got.Title != title {
			t.Errorf("expected title %s but got %s", title, got.Title)
		}
		if got.Description != description {
			t.Errorf("expected description %s but got %s", description, got.Description)
		}
	})
	t.Run("Valid task ID", func(t *testing.T) {
		title := "Test task"
		description := "Test task description"
		ts := NewTaskService(NewMockTaskRepository())

		got, _ := ts.CreateTask(title, description)

		if got.Id == "" {
			t.Errorf("expected non empty task ID")
		}
	})
	t.Run("Two tasks with different ID", func(t *testing.T) {
		title := "Test task"
		description := "Test task description"
		ts := NewTaskService(NewMockTaskRepository())

		task1, _ := ts.CreateTask(title, description)
		task2, _ := ts.CreateTask(title, description)

		if task1.Id == task2.Id {
			t.Errorf("Two tasks can't have same ID")
		}
	})
	t.Run("Create task with valid created_at value", func(t *testing.T) {
		title := "Test task"
		description := "Test task description"
		ts := NewTaskService(NewMockTaskRepository())

		task, _ := ts.CreateTask(title, description)

		if task.CreatedAt.IsZero() {
			t.Errorf("created task has zero created_at value")
		}
	})
	t.Run("Create task with valid updated_at value", func(t *testing.T) {
		title := "Test task"
		description := "Test task description"
		ts := NewTaskService(NewMockTaskRepository())

		task, _ := ts.CreateTask(title, description)

		if task.UpdatedAt.IsZero() {
			t.Errorf("created task has zero updated_at value")
		}
	})
	t.Run("Fail to create task with empty title", func(t *testing.T) {
		title := ""
		description := "Test task description"
		ts := NewTaskService(NewMockTaskRepository())

		_, err := ts.CreateTask(title, description)
		inputInvalidError, ok := err.(errors.Error)
		if !ok {
			t.Errorf("expected `Error` on create task with empty title")
		}
		if inputInvalidError.Type != "INVALID_INPUT" {
			t.Errorf("expected `INVALID_INPUT` type error on create task with empty title")
		}
	})
	t.Run("Fail to create task with empty description", func(t *testing.T) {
		title := "Test task"
		description := ""
		ts := NewTaskService(NewMockTaskRepository())

		_, err := ts.CreateTask(title, description)
		inputInvalidError, ok := err.(errors.Error)
		if !ok {
			t.Errorf("expected `Error` on create task with empty description")
		}
		if inputInvalidError.Type != "INVALID_INPUT" {
			t.Errorf("expected `INVALID_INPUT` type error on create task with empty description")
		}
	})
}

func TestGetTask(t *testing.T) {
	t.Run("Get task ID is correct after creating", func(t *testing.T) {
		title := "Test task"
		description := "Test task description"
		ts := NewTaskService(NewMockTaskRepository())
		created, _ := ts.CreateTask(title, description)

		task, _ := ts.GetTask(created.Id)

		if task.Id != created.Id {
			t.Errorf("Task ID does not match, expected %s but got %s", created.Id, task.Id)
		}
	})
	t.Run("Get task title is correct after creating - 1", func(t *testing.T) {
		title := "Test task 1"
		description := "Test task description"
		ts := NewTaskService(NewMockTaskRepository())
		created, _ := ts.CreateTask(title, description)

		task, _ := ts.GetTask(created.Id)

		if task.Title != title {
			t.Errorf("Task title does not match, expected %s but got %s", title, task.Title)
		}
	})
	t.Run("Get task title is correct after creating - 2", func(t *testing.T) {
		title := "Test task 2"
		description := "Test task description"
		ts := NewTaskService(NewMockTaskRepository())
		created, _ := ts.CreateTask(title, description)

		task, _ := ts.GetTask(created.Id)

		if task.Title != title {
			t.Errorf("Task title does not match, expected %s but got %s", title, task.Title)
		}
	})
}
