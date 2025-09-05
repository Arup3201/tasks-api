package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Arup3201/gotasks/internal/models"
)

func CreateStorage(file string) (*os.File, error) {
	f, err := os.Create(file)
	if err != nil {
		return nil, fmt.Errorf("failed to create file %s: %v", file, err)
	}

	_, err = f.Write([]byte{'[',']'})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize the file %s: %v", file, err)
	}

	return f, nil
}

func RemoveStorage(file string) (error) {
	err := os.Remove(file)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", file, err)
	}

	return nil
}

func GetAll(storage string) ([]models.Task, error) {
	data, err := os.ReadFile(storage)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %v", storage, err)
	}

	allTasks := []models.Task{}
	err = json.Unmarshal(data, &allTasks)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s: %v", storage, err)
	}

	return allTasks, nil
}

func Get(storage, id string) (models.Task, error) {
	data, err := os.ReadFile(storage)
	if err != nil {
		return models.Task{}, fmt.Errorf("failed to open %s: %v", storage, err)
	}

	allTasks := []models.Task{}
	err = json.Unmarshal(data, &allTasks)
	if err != nil {
		return models.Task{}, fmt.Errorf("failed to unmarshal %s: %v", storage, err)
	}

	for _, task := range allTasks {
		if task.Id == id {
			return task, nil
		}
	}

	return models.Task{}, fmt.Errorf("storage get failed")
}

func Save(tasks []models.Task, storage string) error {
	data, err := json.Marshal(tasks)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	err = os.WriteFile(storage, data, 0666)
	if err != nil {
		return fmt.Errorf("failed to write data to %s: %v", storage, err)
	}

	return nil
}

func Add(task models.Task, storage string) (error) {
	allTasks, err := GetAll(storage)
	if err != nil {
		return fmt.Errorf("failed to get all tasks: %v", err)
	}

	allTasks = append(allTasks, task)
	err = Save(allTasks, storage)
	if err != nil {
		return fmt.Errorf("failed to save data: %v", err)
	}

	return nil
}

func Edit(id string, data models.Task, storage string) (error) {
	allTasks, err := GetAll(storage)
	if err != nil {
		return fmt.Errorf("failed to get all tasks: %v", err)
	}
	
	for i, task := range allTasks {
		if task.Id == id {
			task.Title = data.Title
			task.Description = data.Description
			task.Completed = data.Completed
			task.CreatedAt = data.CreatedAt
			task.UpdatedAt = time.Now()

			allTasks[i] = task
		}
	}

	err = Save(allTasks, storage)
	if err != nil {
		return fmt.Errorf("failed to save data: %v", err)
	}

	return nil
}

func Delete(id, storage string) (error) {
	allTasks, err := GetAll(storage)
	if err != nil {
		return fmt.Errorf("failed to get all tasks: %v", err)
	}
	
	newTasks := []models.Task{}
	for _, task := range allTasks {
		if task.Id == id {
			continue
		}
		newTasks = append(newTasks, task)
	}

	err = Save(newTasks, storage)
	if err != nil {
		return fmt.Errorf("failed to save data: %v", err)
	}

	return nil
}
