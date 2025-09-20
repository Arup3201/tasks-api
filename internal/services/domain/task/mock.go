package task

import (
	"fmt"
	"reflect"
	"time"

	"github.com/Arup3201/gotasks/internal/entities/task"
	"github.com/Arup3201/gotasks/internal/errors"
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

func (tr *mockTaskRepository) Get(taskId int) (*task.Task, error) {
	for _, task := range tr.tasks {
		if task.Id == taskId {
			return &task, nil
		}
	}

	return nil, errors.NotFoundError(fmt.Sprintf("Task with ID %d not found", taskId))
}

func (tr *mockTaskRepository) Insert(id int, title, description string) (*task.Task, error) {
	task := task.Task{
		Id:          id,
		Title:       title,
		Description: description,
		IsCompleted: false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	tr.tasks = append(tr.tasks, task)
	return &task, nil
}

func (tr *mockTaskRepository) Update(taskId int, data map[string]any) (*task.Task, error) {
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
			return &task, nil
		}
	}

	return nil, errors.NotFoundError(fmt.Sprintf("Task with ID %d not found", taskId))
}

func (tr *mockTaskRepository) List() ([]task.Task, error) {
	return tr.tasks, nil
}
