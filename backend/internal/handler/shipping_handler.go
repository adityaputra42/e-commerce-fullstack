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

type ShippingHandler struct {
	shippingService services.ShippingService
}

func NewShippingHandler(shippingService services.ShippingService) *ShippingHandler {
	return &ShippingHandler{
		shippingService: shippingService,
	}
}

// GetAllShipping - GET /api/v1/shipping
// @Summary List all shipping methods
// @Description Get a paginated list of all shipping methods
// @Tags Shipping
// @Accept json
// @Produce json
// @Param sort_by query string false "Sort by field"
// @Param search query string false "Search query"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} utils.Response{data=[]models.Shipping} "Success"
// @Router /shipping [get]
// @Security Bearer
func (h *ShippingHandler) GetAllShipping(w http.ResponseWriter, r *http.Request) {

	sortBy := r.URL.Query().Get("sort_by")
	search := r.URL.Query().Get("search")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	param := models.ShippingListRequest{
		SortBy: sortBy,
		Search: search,
		Page:   page,
		Limit:  limit,
	}

	shippings, err := h.shippingService.FindAllShipping(param)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch shipping methods", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Shipping methods retrieved successfully", shippings)
}

// GetShippingByID - GET /api/v1/shipping/{id}
// @Summary Get shipping method by ID
// @Description Get detailed information about a shipping method
// @Tags Shipping
// @Accept json
// @Produce json
// @Param id path int true "Shipping ID"
// @Success 200 {object} utils.Response{data=models.Shipping} "Success"
// @Failure 404 {object} utils.Response "Shipping method not found"
// @Router /shipping/{id} [get]
// @Security Bearer
func (h *ShippingHandler) GetShippingByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid shipping ID", err)
		return
	}

	shipping, err := h.shippingService.FindByID(id)
	if err != nil {
		if err.Error() == "invalid shipping id" {
			utils.WriteError(w, http.StatusBadRequest, "Invalid shipping ID", err)
			return
		}
		utils.WriteError(w, http.StatusNotFound, "Shipping method not found", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Shipping method retrieved successfully", shipping)
}

// CreateShipping - POST /api/v1/shipping
// @Summary Create a new shipping method
// @Description Create a new shipping method for orders
// @Tags Shipping
// @Accept json
// @Produce json
// @Param request body models.CreateShipping true "Shipping request"
// @Success 201 {object} utils.Response{data=models.Shipping} "Shipping method created successfully"
// @Router /shipping [post]
// @Security Bearer
func (h *ShippingHandler) CreateShipping(w http.ResponseWriter, r *http.Request) {
	var input models.CreateShipping
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if input.State != "active" && input.State != "inactive" {
		utils.WriteError(w, http.StatusBadRequest, "State must be 'active' or 'inactive'", fmt.Errorf("State must be 'active' or 'inactive'"))
		return
	}

	shipping, err := h.shippingService.CreateShipping(input)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, "Shipping method created successfully", shipping)
}

// @Router /shipping/{id} [put]
func (h *ShippingHandler) UpdateShipping(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid shipping ID", err)
		return
	}

	var input models.UpdateShipping
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	input.ID = id

	if input.State != "active" && input.State != "inactive" {
		utils.WriteError(w, http.StatusBadRequest, "State must be 'active' or 'inactive'", err)
		return
	}

	shipping, err := h.shippingService.UpdateShipping(input)
	if err != nil {
		if err.Error() == "invalid shipping id" {
			utils.WriteError(w, http.StatusBadRequest, "Invalid shipping ID", err)
			return
		}
		utils.WriteError(w, http.StatusNotFound, "Shipping method not found", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Shipping method updated successfully", shipping)
}

// @Router /shipping/{id} [delete]
func (h *ShippingHandler) DeleteShipping(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid shipping ID", err)
		return
	}

	err = h.shippingService.DeleteShipping(id)
	if err != nil {
		if err.Error() == "invalid shipping id" {
			utils.WriteError(w, http.StatusBadRequest, "Invalid shipping ID", err)
			return
		}
		utils.WriteError(w, http.StatusNotFound, "Shipping method not found", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Shipping method deleted successfully", nil)
}
