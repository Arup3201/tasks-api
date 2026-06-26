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
		id, userID string) (*TaskModel, error)
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

	return task, nil
}
