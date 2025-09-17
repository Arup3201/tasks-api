package storages

import "github.com/Arup3201/gotasks/internal/entities/task"

type TaskRepository interface {
	Get(taskId string) *task.Task
	Insert(taskTitle, taskDesc string) (*task.Task, error)
	Update(taskId string, data map[string]any) *task.Task
	List() []task.Task
}
