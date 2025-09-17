package storages

import "github.com/Arup3201/gotasks/internal/entities/task"

type TaskRepository interface {
	Get(taskId string) *task.Task
	Insert(taskTitle, taskDesc string) (task.Task, error)
}
