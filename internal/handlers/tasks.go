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

type CreatePayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func AddTask(w http.ResponseWriter, r *http.Request) {
	var payload CreatePayload
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
