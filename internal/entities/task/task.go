package task

import "time"

type Task struct {
	Id          string
	Title       string
	Description string
	IsCompleted bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
