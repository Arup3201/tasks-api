package task

import (
	"database/sql/driver"
	"fmt"
	"testing"
	"time"

	"github.com/Arup3201/gotasks/internal/errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestPgGet(t *testing.T) {
	t.Run("should get task", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("sqlmock.New error: %v", err)
		}
		defer db.Close()
		uuid, _ := uuid.NewUUID()
		id := uuid.String()
		title, description := "Test task", "Test task description"
		rows := sqlmock.NewRows([]string{"id", "title", "description", "is_completed", "created_at", "updated_at"}).AddRow(id, title, description, false, time.Now(), time.Now())
		mock.ExpectQuery("^SELECT (.+) FROM tasks").WithArgs(id).WillReturnRows(rows)
		pg := NewPgTaskRepository(db)

		task, err := pg.Get(id)

		if err != nil {
			t.Errorf("pg.Get error: %v", err)
		}
		if task.Id != id {
			t.Errorf("expected task ID %s, but got %s", id, task.Id)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
	t.Run("should fail to get task", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("sqlmock.New error: %v", err)
		}
		defer db.Close()
		uuid_, _ := uuid.NewUUID()
		id := uuid_.String()
		mock.ExpectQuery("^SELECT (.+) FROM tasks").WithArgs(id).WillReturnError(fmt.Errorf("ErrNoRows"))
		uuid_, _ = uuid.NewUUID()
		exid := uuid_.String()
		title, description := "Test task", "Test task description"
		pg := NewPgTaskRepository(db)
		pg.Insert(exid, title, description)

		_, err = pg.Get(id)

		if err == nil {
			t.Errorf("expected error but got nothing")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestPgInsert(t *testing.T) {
	t.Run("should create a task", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock sql db: %v", err)
		}
		defer db.Close()
		uuid_, _ := uuid.NewUUID()
		id := uuid_.String()
		title := "Test task"
		description := "Test task description"
		mock.ExpectExec("INSERT INTO tasks").WithArgs(id, title, description, false, AnyTime{}, AnyTime{}).WillReturnResult(sqlmock.NewResult(1, 1))
		pg := NewPgTaskRepository(db)

		task, err := pg.Insert(id, title, description)

		if err != nil {
			t.Errorf("Insert failed with error: %v", err)
			return
		}
		if task.Title != title || task.Description != description {
			t.Errorf("Inserted task expected title and description %s, %s but got %s, %s", task.Title, task.Description, title, description)
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
	t.Run("fail to create task with duplicate ID", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock sql db: %v", err)
		}
		defer db.Close()
		uuid_, _ := uuid.NewUUID()
		id := uuid_.String()
		title := "Test task 2"
		description := "Test task 2 description"
		mock.ExpectExec("INSERT INTO tasks").WithArgs(id, title, description, false, AnyTime{}, AnyTime{}).WillReturnError(fmt.Errorf("DB integrity error"))
		pg := NewPgTaskRepository(db)
		pg.Insert(id, "Test task 1", "Test task 1 description")

		_, err = pg.Insert(id, title, description)

		if err == nil {
			t.Errorf("expecting an error, but there was none")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestPgUpdate(t *testing.T) {
	t.Run("update task success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("sqlmock.New error: %v", err)
		}
		defer db.Close()
		uuid_, _ := uuid.NewUUID()
		id := uuid_.String()
		title, description := "Test task", "Test task description"
		row := sqlmock.NewRows([]string{"id", "title", "description", "is_completed", "created_at", "updated_at"}).AddRow(id, title, description, false, time.Now(), time.Now())
		updateTitle := "Test task (updated)"
		mock.ExpectQuery("SELECT (.+) FROM tasks").WithArgs(id).WillReturnRows(row)
		mock.ExpectExec("UPDATE tasks").WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectQuery("SELECT (.+) FROM tasks").WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "is_completed", "created_at", "updated_at"}).AddRow(id, updateTitle, description, false, time.Now(), time.Now()))
		pg := NewPgTaskRepository(db)

		task, err := pg.Update(id, map[string]any{
			"Title": "Test task (updated)",
		})

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if task.Title != updateTitle {
			t.Errorf("task is not updated, expected %s but got %s", updateTitle, task.Title)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestPgDelete(t *testing.T) {
	t.Run("delete a task", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("sqlmock.New error: %v", err)
		}
		defer db.Close()
		uuid_, _ := uuid.NewUUID()
		id := uuid_.String()
		sqlmock.NewRows([]string{"id", "title", "description", "is_completed", "created_at", "updated_at"}).AddRow(id, "Test task 1", "Test task 1 description", false, time.Now(), time.Now()).AddRow(2, "Test task 2", "Test task 2 description", true, time.Now(), time.Now()).AddRow(3, "Test task 3", "Test task 3 description", false, time.Now(), time.Now())
		mock.ExpectExec("DELETE FROM tasks").WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 1))
		pg := NewPgTaskRepository(db)

		dId, err := pg.Delete(id)

		if err != nil {
			t.Errorf("error occured: %v", err)
			return
		}
		if *dId != id {
			t.Errorf("expected deleted task %s but got %s", id, *dId)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
	t.Run("delete an invalid task fail", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("sqlmock.New error: %v", err)
		}
		defer db.Close()
		uuid_, _ := uuid.NewUUID()
		id := uuid_.String()
		sqlmock.NewRows([]string{"id", "title", "description", "is_completed", "created_at", "updated_at"}).AddRow(1, "Test task 1", "Test task 1 description", false, time.Now(), time.Now()).AddRow(2, "Test task 2", "Test task 2 description", true, time.Now(), time.Now()).AddRow(3, "Test task 3", "Test task 3 description", false, time.Now(), time.Now())
		mock.ExpectExec("DELETE FROM tasks").WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 0))
		pg := NewPgTaskRepository(db)

		_, err = pg.Delete(id)

		if err == nil {
			t.Errorf("expected not found error, but got no error")
		}
		if _, ok := err.(*errors.AppError); !ok {
			t.Errorf("expected not found error, but got some other error: %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestPgList(t *testing.T) {
	t.Run("list all tasks", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("sqlmock.New error: %v", err)
		}
		defer db.Close()
		rows := sqlmock.NewRows([]string{"id", "title", "description", "is_completed", "created_at", "updated_at"}).AddRow(1, "Test task 1", "Test task 1 description", false, time.Now(), time.Now()).AddRow(2, "Test task 2", "Test task 2 description", true, time.Now(), time.Now()).AddRow(3, "Test task 3", "Test task 3 description", false, time.Now(), time.Now())
		mock.ExpectQuery("^SELECT (.+) FROM tasks$").WillReturnRows(rows)
		pg := NewPgTaskRepository(db)

		tasks, err := pg.List()

		if err != nil {
			t.Errorf("error occured: %v", err)
			return
		}
		if tasks == nil {
			t.Errorf("expected list of tasks but got nothing")
			return
		}
		if len(tasks) == 0 {
			t.Errorf("expected list of tasks but got empty list")
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
