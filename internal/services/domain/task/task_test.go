package task

import (
	"testing"
	"time"

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

		if got.Id == 0 {
			t.Errorf("expected non-zero task ID")
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
			t.Errorf("Task ID does not match, expected %d but got %d", created.Id, task.Id)
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

func TestUpdateTask(t *testing.T) {
	t.Run("update task updates correct task", func(t *testing.T) {
		title := "Test task"
		description := "Test task description"
		ts := NewTaskService(NewMockTaskRepository())
		created, _ := ts.CreateTask(title, description)
		updated_title := "Test task (updated)"

		updated, _ := ts.UpdateTask(created.Id, UpdateTaskData{
			title: &updated_title,
		})

		if updated == nil {
			t.Errorf("updated response is nil")
		}
		if updated != nil && updated.Id != created.Id {
			t.Errorf("wrong updated task expected %d but got %d", created.Id, updated.Id)
		}
	})
	t.Run("update task updates title", func(t *testing.T) {
		title := "Test task"
		description := "Test task description"
		ts := NewTaskService(NewMockTaskRepository())
		created, _ := ts.CreateTask(title, description)
		updated_title := "Test task (updated)"

		updated, _ := ts.UpdateTask(created.Id, UpdateTaskData{
			title: &updated_title,
		})

		if updated.Title != updated_title {
			t.Errorf("updated task title expected %s but got %s", updated_title, updated.Title)
		}
	})
	t.Run("update task updates description", func(t *testing.T) {
		title := "Test task"
		description := "Test task description"
		ts := NewTaskService(NewMockTaskRepository())
		created, _ := ts.CreateTask(title, description)
		updated_description := "Test task description (updated)"

		updated, _ := ts.UpdateTask(created.Id, UpdateTaskData{
			description: &updated_description,
		})

		if updated.Description != updated_description {
			t.Errorf("updated task description expected %s but got %s", updated_description, updated.Description)
		}
	})
	t.Run("update task updates is_completed", func(t *testing.T) {
		title := "Test task"
		description := "Test task description"
		ts := NewTaskService(NewMockTaskRepository())
		created, _ := ts.CreateTask(title, description)
		isCompleted := true

		updated, _ := ts.UpdateTask(created.Id, UpdateTaskData{
			isCompleted: &isCompleted,
		})

		if updated.IsCompleted != isCompleted {
			t.Errorf("updated task is_completed expected %t but got %t", isCompleted, updated.IsCompleted)
		}
	})
	t.Run("update task updates updated_at value", func(t *testing.T) {
		title := "Test task"
		description := "Test task description"
		ts := NewTaskService(NewMockTaskRepository())
		created, _ := ts.CreateTask(title, description)
		time.Sleep(1000 * 2) // 2 secs
		updated_title := "Test task (updated)"

		updated, _ := ts.UpdateTask(created.Id, UpdateTaskData{
			title: &updated_title,
		})

		if updated.UpdatedAt.IsZero() {
			t.Errorf("update task updated_at is zero")
		}
		if updated.UpdatedAt.Equal(created.UpdatedAt) {
			t.Errorf("update task updated_at did not change, got %v same as when created %v", created.UpdatedAt, updated.UpdatedAt)
		}
	})
	t.Run("update task persists the update", func(t *testing.T) {
		title := "Test task"
		description := "Test task description"
		ts := NewTaskService(NewMockTaskRepository())
		created, _ := ts.CreateTask(title, description)
		updated_title := "Test task (updated)"
		updated, _ := ts.UpdateTask(created.Id, UpdateTaskData{
			title: &updated_title,
		})

		task, _ := ts.GetTask(created.Id)

		if task.Title != updated.Title {
			t.Errorf("task update did not persist, expected %s but got %s", updated.Title, task.Title)
		}
	})
}

func TestSearchTasks(t *testing.T) {
	t.Run("Search tasks with match", func(t *testing.T) {
		tasks := []struct {
			title       string
			description string
		}{
			{
				title:       "Learn Golang",
				description: "Learn reflect concept in Golang",
			},
			{
				title:       "Learn Python",
				description: "Learn dict concept in Python",
			},
			{
				title:       "Play Football",
				description: "2 Hrs football time at evening",
			},
			{
				title:       "Play Basketball",
				description: "1 Hr basketball time at evening",
			},
		}
		ts := NewTaskService(NewMockTaskRepository())
		for _, task := range tasks {
			ts.CreateTask(task.title, task.description)
		}
		query := "learn"

		results := ts.SearchTasks(query)

		want := 2
		if got := len(results); got != want {
			t.Errorf("expected searched results %d but got %d", want, got)
		}
	})
	t.Run("Search tasks with no match", func(t *testing.T) {
		tasks := []struct {
			title       string
			description string
		}{
			{
				title:       "Learn Golang",
				description: "Learn reflect concept in Golang",
			},
			{
				title:       "Learn Python",
				description: "Learn dict concept in Python",
			},
			{
				title:       "Play Football",
				description: "2 Hrs football time at evening",
			},
			{
				title:       "Play Basketball",
				description: "1 Hr basketball time at evening",
			},
		}
		ts := NewTaskService(NewMockTaskRepository())
		for _, task := range tasks {
			ts.CreateTask(task.title, task.description)
		}
		query := "nothing"

		results := ts.SearchTasks(query)

		want := 0
		if got := len(results); got != want {
			t.Errorf("expected searched results %d but got %d", want, got)
		}
	})
	t.Run("Search tasks with multi word query", func(t *testing.T) {
		tasks := []struct {
			title       string
			description string
		}{
			{
				title:       "Learn DSA in Golang ",
				description: "Learn reflect concept in Golang",
			},
			{
				title:       "Learn REST API design in Python",
				description: "Learn dict concept in Python",
			},
			{
				title:       "Play Football for 2 hrs",
				description: "2 Hrs football time at evening",
			},
			{
				title:       "Play Basketball 1hr",
				description: "1 Hr basketball time at evening",
			},
		}
		ts := NewTaskService(NewMockTaskRepository())
		for _, task := range tasks {
			ts.CreateTask(task.title, task.description)
		}
		query := "learn golang"

		results := ts.SearchTasks(query)

		want := 1
		if got := len(results); got != want {
			t.Errorf("expected searched results %d but got %d", want, got)
		}
	})
	t.Run("Search tasks with multi word query", func(t *testing.T) {
		tasks := []struct {
			title       string
			description string
		}{
			{
				title:       "Learn DSA in Go language",
				description: "Learn reflect concept in Golang",
			},
			{
				title:       "Learn REST API design in Python language",
				description: "Learn dict concept in Python",
			},
			{
				title:       "Play Football for 1 hr",
				description: "1 Hr football at evening",
			},
			{
				title:       "Play Basketball 1 hr",
				description: "1 Hr basketball at evening",
			},
		}
		ts := NewTaskService(NewMockTaskRepository())
		for _, task := range tasks {
			ts.CreateTask(task.title, task.description)
		}
		query := "learn language"

		results := ts.SearchTasks(query)

		want := 2
		if got := len(results); got != want {
			t.Errorf("expected searched results %d but got %d", want, got)
		}
	})
	t.Run("Search tasks with multi word query", func(t *testing.T) {
		tasks := []struct {
			title       string
			description string
		}{
			{
				title:       "Learn DSA in Go language",
				description: "Learn reflect concept in Golang",
			},
			{
				title:       "Learn REST API design in Python language",
				description: "Learn dict concept in Python",
			},
			{
				title:       "Play Football for 2 hrs",
				description: "1 Hr football at evening",
			},
			{
				title:       "Play Basketball 1 hr",
				description: "1 Hr basketball at evening",
			},
		}
		ts := NewTaskService(NewMockTaskRepository())
		for _, task := range tasks {
			ts.CreateTask(task.title, task.description)
		}
		query := "play hr"

		results := ts.SearchTasks(query)

		want := 2 // tasks[2].title has 'hrs' which has 'hr' in it
		if got := len(results); got != want {
			t.Errorf("expected searched results %d but got %d", want, got)
		}
	})
}
