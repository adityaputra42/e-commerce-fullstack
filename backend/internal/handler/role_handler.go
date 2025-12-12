package handler

import (
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/services"
	"e-commerce/backend/internal/utils"
	"encoding/json"
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

// @Router /roles [get]
func (h *RoleHandler) GetAllRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := h.roleService.FindAllRole()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch roles")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Roles retrieved successfully",
		"data":    roles,
	})
}

// @Router /roles/{id} [get]
func (h *RoleHandler) GetRoleByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid role ID")
		return
	}

	role, err := h.roleService.FindById(uint(id))
	if err != nil {
		if err.Error() == "role not found" {
			utils.WriteError(w, http.StatusNotFound, "Role not found")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch role")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Role retrieved successfully",
		"data":    role,
	})
}

// @Router /roles [post]
func (h *RoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var input models.RoleInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if input.Name == "" {
		utils.WriteError(w, http.StatusBadRequest, "Role name is required")
		return
	}

	role, err := h.roleService.CreateRole(&input)
	if err != nil {
		if err.Error() == "role with this name already exists" {
			utils.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create role")
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Role created successfully",
		"data":    role,
	})
}

// @Router /roles/{id} [put]
func (h *RoleHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid role ID")
		return
	}

	var input models.RoleInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	role, err := h.roleService.UpdateRole(uint(id), &input)
	if err != nil {
		if err.Error() == "role not found" {
			utils.WriteError(w, http.StatusNotFound, "Role not found")
			return
		}
		if err.Error() == "cannot change system role name" {
			utils.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update role")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Role updated successfully",
		"data":    role,
	})
}

// @Router /roles/{id} [delete]
func (h *RoleHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid role ID")
		return
	}

	err = h.roleService.DeleteRole(uint(id))
	if err != nil {
		if err.Error() == "role not found" {
			utils.WriteError(w, http.StatusNotFound, "Role not found")
			return
		}
		if err.Error() == "cannot delete system role" || err.Error() == "cannot delete role that is assigned to users" {
			utils.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete role")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Role deleted successfully",
	})
}

// @Router /roles/permissions [get]
func (h *RoleHandler) GetAllPermissions(w http.ResponseWriter, r *http.Request) {
	permissions, err := h.roleService.GetPermissions()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch permissions")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Permissions retrieved successfully",
		"data":    permissions,
	})
}

// @Router /roles/{id}/permissions [get]
func (h *RoleHandler) GetRolePermissions(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid role ID")
		return
	}

	permissions, err := h.roleService.GetRolePermissions(uint(id))
	if err != nil {
		if err.Error() == "role not found" {
			utils.WriteError(w, http.StatusNotFound, "Role not found")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch role permissions")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Role permissions retrieved successfully",
		"data":    permissions,
	})
}

// @Router /roles/{id}/permissions [post]
func (h *RoleHandler) AssignPermissions(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid role ID")
		return
	}

	var input struct {
		PermissionIDs []uint `json:"permission_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(input.PermissionIDs) == 0 {
		utils.WriteError(w, http.StatusBadRequest, "Permission IDs are required")
		return
	}

	err = h.roleService.AssignPermissions(uint(id), input.PermissionIDs)
	if err != nil {
		if err.Error() == "role not found" {
			utils.WriteError(w, http.StatusNotFound, "Role not found")
			return
		}
		if err.Error() == "some permissions not found" {
			utils.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to assign permissions")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Permissions assigned successfully",
	})
}
