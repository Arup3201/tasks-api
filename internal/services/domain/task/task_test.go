package task

import (
	"testing"

	"github.com/Arup3201/gotasks/internal/services/errors"
)

func TestAddTask(t *testing.T) {
	t.Run("Create a task - 1", func(t *testing.T) {
		title := "Test task"
		description := "Test task description"

		got, _ := CreateTask(title, description)

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

		got, _ := CreateTask(title, description)

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

		got, _ := CreateTask(title, description)

		if got.Id == "" {
			t.Errorf("expected non empty task ID")
		}
	})
	t.Run("Two tasks with different ID", func(t *testing.T) {
		title := "Test task"
		description := "Test task description"

		task1, _ := CreateTask(title, description)
		task2, _ := CreateTask(title, description)

		if task1.Id == task2.Id {
			t.Errorf("Two tasks can't have same ID")
		}
	})
	t.Run("Create task with valid created_at value", func(t *testing.T) {
		title := "Test task"
		description := "Test task description"

		task, _ := CreateTask(title, description)

		if task.CreatedAt.IsZero() {
			t.Errorf("created task has zero created_at value")
		}
	})
	t.Run("Create task with valid updated_at value", func(t *testing.T) {
		title := "Test task"
		description := "Test task description"

		task, _ := CreateTask(title, description)

		if task.UpdatedAt.IsZero() {
			t.Errorf("created task has zero updated_at value")
		}
	})
	t.Run("Fail to create task with empty title", func(t *testing.T) {
		title := ""
		description := "Test task description"

		_, err := CreateTask(title, description)
		if _, ok := err.(errors.InputValidationError); !ok {
			t.Errorf("expected InputValidationError in task creation for providing empty title")
		}
	})
	t.Run("Fail to create task with empty description", func(t *testing.T) {
		title := "Test task"
		description := ""

		_, err := CreateTask(title, description)
		if _, ok := err.(errors.InputValidationError); !ok {
			t.Errorf("expected InputValidationError in task creation for providing empty description")
		}
	})
}
