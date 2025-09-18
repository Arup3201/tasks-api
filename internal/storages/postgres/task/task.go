package task

import (
	"database/sql"
	"time"

	"github.com/Arup3201/gotasks/internal/entities/task"
)

type PgTaskRepository struct {
	db *sql.DB
}

func NewPgTaskRepository(db *sql.DB) *PgTaskRepository {
	return &PgTaskRepository{
		db: db,
	}
}

func (pg *PgTaskRepository) Insert(taskId, taskTitle, taskDesc string) (*task.Task, error) {
	task := task.Task{
		Id:          taskId,
		Title:       taskTitle,
		Description: taskDesc,
		IsCompleted: false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_, err := pg.db.Exec("INSERT INTO tasks(id, title, description, is_completed, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)", task.Id, task.Title, task.Description, task.IsCompleted, task.CreatedAt, task.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &task, nil
}
