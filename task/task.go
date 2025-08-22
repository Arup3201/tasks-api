package task

import (
	"time"
	"github.com/Arup3201/gotasks/db"
)

type Task struct {
	Title string
	Description string
	Deadline time.Time
	CreatedAt time.Time `json:"created_at"`
}

func (t Task) Save() {
	db.Add("tasks", t)
}

func NewTask(title, description string, deadline time.Time) Task {
	return Task{
		Title: title, 
		Description: description,
		Deadline: deadline,
		CreatedAt: time.Now(),
	}
}

func GetAllTasks() []Task {
	tasks := db.Get("tasks")
	return tasks
}

