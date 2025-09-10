package models

import (
	"fmt"
	"time"

	"github.com/Arup3201/gotasks/internal/handlers/apperr"
)

type Task struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateTask struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

type UpdateTask struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Completed   *bool   `json:"completed,omitempty"`
}

type EditTask struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
}

type MarkTask struct {
	Completed *bool `json:"completed,omitempty"`
}

type TaskModel struct {
	store TaskStore
}

func NewTaskModel(store TaskStore) TaskModel {
	return TaskModel{store: store}
}

func (model TaskModel) AllTasks() ([]Task, error) {
	tasks, err := model.store.List()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (model TaskModel) GetTaskByID(id string) (Task, error) {
	task, err := model.store.Get(id)
	if err != nil {
		return task, err
	}

	return task, nil
}

func (model TaskModel) AddTask(data CreateTask) (string, error) {
	if data.Title == nil {
		return "", apperr.MissingBodyError(apperr.ErrorField{
			Reason: "The body property 'title' is required",
			Field:  "title",
		})
	}

	if data.Description == nil {
		return "", apperr.MissingBodyError(apperr.ErrorField{
			Reason: "The body property 'description' is required",
			Field:  "description",
		})
	}

	if *data.Title == "" {
		return "", apperr.InvalidBodyError(apperr.ErrorField{
			Reason: "Body property 'title' can't be empty",
			Field:  "title",
		})
	}

	if *data.Description == "" {
		return "", apperr.InvalidBodyError(apperr.ErrorField{
			Reason: "Body property 'description' can't be empty",
			Field:  "description",
		})
	}

	newTask, err := model.store.Insert(*data.Title, *data.Description)
	if err != nil {
		return "", fmt.Errorf("store.Insert error: %v", err)
	}

	return newTask, nil
}

func (model TaskModel) EditTask(id string, data EditTask) (Task, error) {
	var task Task

	if data.Title == nil && data.Description == nil {
		return task, apperr.InvalidBodyError(apperr.ErrorField{
			Reason: "Atleast one body property 'title' or 'description' is required",
			Field:  "title, description",
		})
	}

	if data.Title != nil && *data.Title == "" {
		return task, apperr.InvalidBodyError(apperr.ErrorField{
			Reason: "Body property 'title' can't be empty",
			Field:  "title",
		})
	}

	if data.Description != nil && *data.Description == "" {
		return task, apperr.InvalidBodyError(apperr.ErrorField{
			Reason: "Body property 'description' can't be empty",
			Field:  "description",
		})
	}

	task, err := model.store.Update(id, data.Title, data.Description, nil)
	if err != nil {
		return task, err
	}

	return task, nil
}

func (model TaskModel) MarkTask(id string, data MarkTask) (Task, error) {
	var task Task

	if data.Completed == nil {
		return task, apperr.MissingBodyError(apperr.ErrorField{
			Reason: "Body property 'completed' is required",
			Field:  "completed",
		})
	}

	task, err := model.store.Update(id, nil, nil, data.Completed)
	if err != nil {
		return task, err
	}

	return task, nil
}

func (model TaskModel) DeleteTask(id string) (string, error) {
	taskId, err := model.store.Delete(id)
	if err != nil {
		return "", err
	}

	return taskId, nil
}

func (model TaskModel) SearchTask(query string) ([]Task, error) {
	if query == "" {
		return nil, apperr.InvalidRequestParamError(apperr.ErrorField{
			Reason: "Search query can't be empty",
			Field:  "q",
		})
	}

	tasks, err := model.store.Search("title", query)
	if err != nil {
		return nil, fmt.Errorf("TaskModel.SearchTask failed with error: %v", err)
	}

	return tasks, nil
}
