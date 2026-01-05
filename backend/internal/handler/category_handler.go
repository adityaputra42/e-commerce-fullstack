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

type CategoryHandler struct {
	service services.CategoryService
}

func NewCategoryHandler(svc services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		service: svc,
	}
}

// CreateCategory handles POST /api/v1/categories
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req models.Category

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	result, err := h.service.Create(req)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to create category", err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, "Category created successfully", result)
}

// UpdateCategory handles PUT /api/v1/categories/{id}
func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid category ID", err)
		return
	}

	var req models.Category
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	req.ID = id

	result, err := h.service.Update(req)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "category not found" {
			statusCode = http.StatusNotFound
		}
		utils.WriteError(w, statusCode, "Failed to update category", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Category updated successfully", result)
}

// DeleteCategory handles DELETE /api/v1/categories/{id}
func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid category ID", err)
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "category not found" {
			statusCode = http.StatusNotFound
		}
		utils.WriteError(w, statusCode, "Failed to delete category", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Category deleted successfully", nil)
}

// GetCategoryById handles GET /api/v1/categories/{id}
func (h *CategoryHandler) GetCategoryById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid category ID", err)
		return
	}

	result, err := h.service.GetById(id)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "category not found" {
			statusCode = http.StatusNotFound
		}
		utils.WriteError(w, statusCode, "Failed to get category", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Category retrieved successfully", result)
}

// GetAllCategories handles GET /api/v1/categories
func (h *CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	page := int64(1)
	limit := int64(10)
	search := ""
	sortBy := "id"

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.ParseInt(pageStr, 10, 64); err == nil {
			page = p
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 64); err == nil {
			limit = l
		}
	}

	if s := r.URL.Query().Get("search"); s != "" {
		search = s
	}

	if sort := r.URL.Query().Get("sort_by"); sort != "" {
		sortBy = sort
	}

	param := models.CategoryListRequest{
		Page:   int(page),
		Limit:  int(limit),
		Search: search,
		SortBy: sortBy,
	}

	result, err := h.service.GetAll(param)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get categories", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Categories retrieved successfully", result)
}
