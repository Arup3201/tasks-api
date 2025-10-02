package task

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Arup3201/gotasks/internal/entities/task"
	"github.com/Arup3201/gotasks/internal/errors"
)

type PgTaskRepository struct {
	db *sql.DB
}

func NewPgTaskRepository(db *sql.DB) *PgTaskRepository {
	return &PgTaskRepository{
		db: db,
	}
}

func (pg *PgTaskRepository) Get(taskId int) (*task.Task, error) {
	var task task.Task
	if err := pg.db.QueryRow("SELECT * FROM tasks WHERE id = ($1)", taskId).Scan(&task.Id, &task.Title, &task.Description, &task.IsCompleted, &task.CreatedAt, &task.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NotFoundError(fmt.Sprintf("Task with ID %d not found", taskId))
		}
		return nil, err
	}
	return &task, nil
}

func (pg *PgTaskRepository) Insert(taskId int, taskTitle, taskDesc string) (*task.Task, error) {
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

func (pg *PgTaskRepository) Update(taskId int, data map[string]any) (*task.Task, error) {
	_, err := pg.Get(taskId)
	if err != nil {
		return nil, err
	}

	setFields := []string{}

	title, ok := data["Title"]
	if ok {
		setFields = append(setFields, fmt.Sprintf("title='%s'", title))
	}
	description, ok := data["Description"]
	if ok {
		setFields = append(setFields, fmt.Sprintf("description='%s'", description))
	}
	isCompleted, ok := data["IsCompleted"]
	if ok {
		setFields = append(setFields, fmt.Sprintf("is_completed=%t", isCompleted))
	}

	if len(setFields) == 0 {
		return nil, errors.NoOp("Found no fields to update")
	}

	execString := fmt.Sprintf("UPDATE tasks SET %s WHERE id=($1)", strings.Join(setFields, ", "))
	_, err = pg.db.Exec(execString, taskId)
	if err != nil {
		return nil, err
	}

	return pg.Get(taskId)
}

func (pg *PgTaskRepository) Delete(taskId int) (*int, error) {
	res, err := pg.db.Exec("DELETE FROM tasks WHERE id=($1)", taskId)
	if err != nil {
		return nil, err
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return nil, errors.NotFoundError(fmt.Sprintf("Task with ID %d not found", taskId))
	}

	return &taskId, nil
}

func (pg *PgTaskRepository) List() ([]task.Task, error) {
	var tasks []task.Task
	rows, err := pg.db.Query("SELECT * FROM tasks")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		t := &task.Task{}
		if err := rows.Scan(&t.Id, &t.Title, &t.Description, &t.IsCompleted, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, *t)
	}

	return tasks, nil
}

func (pg *PgTaskRepository) Close() error {
	return pg.db.Close()
}
