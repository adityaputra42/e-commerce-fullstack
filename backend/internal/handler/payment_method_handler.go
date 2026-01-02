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

// GetAllPaymentMethods - GET /api/v1/payment-methods
// @Summary List all payment methods
// @Description Get a paginated list of all payment methods
// @Tags Payment Method
// @Accept json
// @Produce json
// @Param sort_by query string false "Sort by field"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} utils.Response{data=[]models.PaymentMethod} "Success"
// @Router /payment-methods [get]
// @Security Bearer
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
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch payment methods", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Payment methods retrieved successfully", paymentMethods)
}

// GetPaymentMethodByID - GET /api/v1/payment-methods/{id}
// @Summary Get payment method by ID
// @Description Get detailed information about a payment method
// @Tags Payment Method
// @Accept json
// @Produce json
// @Param id path int true "Payment Method ID"
// @Success 200 {object} utils.Response{data=models.PaymentMethod} "Success"
// @Failure 404 {object} utils.Response "Payment method not found"
// @Router /payment-methods/{id} [get]
// @Security Bearer
func (h *PaymentMethodHandler) GetPaymentMethodByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid payment method ID", err)
		return
	}

	paymentMethod, err := h.paymentMethodService.FindById(id)
	if err != nil {
		if err.Error() == "payment method not found" || err.Error() == "invalid payment method id" {
			utils.WriteError(w, http.StatusNotFound, "Payment method not found", err)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch payment method", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Payment method retrieved successfully", paymentMethod)
}

// CreatePaymentMethod - POST /api/v1/payment-methods
// @Summary Create a new payment method
// @Description Create a new payment method with bank details and image
// @Tags Payment Method
// @Accept multipart/form-data
// @Produce json
// @Param account_name formData string true "Account holder name"
// @Param account_number formData string true "Account number"
// @Param bank_name formData string true "Bank name"
// @Param bank_image formData file true "Bank logo/image"
// @Success 201 {object} utils.Response{data=models.PaymentMethod} "Payment method created successfully"
// @Router /payment-methods [post]
// @Security Bearer
func (h *PaymentMethodHandler) CreatePaymentMethod(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form (max 10MB)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to parse form data", err)
		return
	}

	accountName := r.FormValue("account_name")
	accountNumber := r.FormValue("account_number")
	bankName := r.FormValue("bank_name")

	if accountName == "" || accountNumber == "" || bankName == "" {
		utils.WriteError(w, http.StatusBadRequest, "All fields are required", err)
		return
	}

	file, fileHeader, err := r.FormFile("bank_image")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Bank image is required", err)
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
		utils.WriteError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, "Payment method created successfully", paymentMethod)
}

// @Router /payment-methods/{id} [put]
func (h *PaymentMethodHandler) UpdatePaymentMethod(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid payment method ID", err)
		return
	}

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to parse form data", err)
		return
	}

	accountName := r.FormValue("account_name")
	accountNumber := r.FormValue("account_number")
	bankName := r.FormValue("bank_name")

	var isActive *bool
	isActiveStr := r.FormValue("is_active")
	if isActiveStr != "" {
		b, err := strconv.ParseBool(isActiveStr)
		if err == nil {
			isActive = &b
		}
	}

	param := models.UpdatePaymentMethod{
		ID:            id,
		AccountName:   accountName,
		AccountNumber: accountNumber,
		BankName:      bankName,
		IsActive:      isActive,
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
			utils.WriteError(w, http.StatusNotFound, "Payment method not found", err)
			return
		}
		utils.WriteError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Payment method updated successfully", paymentMethod)
}

// @Router /payment-methods/{id} [delete]
func (h *PaymentMethodHandler) DeletePaymentMethod(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid payment method ID", err)
		return
	}

	err = h.paymentMethodService.DeletePaymentMethod(id)
	if err != nil {
		if err.Error() == "payment method not found" || err.Error() == "invalid payment method id" {
			utils.WriteError(w, http.StatusNotFound, "Payment method not found", err)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete payment method", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Payment method deleted successfully", nil)
}
