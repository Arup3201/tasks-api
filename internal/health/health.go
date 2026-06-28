package health

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"gorm.io/gorm"
)

var startTime = time.Now()

type HealthChecker struct {
	db *gorm.DB
}

func NewHealthChecker(db *gorm.DB) *HealthChecker {
	return &HealthChecker{db}
}

func (h *HealthChecker) HealthHandler(w http.ResponseWriter, r *http.Request) {
	status := h.checkHealth()

	w.Header().Set("Content-Type", "application/json")

	if status.Status == "healthy" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	json.NewEncoder(w).Encode(status)
}

type HealthStatus struct {
	Status string  `json:"status"`
	Uptime string  `json:"uptime"`
	Checks []Check `json:"checks"`
}

type Check struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Latency string `json:"latency"`
}

func (h *HealthChecker) checkHealth() HealthStatus {

	dbCheck := h.checkDatabase()

	status := "healthy"
	if dbCheck.Status != "healthy" {
		status = "unhealthy"
	}

	return HealthStatus{
		Status: status,
		Uptime: time.Since(startTime).String(),
		Checks: []Check{dbCheck},
	}
}

func (h *HealthChecker) checkDatabase() Check {
	start := time.Now()

	sql, err := h.db.DB()
	if err != nil {
		return Check{
			Status:  "unhealthy",
			Message: "Database is missing",
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	err = sql.PingContext(ctx)
	defer cancel()
	if err != nil {
		return Check{
			Status:  "unhealthy",
			Message: err.Error(),
		}
	}

	return Check{
		Status:  "healthy",
		Message: "Database is healthy",
		Latency: time.Since(start).String(),
	}
}
