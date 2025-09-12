package storage

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Arup3201/gotasks/internal/handlers/apperr"
	"github.com/Arup3201/gotasks/internal/models"
	"github.com/Arup3201/gotasks/internal/utils"
	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
}

func (p *Postgres) Get(id string) (models.Task, error) {
	var task models.Task

	row := p.db.QueryRow("SELECT * FROM tasks WHERE id = $1", id)
	if err := row.Scan(&task.Id, &task.Title, &task.Description, &task.Completed, &task.CreatedAt, &task.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return task, apperr.NotFoundError()
		}
		return task, fmt.Errorf("Postgres.Get %s: %v", id, err)
	}

	return task, nil
}
func (p *Postgres) Insert(title, description string) (string, error) {
	taskId := utils.GenerateID("TASK_")
	result, err := p.db.Exec("INSERT INTO tasks (id, title, description, completed, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)", taskId, title, description, false, time.Now(), time.Now())
	if err != nil {
		return "", fmt.Errorf("Postgres.Insert: %v", err)
	}

	if rows, _ := result.RowsAffected(); rows < 1 {
		return "", fmt.Errorf("Postgres.Insert failed to insert row")
	}

	return taskId, nil
}
func (p *Postgres) Update(id string, title, description *string, completed *bool) (models.Task, error) {
	task, err := p.Get(id)
	if err != nil {
		return task, err
	}

	if title == nil && description == nil && completed == nil {
		return task, fmt.Errorf("Postgres.Update error: expected any one of 'title', 'description' or 'completed'")
	}

	if title != nil {
		result, err := p.db.Exec("UPDATE tasks SET title = $1 WHERE id = $2", *title, id)
		if err != nil {
			return task, fmt.Errorf("Postgres.Update: %v", err)
		}

		if rows, _ := result.RowsAffected(); rows < 1 {
			return task, fmt.Errorf("Postgres.Update failed to update title column for row id %s", id)
		}
	}

	if description != nil {
		result, err := p.db.Exec("UPDATE tasks SET description = $1 WHERE id = $2", *description, id)
		if err != nil {
			return task, fmt.Errorf("Postgres.Update: %v", err)
		}

		if rows, _ := result.RowsAffected(); rows < 1 {
			return task, fmt.Errorf("Postgres.Update failed to update description column for row id %s", id)
		}
	}

	if completed != nil {
		result, err := p.db.Exec("UPDATE tasks SET completed = $1 WHERE id = $2", *completed, id)
		if err != nil {
			return task, fmt.Errorf("Postgres.Update: %v", err)
		}

		if rows, _ := result.RowsAffected(); rows < 1 {
			return task, fmt.Errorf("Postgres.Update failed to update completed column for row id %s", id)
		}
	}

	editedTask, _ := p.Get(id)
	return editedTask, nil
}
func (p *Postgres) Delete(id string) (string, error) {
	result, err := p.db.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		return "", fmt.Errorf("Postgres.Delete: %v", err)
	}

	if rows, _ := result.RowsAffected(); rows < 1 {
		return "", fmt.Errorf("Postgres.Delete failed to delete row with id %s", id)
	}

	return id, nil
}
func (p *Postgres) List() ([]models.Task, error) {
	var tasks []models.Task

	rows, err := p.db.Query("SELECT * FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("Postgres.List: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Task

		if err := rows.Scan(&t.Id, &t.Title, &t.Description, &t.Completed, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("Postgres.List: %v", err)
		}

		tasks = append(tasks, t)
	}

	return tasks, nil
}
func (p *Postgres) Search(by models.FieldName, query string) ([]models.Task, error) {
	var tasks []models.Task

	rows, err := p.db.Query(fmt.Sprintf("SELECT * FROM tasks WHERE LOWER(%s) LIKE '%%%s%%'", by, strings.ToLower(query)))
	if err != nil {
		return nil, fmt.Errorf("Postgres.Search: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Task

		if err := rows.Scan(&t.Id, &t.Title, &t.Description, &t.Completed, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("Postgres.Search: %v", err)
		}

		tasks = append(tasks, t)
	}

	return tasks, nil
}

func NewPostgres() (*Postgres, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", utils.Config.DBUser, utils.Config.DBPass, utils.Config.DBHost, utils.Config.DBName))
	if err != nil {
		return nil, fmt.Errorf("NewPostgres failed to open database: %v", err)
	}

	return &Postgres{
		db: db,
	}, nil
}
