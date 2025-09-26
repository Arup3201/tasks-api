package storages

import (
	"database/sql"
	"fmt"

	"github.com/Arup3201/gotasks/internal/entities/task"
	postgres "github.com/Arup3201/gotasks/internal/storages/postgres/task"
	. "github.com/Arup3201/gotasks/internal/utils"
	_ "github.com/lib/pq"
)

const (
	Postgres = "Postgres"
)

func New(dbType string) (TaskRepository, error) {
	var repo TaskRepository
	switch dbType {
	case Postgres:
		db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", Config.DBUser, Config.DBPass, Config.DBHost, Config.DBName))
		if err != nil {
			return nil, fmt.Errorf("sql.Open error: %v", err)
		}
		repo = postgres.NewPgTaskRepository(db)
	default:
		// in-memory
	}

	return repo, nil
}

type TaskRepository interface {
	Get(taskId int) (*task.Task, error)
	Insert(taskId int, taskTitle, taskDesc string) (*task.Task, error)
	Update(taskId int, data map[string]any) (*task.Task, error)
	Delete(taskId int) (*int, error)
	List() ([]task.Task, error)
	Close() error
}
