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

type TransactionHandler struct {
	transactionService services.TransactionService
}

func NewTransactionHandler(transactionService services.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

// @Router /transactions [post]
func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var input models.CreateTransaction
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if input.AddressID <= 0 {
		utils.WriteError(w, http.StatusBadRequest, "Address ID is required")
		return
	}
	if input.ShippingID <= 0 {
		utils.WriteError(w, http.StatusBadRequest, "Shipping ID is required")
		return
	}
	if input.PaymentMethodID <= 0 {
		utils.WriteError(w, http.StatusBadRequest, "Payment method ID is required")
		return
	}
	if len(input.ProductOrders) == 0 {
		utils.WriteError(w, http.StatusBadRequest, "At least one product order is required")
		return
	}

	// Get context from request
	ctx := r.Context()

	// Create transaction
	transaction, err := h.transactionService.CreateTransaction(ctx, input)
	if err != nil {
		// Handle specific errors
		errMsg := err.Error()
		if contains(errMsg, "not found") {
			utils.WriteError(w, http.StatusNotFound, errMsg)
			return
		}
		if contains(errMsg, "insufficient stock") {
			utils.WriteError(w, http.StatusBadRequest, errMsg)
			return
		}
		if contains(errMsg, "context canceled") {
			utils.WriteError(w, http.StatusRequestTimeout, "Request timeout")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create transaction")
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Transaction created successfully",
		"data":    transaction,
	})
}

// @Router /transactions [get]
func (h *TransactionHandler) GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	sortBy := r.URL.Query().Get("sort_by")
	search := r.URL.Query().Get("search")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	// Default pagination
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

	param := models.TransactionListRequest{
		SortBy: sortBy,
		Search: search,
		Page:   page,
		Limit:  limit,
	}

	transactions, err := h.transactionService.FindAllTransaction(param)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch transactions")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Transactions retrieved successfully",
		"data":    transactions,
	})
}

// @Router /transactions/{tx_id} [get]
func (h *TransactionHandler) GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	txID := chi.URLParam(r, "tx_id")
	if txID == "" {
		utils.WriteError(w, http.StatusBadRequest, "Transaction ID is required")
		return
	}

	transaction, err := h.transactionService.FindTransactionById(txID)
	if err != nil {
		if err.Error() == "transaction not found" {
			utils.WriteError(w, http.StatusNotFound, "Transaction not found")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch transaction")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Transaction retrieved successfully",
		"data":    transaction,
	})
}

// @Router /transactions/{tx_id} [put]
func (h *TransactionHandler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	txID := chi.URLParam(r, "tx_id")
	if txID == "" {
		utils.WriteError(w, http.StatusBadRequest, "Transaction ID is required")
		return
	}

	var input models.UpdateTransaction
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Set ID from URL
	input.TxID = txID

	// Validate status
	validStatuses := []string{"waiting_payment", "paid", "processing", "shipped", "completed", "cancelled", "refunded"}
	isValid := false
	for _, status := range validStatuses {
		if input.Status == status {
			isValid = true
			break
		}
	}

	if !isValid {
		utils.WriteError(w, http.StatusBadRequest, "Invalid transaction status")
		return
	}

	transaction, err := h.transactionService.UpdateTransaction(input)
	if err != nil {
		if err.Error() == "transaction not found" {
			utils.WriteError(w, http.StatusNotFound, "Transaction not found")
			return
		}
		if err.Error() == "invalid status transition" {
			utils.WriteError(w, http.StatusBadRequest, "Invalid status transition")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update transaction")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Transaction updated successfully",
		"data":    transaction,
	})
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr ||
		len(s) > len(substr) && indexOf(s, substr) >= 0
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
