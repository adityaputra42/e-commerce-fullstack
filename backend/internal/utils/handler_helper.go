package utils

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func WriteError(w http.ResponseWriter, statusCode int, message string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}

	json.NewEncoder(w).Encode(Response{
		Success: false,
		Message: message,
		Error:   errorMsg,
	})
}
func ExtractIDFromPath(path string) (uint, error) {
	pathParts := strings.Split(strings.Trim(path, "/"), "/")
	if len(pathParts) < 2 {
		return 0, nil
	}

	idStr := pathParts[len(pathParts)-1]
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}
