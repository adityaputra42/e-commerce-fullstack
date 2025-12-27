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

type PaymentHandler struct {
	paymentService services.PaymentService
}

func NewPaymentHandler(paymentService services.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// @Router /payments [post]
func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var input models.CreatePayment
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate input
	if input.TransactionID == "" {
		utils.WriteError(w, http.StatusBadRequest, "Transaction ID is required", fmt.Errorf("Transaction ID is required"))
		return
	}
	if input.TotalPayment <= 0 {
		utils.WriteError(w, http.StatusBadRequest, "Total payment must be greater than 0", fmt.Errorf("Total payment must be greater than 0"))
		return
	}

	payment, err := h.paymentService.CreatePayment(input)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "transaction not found" {
			utils.WriteError(w, http.StatusNotFound, errMsg, err)

		}
		if errMsg == "total payment didn't match with transaction total price" {
			utils.WriteError(w, http.StatusBadRequest, errMsg, err)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create payment", err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, "Payment created successfully", payment)
}

// @Router /payments [get]
func (h *PaymentHandler) GetAllPayments(w http.ResponseWriter, r *http.Request) {

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

	param := models.PaymentListRequest{

		SortBy: sortBy,
		Page:   page,
		Limit:  limit,
	}

	payments, err := h.paymentService.FindAllPayment(param)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch payments", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Payments retrieved successfully", payments)
}

// @Router /payments/{id} [get]
func (h *PaymentHandler) GetPaymentByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid payment ID", err)
		return
	}

	payment, err := h.paymentService.FindById(id)
	if err != nil {
		if err.Error() == "payment not found" || err.Error() == "invalid payment id" {
			utils.WriteError(w, http.StatusNotFound, "Payment not found", err)
			return
		}
		if err.Error() == "transaction not found" {
			utils.WriteError(w, http.StatusNotFound, "Related transaction not found", err)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch payment", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Payment retrieved successfully", payment)
}

// @Router /payments/{id} [put]
func (h *PaymentHandler) UpdatePayment(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid payment ID", err)
		return
	}

	var input models.UpdatePayment
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Set ID from URL
	input.ID = id

	// Validate status
	validStatuses := []string{"pending", "confirmed", "rejected", "cancelled", "completed", "refunded"}
	isValid := false
	for _, status := range validStatuses {
		if input.Status == status {
			isValid = true
			break
		}
	}

	if !isValid {
		utils.WriteError(w, http.StatusBadRequest, "Invalid payment status", err)
		return
	}

	payment, err := h.paymentService.UpdatePayment(input)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "payment not found" || errMsg == "invalid payment id" {
			utils.WriteError(w, http.StatusNotFound, "Payment not found", err)
			return
		}
		if errMsg == "payment status is already the same, no changes needed" {
			utils.WriteError(w, http.StatusBadRequest, errMsg, err)
			return
		}
		if contains(errMsg, "invalid status transition") {
			utils.WriteError(w, http.StatusBadRequest, errMsg, err)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update payment", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Payment updated successfully", payment)
}

// @Router /payments/{id} [delete]
func (h *PaymentHandler) DeletePayment(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid payment ID", err)
		return
	}

	err = h.paymentService.DeletePayment(id)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "payment not found" || errMsg == "invalid payment id" {
			utils.WriteError(w, http.StatusNotFound, "Payment not found", err)
			return
		}
		if errMsg == "cannot delete payment with completed or confirmed status" {
			utils.WriteError(w, http.StatusBadRequest, errMsg, err)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete payment", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Payment deleted successfully", nil)
}
