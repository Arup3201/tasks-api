package storage

import (
	"os"
	"testing"
	"time"

	"github.com/Arup3201/gotasks/internal/models"
	"github.com/Arup3201/gotasks/internal/storage"
)

func TestCreateStorage(t *testing.T) {
	testFile := "../data/test.json"
	_, err := storage.CreateStorage(testFile)
	if err != nil {
		t.Errorf("storage.CreateStorage failed: %v", err)
	}

	_, err = os.Open("../data/test.json")
	if err != nil {
		t.Errorf("storage.CreateStorage did not create file")
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

	err := storage.Add(testTask, testFile)
	if err != nil {
		t.Errorf("storage.Add() failed to add: %v", err)
	}

	tasks, _ := storage.GetAll(testFile)
	saved := false
	for _, task := range tasks {
		if task.Id == testTask.Id {
			saved = true
		}
	}

	if !saved {
		t.Errorf("storage.Add() did not save")
	}
}

func TestRemoveStorage(t *testing.T) {
	err := storage.RemoveStorage("../data/test.json")
	if err != nil {
		t.Errorf("storage.RemoveStorage failed: %v", err)
	}

	_, err = os.Open("../data/test.json")
	if err==nil {
		t.Errorf("storage.RemoveStorage did not remove the storage file")
	}
}
