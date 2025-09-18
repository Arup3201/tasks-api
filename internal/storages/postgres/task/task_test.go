package task

import (
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestPgInsert(t *testing.T) {
	t.Run("should create a task", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock sql db: %v", err)
		}
		defer db.Close()
		pg := NewPgTaskRepository(db)
		id := uuid.New().String()
		title := "Test task"
		description := "Test task description"
		mock.ExpectExec("INSERT INTO tasks").WithArgs(id, title, description, false, AnyTime{}, AnyTime{}).WillReturnResult(sqlmock.NewResult(1, 1))

		task, err := pg.Insert(id, title, description)

		if err != nil {
			t.Errorf("Insert failed with error: %v", err)
			return
		}
		if task.Title != title || task.Description != description {
			t.Errorf("Inserted task expected title and description %s, %s but got %s, %s", task.Title, task.Description, title, description)
			return
		}
		if task.Id == "" {
			t.Errorf("Inserted task must have an Id")
			return
		}
		if task.CreatedAt.IsZero() || task.UpdatedAt.IsZero() {
			t.Errorf("Inserted task must have created_at and updated_at field")
		}
		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
