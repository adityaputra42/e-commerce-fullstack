package handler

import (
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/services"
	"e-commerce/backend/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type RoleHandler struct {
	roleService services.RoleService
}

func NewRoleHandler(roleService services.RoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

// GetAllRoles - GET /api/v1/roles
// @Summary List all roles
// @Description Get all available user roles
// @Tags Role
// @Produce json
// @Success 200 {object} utils.Response{data=[]models.Role} "Success"
// @Router /roles [get]
// @Security Bearer
func (h *RoleHandler) GetAllRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := h.roleService.FindAllRole()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch roles", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Roles retrieved successfully", roles)
}

// GetRoleByID - GET /api/v1/roles/{id}
// @Summary Get role by ID
// @Description Get detail info about a role
// @Tags Role
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} utils.Response{data=models.Role} "Success"
// @Router /roles/{id} [get]
// @Security Bearer
func (h *RoleHandler) GetRoleByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid role ID", err)
		return
	}

	role, err := h.roleService.FindById(uint(id))
	if err != nil {
		if err.Error() == "role not found" {
			utils.WriteError(w, http.StatusNotFound, "Role not found", err)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch role", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Role retrieved successfully", role)
}

// CreateRole - POST /api/v1/roles
// @Summary Create a new role
// @Description Create a new system role
// @Tags Role
// @Accept json
// @Produce json
// @Param request body models.RoleInput true "Role request"
// @Success 201 {object} utils.Response{data=models.Role} "Role created successfully"
// @Router /roles [post]
// @Security Bearer
func (h *RoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var input models.RoleInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if input.Name == "" {
		utils.WriteError(w, http.StatusBadRequest, "Role name is required", fmt.Errorf("Role name is required"))
		return
	}

	role, err := h.roleService.CreateRole(&input)
	if err != nil {
		if err.Error() == "role with this name already exists" {
			utils.WriteError(w, http.StatusBadRequest, err.Error(), err)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create role", err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, "Role created successfully", role)
}

// @Router /roles/{id} [put]
func (h *RoleHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid role ID", err)
		return
	}

	var input models.RoleInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	role, err := h.roleService.UpdateRole(uint(id), &input)
	if err != nil {
		if err.Error() == "role not found" {
			utils.WriteError(w, http.StatusNotFound, "Role not found", err)
			return
		}
		if err.Error() == "cannot change system role name" {
			utils.WriteError(w, http.StatusBadRequest, err.Error(), err)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update role", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Role updated successfully", role)
}

// @Router /roles/{id} [delete]
func (h *RoleHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid role ID", err)
		return
	}

	err = h.roleService.DeleteRole(uint(id))
	if err != nil {
		if err.Error() == "role not found" {
			utils.WriteError(w, http.StatusNotFound, "Role not found", err)
			return
		}
		if err.Error() == "cannot delete system role" || err.Error() == "cannot delete role that is assigned to users" {
			utils.WriteError(w, http.StatusBadRequest, err.Error(), err)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete role", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Role deleted successfully", nil)
}

// @Router /roles/permissions [get]
func (h *RoleHandler) GetAllPermissions(w http.ResponseWriter, r *http.Request) {
	permissions, err := h.roleService.GetPermissions()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch permissions", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Permissions retrieved successfully", permissions)
}

// @Router /roles/{id}/permissions [get]
func (h *RoleHandler) GetRolePermissions(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid role ID", err)
		return
	}

	permissions, err := h.roleService.GetRolePermissions(uint(id))
	if err != nil {
		if err.Error() == "role not found" {
			utils.WriteError(w, http.StatusNotFound, "Role not found", err)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch role permissions", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Role permissions retrieved successfully", permissions)
}

// @Router /roles/{id}/permissions [post]
func (h *RoleHandler) AssignPermissions(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid role ID", err)
		return
	}

	var input struct {
		PermissionIDs []uint `json:"permission_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if len(input.PermissionIDs) == 0 {
		utils.WriteError(w, http.StatusBadRequest, "Permission IDs are required", err)
		return
	}

	err = h.roleService.AssignPermissions(uint(id), input.PermissionIDs)
	if err != nil {
		if err.Error() == "role not found" {
			utils.WriteError(w, http.StatusNotFound, "Role not found", err)
			return
		}
		if err.Error() == "some permissions not found" {
			utils.WriteError(w, http.StatusBadRequest, err.Error(), err)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to assign permissions", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Permissions assigned successfully", nil)
}
