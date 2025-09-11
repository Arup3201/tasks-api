package models

type FieldName string

type TaskStore interface {
	Get(id string) (Task, error)                                                 // get task with id if something wrong return error
	Insert(title, description string) (string, error)                            // create task with data and return the created task id but on error return that
	Update(id string, title, description *string, completed *bool) (Task, error) // update task with id by data and return the updated task in case of error return that error
	Delete(id string) (string, error)                                            // delete task with id and return the deleted task id, but on error return that
	List() ([]Task, error)                                                       // return list of all tasks, on error return that error
	Search(by FieldName, query string) ([]Task, error)                           // search tasks with `query` on field `by`
}
