package handler

import (
	"e-commerce/backend/internal/middleware"
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/services"
	"e-commerce/backend/internal/utils"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Helper function untuk extract ID dari URL
func extractIDFromPath(path string) (uint, error) {
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

// GetUsers - GET /api/users
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 10
	}

	req := models.UserListRequest{
		Page:  page,
		Limit: limit,
	}

	users, err := h.userService.GetUsers(req)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Users retrieved successfully",
		"data":    users,
	})
}

// GetUserById - GET /api/users/:id
func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(r.URL.Path)
	if err != nil || id == 0 {
		utils.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.userService.GetUserById(id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "User retrieved successfully",
		"data":    user,
	})
}

// GetCurrentUser - GET /api/users/me
func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserIDFromContext(r)
	if userID == 0 {
		utils.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	user, err := h.userService.GetUserById(userID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Current user retrieved successfully",
		"data":    user,
	})
}

// CreateUser - POST /api/users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.UserInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.userService.CreateUser(&req)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "User created successfully",
		"data":    user,
	})
}

// UpdateUser - PUT /api/users/:id
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(r.URL.Path)
	if err != nil || id == 0 {
		utils.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req models.UserUpdateInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.userService.UpdateUser(id, &req)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "User updated successfully",
		"data":    user,
	})
}

// DeleteUser - DELETE /api/users/:id
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(r.URL.Path)
	if err != nil || id == 0 {
		utils.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	err = h.userService.DeleteUser(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "User deleted successfully",
	})
}

// ActivateUser - PUT /api/users/:id/activate
func (h *UserHandler) ActivateUser(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(r.URL.Path)
	if err != nil || id == 0 {
		utils.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.userService.ActivateUser(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "User activated successfully",
		"data":    user,
	})
}

// DeactivateUser - PUT /api/users/:id/deactivate
func (h *UserHandler) DeactivateUser(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(r.URL.Path)
	if err != nil || id == 0 {
		utils.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.userService.DeactivateUser(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "User deactivated successfully",
		"data":    user,
	})
}

// BulkUserActions - POST /api/users/bulk
func (h *UserHandler) BulkUserActions(w http.ResponseWriter, r *http.Request) {
	var req services.BulkActionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.userService.BulkUserActions(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Bulk action completed successfully",
	})
}

// UpdatePassword - PUT /api/users/:id/password
func (h *UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(r.URL.Path)
	if err != nil || id == 0 {
		utils.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req models.PasswordUpdateInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = h.userService.UpdatePassword(id, &req)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Password updated successfully",
	})
}

// UpdateCurrentUserPassword - PUT /api/users/me/password
func (h *UserHandler) UpdateCurrentUserPassword(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID from context
	userID := middleware.GetUserIDFromContext(r)
	if userID == 0 {
		utils.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req models.PasswordUpdateInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.userService.UpdatePassword(userID, &req)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Password updated successfully",
	})
}
