package models

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type Task struct {
	ID                   string `gorm:"primaryKey"`
	UserID               string
	Title, Description   string
	IsCompleted          bool
	CreatedAt, UpdatedAt time.Time
}

type TaskStore struct {
	db *gorm.DB
}

func NewTaskStore(db *gorm.DB) *TaskStore {
	return &TaskStore{db: db}
}

func (ts *TaskStore) Create(ctx context.Context,
	id, userID, title, description string,
	isCompleted bool) error {
	task := Task{
		ID:          id,
		UserID:      userID,
		Title:       title,
		Description: description,
		IsCompleted: isCompleted,
	}

	err := gorm.G[Task](ts.db).Create(ctx, &task)
	if err != nil {
		return err
	}

	return nil
}

func (ts *TaskStore) Get(ctx context.Context,
	id, userID string) (*Task, error) {

	task, err := gorm.
		G[Task](ts.db).
		Where("id = ? AND user_id = ?", id, userID).
		First(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrTaskNotFound
	} else if err != nil {
		return nil, err
	}

	return &task, nil
}

func (ts *TaskStore) Update(ctx context.Context,
	task *Task) error {
	if err := ts.db.Save(task).Error; err != nil {
		return err
	}

	return nil
}
