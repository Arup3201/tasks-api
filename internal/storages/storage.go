package storages

import "github.com/Arup3201/gotasks/internal/entities/task"

type TaskRepository interface {
	Get(taskId int) (*task.Task, error)
	Insert(taskId int, taskTitle, taskDesc string) (*task.Task, error)
	Update(taskId int, data map[string]any) (*task.Task, error)
	List() ([]task.Task, error)
}
