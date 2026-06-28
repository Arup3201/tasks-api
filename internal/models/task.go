package models

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidTask = errors.New("invalid task value")
)

type TaskModel struct {
	ID, UserID, Title, Description string
	IsCompleted                    bool
	CreatedAt, UpdatedAt           time.Time
}

type TaskStoreInterface interface {
	Create(ctx context.Context,
		id, userID, title, description string,
		isCompleted bool) error
	Get(ctx context.Context,
		id, userID string) (*Task, error)
	Update(ctx context.Context,
		task *Task) error
	List(ctx context.Context,
		userID string) ([]Task, error)
	Delete(ctx context.Context,
		id, userID string) error
}

type TaskService struct {
	store TaskStoreInterface
}

func NewTaskService(store TaskStoreInterface) *TaskService {
	return &TaskService{store: store}
}

func (ts *TaskService) CreateTask(ctx context.Context,
	userID, title, description string) (*TaskModel, error) {
	if strings.TrimSpace(title) == "" {
		return nil, ErrInvalidTask
	}

	if strings.TrimSpace(description) == "" {
		return nil, ErrInvalidTask
	}

	id := uuid.NewString()
	err := ts.store.Create(ctx, id, userID, title, description, false)
	if err != nil {
		return nil, err
	}

	task, err := ts.store.Get(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	return &TaskModel{
		ID:          task.ID,
		UserID:      task.UserID,
		Title:       task.Title,
		Description: task.Description,
		IsCompleted: task.IsCompleted,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}, nil
}

func (ts *TaskService) UpdateTask(ctx context.Context,
	id, userID string,
	title, description *string,
	isCompleted *bool) (*TaskModel, error) {

	task, err := ts.store.Get(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	if title != nil {
		if strings.TrimSpace(*title) == "" {
			return nil, ErrInvalidTask
		}

		task.Title = *title
	}

	if description != nil {
		if strings.TrimSpace(*description) == "" {
			return nil, ErrInvalidTask
		}

		task.Description = *description
	}

	if isCompleted != nil {
		task.IsCompleted = *isCompleted
	}

	err = ts.store.Update(ctx, task)
	if err != nil {
		return nil, err
	}

	task, err = ts.store.Get(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	return &TaskModel{
		ID:          task.ID,
		UserID:      task.UserID,
		Title:       task.Title,
		Description: task.Description,
		IsCompleted: task.IsCompleted,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}, nil
}

func (ts *TaskService) ListTasks(ctx context.Context,
	userID string) ([]TaskModel, error) {

	rows, err := ts.store.List(ctx, userID)
	if err != nil {
		return nil, err
	}

	tasks := []TaskModel{}
	for _, r := range rows {
		tasks = append(tasks, TaskModel{
			ID:          r.ID,
			UserID:      r.UserID,
			Title:       r.Title,
			Description: r.Description,
			IsCompleted: r.IsCompleted,
			CreatedAt:   r.CreatedAt,
			UpdatedAt:   r.UpdatedAt,
		})
	}

	return tasks, nil
}

func (ts *TaskService) DeleteTask(ctx context.Context,
	id, userID string) error {

	err := ts.store.Delete(ctx, id, userID)
	if err != nil {
		return err
	}

	return nil
}
