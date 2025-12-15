package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Version   string    `json:"version"`
}

// HealthCheck handler - support GET and HEAD method
func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Support both GET and HEAD methods
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
		Service:   "e-commerce-api",
		Version:   "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Only write body for GET requests (HEAD should only return headers)
	if r.Method == http.MethodGet {
		json.NewEncoder(w).Encode(response)
	}
}

// Root handler
func (h *HealthHandler) Root(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "E-Commerce API Server",
		"status":  "running",
		"time":    time.Now().Format(time.RFC3339),
		"version": "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
