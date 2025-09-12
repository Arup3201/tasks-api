package storage

import (
	"testing"
)

func TestNewPostgres(t *testing.T) {
	_, err := NewPostgres()

	if err != nil {
		t.Errorf("NewPotgres test failed with error: %v", err)
	}
}

func TestPostgresGet(t *testing.T) {
	p, _ := NewPostgres()
	testTaskId := "TASK_TEST1"

	task, err := p.Get(testTaskId)

	if err != nil {
		t.Errorf("TestPostgresGet failed with error: %v", err)
	}
	if task.Id != testTaskId {
		t.Errorf("Postgres.Get(%s)={Id: %s, ...}, expected %s", testTaskId, task.Id, testTaskId)
	}
}

func TestPostgresInsert(t *testing.T) {
	p, _ := NewPostgres()
	title, description := "Task for Testing", "This task is for testing purpose"

	taskId, err := p.Insert(title, description)

	if err != nil {
		t.Errorf("TestPostgresInsert failed with error: %v", err)
	}
	_, err = p.Get(taskId)
	if err != nil {
		t.Errorf("Postgres.Insert failed to insert")
	}
}

func TestPostgresUpdate(t *testing.T) {
	p, _ := NewPostgres()
	testTaskId := "TASK_TEST1"
	newTitle := "Test Task 1 (edited 3 times) Title"

	task, err := p.Update(testTaskId, &newTitle, nil, nil)

	if err != nil {
		t.Errorf("TestPostgresUpdate failed with error: %v", err)
	}
	if err != nil {
		t.Errorf("Postgres.Update failed to Update")
	}
	if task.Title != newTitle {
		t.Errorf("Postgres.Update => %s, expected %s", task.Title, newTitle)
	}
}

func TestPostgresDelete(t *testing.T) {
	p, _ := NewPostgres()
	testTaskId := "TASK_TEST1"

	deletedId, err := p.Delete(testTaskId)

	if err != nil {
		t.Errorf("TestPostgresDelete failed with error: %v", err)
	}
	if deletedId != testTaskId {
		t.Errorf("p.Delete(%s)=%s, expected %s", testTaskId, deletedId, testTaskId)
	}
}

func TestPostgresList(t *testing.T) {
	p, _ := NewPostgres()

	tasks, err := p.List()

	if err != nil {
		t.Errorf("TestPostgresList failed with error: %v", err)
	}
	if tasks == nil {
		t.Errorf("p.List()=nil, expected Slice([])")
	}
}
