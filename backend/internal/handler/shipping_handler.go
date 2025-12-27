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

// @Router /shipping [get]
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

// @Router /shipping/{id} [get]
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

// @Router /shipping [post]
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
