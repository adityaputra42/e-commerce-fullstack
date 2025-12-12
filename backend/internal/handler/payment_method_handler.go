package handler

import (
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/services"
	"e-commerce/backend/internal/utils"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type PaymentMethodHandler struct {
	paymentMethodService services.PaymentMethodService
}

func NewPaymentMethodHandler(paymentMethodService services.PaymentMethodService) *PaymentMethodHandler {
	return &PaymentMethodHandler{
		paymentMethodService: paymentMethodService,
	}
}

// @Router /payment-methods [get]
func (h *PaymentMethodHandler) GetAllPaymentMethods(w http.ResponseWriter, r *http.Request) {

	sortBy := r.URL.Query().Get("sort_by")
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

	param := models.PaymentMethodListRequest{
		SortBy: sortBy,
		Page:   page,
		Limit:  limit,
	}

	paymentMethods, err := h.paymentMethodService.FindAllPaymentMethod(param)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch payment methods")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Payment methods retrieved successfully",
		"data":    paymentMethods,
	})
}

// @Router /payment-methods/{id} [get]
func (h *PaymentMethodHandler) GetPaymentMethodByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid payment method ID")
		return
	}

	paymentMethod, err := h.paymentMethodService.FindById(id)
	if err != nil {
		if err.Error() == "payment method not found" || err.Error() == "invalid payment method id" {
			utils.WriteError(w, http.StatusNotFound, "Payment method not found")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch payment method")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Payment method retrieved successfully",
		"data":    paymentMethod,
	})
}

// @Router /payment-methods [post]
func (h *PaymentMethodHandler) CreatePaymentMethod(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form (max 10MB)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to parse form data")
		return
	}

	accountName := r.FormValue("account_name")
	accountNumber := r.FormValue("account_number")
	bankName := r.FormValue("bank_name")

	if accountName == "" || accountNumber == "" || bankName == "" {
		utils.WriteError(w, http.StatusBadRequest, "All fields are required")
		return
	}

	file, fileHeader, err := r.FormFile("bank_image")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Bank image is required")
		return
	}
	defer file.Close()

	param := models.CreatePaymentMethod{
		AccountName:   accountName,
		AccountNumber: accountNumber,
		BankName:      bankName,
	}

	paymentMethod, err := h.paymentMethodService.CreatePaymentMethod(param, fileHeader)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Payment method created successfully",
		"data":    paymentMethod,
	})
}

// @Router /payment-methods/{id} [put]
func (h *PaymentMethodHandler) UpdatePaymentMethod(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid payment method ID")
		return
	}

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to parse form data")
		return
	}

	accountName := r.FormValue("account_name")
	accountNumber := r.FormValue("account_number")
	bankName := r.FormValue("bank_name")

	if accountName == "" || accountNumber == "" || bankName == "" {
		utils.WriteError(w, http.StatusBadRequest, "All fields are required")
		return
	}

	param := models.UpdatePaymentMethod{
		ID:            id,
		AccountName:   accountName,
		AccountNumber: accountNumber,
		BankName:      bankName,
	}

	var fileHeader *multipart.FileHeader
	file, header, err := r.FormFile("bank_image")
	if err == nil {
		defer file.Close()
		fileHeader = header
	}

	paymentMethod, err := h.paymentMethodService.UpdatePaymentMethod(param, fileHeader)
	if err != nil {
		if err.Error() == "payment method not found" {
			utils.WriteError(w, http.StatusNotFound, "Payment method not found")
			return
		}
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Payment method updated successfully",
		"data":    paymentMethod,
	})
}

// @Router /payment-methods/{id} [delete]
func (h *PaymentMethodHandler) DeletePaymentMethod(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid payment method ID")
		return
	}

	err = h.paymentMethodService.DeletePaymentMethod(id)
	if err != nil {
		if err.Error() == "payment method not found" || err.Error() == "invalid payment method id" {
			utils.WriteError(w, http.StatusNotFound, "Payment method not found")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete payment method")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Payment method deleted successfully",
	})
}
