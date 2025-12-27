package handler

import (
	"e-commerce/backend/internal/middleware"
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/services"
	"e-commerce/backend/internal/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type AddressHandler struct {
	addressService services.AddressService
}

func NewAddressHandler(addressService services.AddressService) *AddressHandler {
	return &AddressHandler{
		addressService: addressService,
	}
}

// CreateAddress handles POSaT /addresses
func (h *AddressHandler) CreateAddress(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)

	var param models.CreateAddress
	if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	result, err := h.addressService.CreateAddress(int64(userID), param)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, "Address created successfully", result)
}

// UpdateAddress handles PUT /addresses/:id
func (h *AddressHandler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)

	// Extract ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		utils.WriteError(w, http.StatusBadRequest, "Invalid address ID", errors.New("Invalid address ID"))
		return
	}

	addressID, err := strconv.ParseInt(pathParts[len(pathParts)-1], 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid address ID", err)
		return
	}

	var param models.UpdateAddress
	if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	result, err := h.addressService.UpdateAddress(addressID, int64(userID), param)
	if err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			utils.WriteError(w, http.StatusForbidden, err.Error(), err)
			return
		}
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, err.Error(), err)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Address updated successfully", result)
}

// GetAddresses handles GET /addresses
func (h *AddressHandler) GetAddresses(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)

	// Parse query parameters
	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	limit, _ := strconv.Atoi(query.Get("limit"))

	param := models.AddressListRequest{
		UserId: &userID,
		Page:   page,
		Limit:  limit,
	}

	result, err := h.addressService.FindAllAddress(param)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Addresses retrieved successfully", result)
}

func (h *AddressHandler) GetAddressByID(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)

	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		utils.WriteError(w, http.StatusBadRequest, "Invalid path url", fmt.Errorf("Invalid path url"))
		return
	}

	addressID, err := strconv.ParseInt(pathParts[len(pathParts)-1], 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid address ID", err)
		return
	}

	result, err := h.addressService.FindById(addressID, int64(userID))
	if err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			utils.WriteError(w, http.StatusForbidden, err.Error(), err)
			return
		}
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, err.Error(), err)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Address retrieved successfully", result)
}

// DeleteAddress handles DELETE /addresses/:id
func (h *AddressHandler) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)

	// Extract ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		utils.WriteError(w, http.StatusBadRequest, "Invalid path url", fmt.Errorf("Invalid path url"))
		return
	}

	addressID, err := strconv.ParseInt(pathParts[len(pathParts)-1], 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid address ID", err)
		return
	}

	// Verify ownership before deleting
	_, err = h.addressService.FindById(addressID, int64(userID))
	if err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			utils.WriteError(w, http.StatusForbidden, err.Error(), err)
			return
		}
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, err.Error(), err)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	if err := h.addressService.DeleteAddress(addressID); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Address deleted successfully", nil)
}
