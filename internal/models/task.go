package models

import "time"

type Task struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateTask struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

type UpdateTask struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Completed   *bool   `json:"completed,omitempty"`
}

type EditTask struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
}

type MarkTask struct {
	Completed *bool `json:"completed,omitempty"`
}
