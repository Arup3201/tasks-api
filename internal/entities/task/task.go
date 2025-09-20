package task

import "time"

type Task struct {
	Id          int
	Title       string
	Description string
	IsCompleted bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
