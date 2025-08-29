package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Arup3201/gotasks/internal/storage"
)

func GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := storage.GetAll("data/tasks.json")
	if err != nil {
		http.Error(w, "Failed to open tasks.json", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
