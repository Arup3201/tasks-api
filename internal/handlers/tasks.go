package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Arup3201/gotasks/internal/models"
	"github.com/Arup3201/gotasks/internal/storage"
)

func GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := storage.GetAll("data/tasks.json")
	if err != nil {
		log.Fatalf("GetTask failed: %v", err)
		http.Error(w, "SERVER_ERROR", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tasks)
	if err != nil {
		log.Fatalf("GetTask failed: %v", err)
		http.Error(w, "SERVER_ERROR", http.StatusInternalServerError)
		return
	}
}

func AddTask(w http.ResponseWriter, r *http.Request) {
	var payload models.CreateTask
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Fatalf("AddTask failed: %v", err)
		http.Error(w, "BAD_REQUEST", http.StatusBadRequest)
		return
	}

	if payload.Title == "" || payload.Description == "" {
		log.Fatalf("payload is invalid")
		http.Error(w, "BAD_REQUEST", http.StatusBadRequest)
		return
	}

	task := models.Task{
		Id:          "T" + time.Time.String(time.Now()),
		Title:       payload.Title,
		Description: payload.Description,
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err = storage.Add(task, "data/tasks.json")
	if err != nil {
		log.Fatalf("AddTask failed: %v", err)
		http.Error(w, "SERVER_ERROR", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func EditTask(w http.ResponseWriter, r *http.Request) {
	var payload models.EditTask
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Fatalf("EditTask decode error: %v", err)
		http.Error(w, "SERVER_ERROR", http.StatusInternalServerError)
		return
	}

	if payload.Id == "" {
		log.Fatalf("EditTask id missing in payload")
		http.Error(w, "BAD_REQUEST", http.StatusBadRequest)
		return
	}

	if payload.Title==nil && payload.Description==nil && payload.Completed==nil {
		log.Fatalf("EditTask payload is empty")
		http.Error(w, "NOT_MODIFIED", http.StatusNotModified)
		return
	}

	task, err := storage.Get(payload.Id, "data/tasks.json")
	if err != nil {
		log.Fatalf("EditTask storage.Get failed: %v", err)
		http.Error(w, "SERVER_ERROR", http.StatusInternalServerError)
		return
	}

	if payload.Title!=nil && *payload.Title!= task.Title {
		task.Title = *payload.Title
	}
	if payload.Description!=nil && *payload.Description!=task.Description {
		task.Description = *payload.Description
	}
	if payload.Completed!=nil && *payload.Completed!=task.Completed {
		task.Completed = *payload.Completed
	}

	err = storage.Edit(payload.Id, task, "data/tasks.json")
	if err!=nil {
		log.Fatalf("EditTask storage.Edit error: %v", err)
		http.Error(w, "SERVER_ERROR", http.StatusInternalServerError)
		return
	}
}
