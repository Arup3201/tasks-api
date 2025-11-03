package storages

import (
	"database/sql"
	"fmt"
	"log"

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

		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tasks(
								id VARCHAR(256) PRIMARY KEY,
								title TEXT NOT NULL, 
								description TEXT NOT NULL, 
								is_completed BOOLEAN NOT NULL, 
								created_at TIMESTAMP WITH TIME ZONE NOT NULL,
								updated_at TIMESTAMP WITH TIME ZONE NOT NULL
							)`)
		if err != nil {
			log.Fatalf("Table create error: %v", err)
		}

		repo = postgres.NewPgTaskRepository(db)
	default:
		// in-memory
	}

	return repo, nil
}

type TaskRepository interface {
	Get(taskId string) (*task.Task, error)
	Insert(taskId string, taskTitle, taskDesc string) (*task.Task, error)
	Update(taskId string, data map[string]any) (*task.Task, error)
	Delete(taskId string) (*string, error)
	List() ([]task.Task, error)
	Close() error
}
