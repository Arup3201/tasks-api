package task

import (
	"reflect"
	"time"

	"github.com/Arup3201/gotasks/internal/entities/task"
	"github.com/google/uuid"
)

/* Mock up of the task repository for test */

type mockTaskRepository struct {
	tasks []task.Task
}

func NewMockTaskRepository() *mockTaskRepository {
	return &mockTaskRepository{
		tasks: make([]task.Task, 0),
	}
}

func (tr *mockTaskRepository) Get(taskId string) *task.Task {
	for _, task := range tr.tasks {
		if task.Id == taskId {
			return &task
		}
	}

	return nil
}

func (tr *mockTaskRepository) Insert(title, description string) (*task.Task, error) {
	task := task.Task{
		Id:          uuid.New().String(),
		Title:       title,
		Description: description,
		IsCompleted: false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	tr.tasks = append(tr.tasks, task)
	return &task, nil
}

func (tr *mockTaskRepository) Update(taskId string, data map[string]any) *task.Task {
	for i, task := range tr.tasks {
		if task.Id == taskId {
			t := reflect.ValueOf(&task).Elem()
			for f, v := range data {
				field := t.FieldByName(f)
				value := reflect.ValueOf(v)
				if field.Kind() == reflect.String {
					field.SetString(value.String())
				} else if field.Kind() == reflect.Bool {
					field.SetBool(value.Bool())
				}
			}
			task.UpdatedAt = time.Now()
			tr.tasks[i] = task
			return &task
		}
	}

	return nil
}
