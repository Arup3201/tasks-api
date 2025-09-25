package httpController

import (
	"fmt"
	"reflect"
	"time"

	entities "github.com/Arup3201/gotasks/internal/entities/task"
	serverErrors "github.com/Arup3201/gotasks/internal/errors"
)

type MockRepository struct {
	tasks []entities.Task
}

func (tr *MockRepository) Get(taskId int) (*entities.Task, error) {
	for _, task := range tr.tasks {
		if task.Id == taskId {
			return &task, nil
		}
	}

	return nil, serverErrors.NotFoundError(fmt.Sprintf("Task with ID %d not found", taskId))
}

func (tr *MockRepository) Insert(id int, title, description string) (*entities.Task, error) {
	task := entities.Task{
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

func (tr *MockRepository) Update(taskId int, data map[string]any) (*entities.Task, error) {
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

	return nil, serverErrors.NotFoundError(fmt.Sprintf("Task with ID %d not found", taskId))
}

func (tr *MockRepository) Delete(taskId int) (*int, error) {
	for i, task := range tr.tasks {
		if task.Id == taskId {
			tr.tasks = append(tr.tasks[:i], tr.tasks[i+1:]...)
			return &task.Id, nil
		}
	}

	return nil, serverErrors.NotFoundError(fmt.Sprintf("Task with ID %d not found", taskId))
}

func (tr *MockRepository) List() ([]entities.Task, error) {
	return tr.tasks, nil
}

func (tr *MockRepository) Close() error {
	return nil
}
