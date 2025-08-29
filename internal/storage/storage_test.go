package storage

import (
	"os"
	"testing"
	"time"

	"github.com/Arup3201/gotasks/internal/models"
)

func TestCreateStorage(t *testing.T) {
	testFile := "../data/test.json"
	_, err := CreateStorage(testFile)
	if err != nil {
		t.Errorf("CreateStorage failed: %v", err)
	}

	_, err = os.Open("../data/test.json")
	if err != nil {
		t.Errorf("CreateStorage did not create file")
	}
}

func TestAdd(t *testing.T) {
	testFile := "../data/test.json"
	testTask := models.Task{
		Id: "T001",
		Title: "Golang practice",
		Description: "Practice golang for 20 mins - focusing on goroutines, and testing",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := Add(testTask, testFile)
	if err != nil {
		t.Errorf("Add() failed to add: %v", err)
	}

	tasks, _ := GetAll(testFile)
	saved := false
	for _, task := range tasks {
		if task.Id == testTask.Id {
			saved = true
		}
	}

	if !saved {
		t.Errorf("Add() did not save")
	}
}

func TestDelete(t *testing.T) {
	testFile := "../data/test.json"
	testTaskId := "T002"

	err := Delete(testTaskId, testFile)
	if err != nil {
		t.Errorf("delete failed: %v", err)
	}

	tasks, _ := GetAll(testFile)
	for _, task := range tasks {
		if task.Id == testTaskId {
			t.Errorf("task was not deleted")
		}
	}
}

func TestRemoveStorage(t *testing.T) {
	err := RemoveStorage("../data/test.json")
	if err != nil {
		t.Errorf("RemoveStorage failed: %v", err)
	}

	_, err = os.Open("../data/test.json")
	if err==nil {
		t.Errorf("RemoveStorage did not remove the file")
	}
}
