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
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch shipping methods")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Shipping methods retrieved successfully",
		"data":    shippings,
	})
}

// @Router /shipping/{id} [get]
func (h *ShippingHandler) GetShippingByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid shipping ID")
		return
	}

	shipping, err := h.shippingService.FindByID(id)
	if err != nil {
		if err.Error() == "invalid shipping id" {
			utils.WriteError(w, http.StatusBadRequest, "Invalid shipping ID")
			return
		}
		utils.WriteError(w, http.StatusNotFound, "Shipping method not found")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Shipping method retrieved successfully",
		"data":    shipping,
	})
}

// @Router /shipping [post]
func (h *ShippingHandler) CreateShipping(w http.ResponseWriter, r *http.Request) {
	var input models.CreateShipping
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if input.State != "active" && input.State != "inactive" {
		utils.WriteError(w, http.StatusBadRequest, "State must be 'active' or 'inactive'")
		return
	}

	shipping, err := h.shippingService.CreateShipping(input)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Shipping method created successfully",
		"data":    shipping,
	})
}

// @Router /shipping/{id} [put]
func (h *ShippingHandler) UpdateShipping(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid shipping ID")
		return
	}

	var input models.UpdateShipping
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	input.ID = id

	if input.State != "active" && input.State != "inactive" {
		utils.WriteError(w, http.StatusBadRequest, "State must be 'active' or 'inactive'")
		return
	}

	shipping, err := h.shippingService.UpdateShipping(input)
	if err != nil {
		if err.Error() == "invalid shipping id" {
			utils.WriteError(w, http.StatusBadRequest, "Invalid shipping ID")
			return
		}
		utils.WriteError(w, http.StatusNotFound, "Shipping method not found")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Shipping method updated successfully",
		"data":    shipping,
	})
}

// @Router /shipping/{id} [delete]
func (h *ShippingHandler) DeleteShipping(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid shipping ID")
		return
	}

	err = h.shippingService.DeleteShipping(id)
	if err != nil {
		if err.Error() == "invalid shipping id" {
			utils.WriteError(w, http.StatusBadRequest, "Invalid shipping ID")
			return
		}
		utils.WriteError(w, http.StatusNotFound, "Shipping method not found")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Shipping method deleted successfully",
	})
}
