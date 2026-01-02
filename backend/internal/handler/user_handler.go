package handler

import (
	"e-commerce/backend/internal/middleware"
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/services"
	"e-commerce/backend/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUsers - GET /api/v1/users
// @Summary List users
// @Description Get a paginated list of users
// @Tags User
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} utils.Response{data=[]models.UserResponse} "Success"
// @Router /users [get]
// @Security Bearer
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
		utils.WriteError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Users retrieved successfully", users)
}

// GetUserById - GET /api/v1/users/{id}
// @Summary Get user by ID
// @Description Get detailed information about a user
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} utils.Response{data=models.UserResponse} "Success"
// @Failure 404 {object} utils.Response "User not found"
// @Router /users/{id} [get]
// @Security Bearer
func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ExtractIDFromPath(r.URL.Path)
	if err != nil || id == 0 {
		utils.WriteError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	user, err := h.userService.GetUserById(id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "User retrieved successfully", user)
}

// GetCurrentUser - GET /api/users/me
func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserIDFromContext(r)
	if userID == 0 {
		utils.WriteError(w, http.StatusUnauthorized, "User not authenticated", fmt.Errorf("User not authenticated"))
		return
	}

	user, err := h.userService.GetUserById(userID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Current user retrieved successfully", user)
}

// CreateUser - POST /api/v1/users
// @Summary Create a new user
// @Description Create a new user with provided details
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.UserInput true "User input"
// @Success 201 {object} utils.Response{data=models.UserResponse} "User created successfully"
// @Router /users [post]
// @Security Bearer
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.UserInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	user, err := h.userService.CreateUser(&req)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, "User created successfully", user)
}

// UpdateUser - PUT /api/users/:id
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ExtractIDFromPath(r.URL.Path)
	if err != nil || id == 0 {
		utils.WriteError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	var req models.UserUpdateInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	user, err := h.userService.UpdateUser(id, &req)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "User updated successfully", user)
}

// DeleteUser - DELETE /api/users/:id
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ExtractIDFromPath(r.URL.Path)
	if err != nil || id == 0 {
		utils.WriteError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	err = h.userService.DeleteUser(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "User deleted successfully", nil)
}

// ActivateUser - PUT /api/users/:id/activate
func (h *UserHandler) ActivateUser(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ExtractIDFromPath(r.URL.Path)
	if err != nil || id == 0 {
		utils.WriteError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	user, err := h.userService.ActivateUser(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "User activated successfully", user)
}

// DeactivateUser - PUT /api/users/:id/deactivate
func (h *UserHandler) DeactivateUser(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ExtractIDFromPath(r.URL.Path)
	if err != nil || id == 0 {
		utils.WriteError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	user, err := h.userService.DeactivateUser(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "User deactivated successfully", user)
}

// BulkUserActions - POST /api/users/bulk
func (h *UserHandler) BulkUserActions(w http.ResponseWriter, r *http.Request) {
	var req services.BulkActionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := h.userService.BulkUserActions(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Bulk action completed successfully", nil)
}

// UpdatePassword - PUT /api/users/:id/password
func (h *UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ExtractIDFromPath(r.URL.Path)
	if err != nil || id == 0 {
		utils.WriteError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	var req models.PasswordUpdateInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	err = h.userService.UpdatePassword(id, &req)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Password updated successfully", nil)
}

// UpdateCurrentUserPassword - PUT /api/users/me/password
func (h *UserHandler) UpdateCurrentUserPassword(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID from context
	userID := middleware.GetUserIDFromContext(r)
	if userID == 0 {
		utils.WriteError(w, http.StatusUnauthorized, "User not authenticated", fmt.Errorf("User not authenticated"))
		return
	}

	var req models.PasswordUpdateInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	err := h.userService.UpdatePassword(userID, &req)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Password updated successfully", nil)
}
