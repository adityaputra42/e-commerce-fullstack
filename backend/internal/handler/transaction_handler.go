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

type TransactionHandler struct {
	transactionService services.TransactionService
}

func NewTransactionHandler(transactionService services.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

// CreateTransaction - POST /api/v1/transactions
// @Summary Create a new transaction
// @Description Create a new transaction with product orders
// @Tags Transaction
// @Accept json
// @Produce json
// @Param request body models.CreateTransaction true "Transaction request"
// @Success 201 {object} utils.Response{data=models.Transaction} "Transaction created successfully"
// @Router /transactions [post]
// @Security Bearer
func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var input models.CreateTransaction
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate input
	if input.AddressID <= 0 {
		utils.WriteError(w, http.StatusBadRequest, "Address ID is required", fmt.Errorf("Address ID is required"))
		return
	}
	if input.ShippingID <= 0 {
		utils.WriteError(w, http.StatusBadRequest, "Shipping ID is required", fmt.Errorf("Shipping ID is required"))
		return
	}
	if input.PaymentMethodID <= 0 {
		utils.WriteError(w, http.StatusBadRequest, "Payment method ID is required", fmt.Errorf("Payment method ID is required"))
		return
	}
	if len(input.ProductOrders) == 0 {
		utils.WriteError(w, http.StatusBadRequest, "At least one product order is required", fmt.Errorf("At least one product order is required"))
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
			utils.WriteError(w, http.StatusNotFound, errMsg, err)
			return
		}
		if contains(errMsg, "insufficient stock") {
			utils.WriteError(w, http.StatusBadRequest, errMsg, err)
			return
		}
		if contains(errMsg, "context canceled") {
			utils.WriteError(w, http.StatusRequestTimeout, "Request timeout", err)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create transaction", err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, "Transaction created successfully", transaction)
}

// GetAllTransactions - GET /api/v1/transactions
// @Summary List all transactions
// @Description Get a paginated list of all transactions
// @Tags Transaction
// @Accept json
// @Produce json
// @Param sort_by query string false "Sort by field"
// @Param search query string false "Search query"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} utils.Response{data=[]models.Transaction} "Success"
// @Router /transactions [get]
// @Security Bearer
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
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch transactions", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Transactions retrieved successfully", transactions)
}

// GetTransactionByID - GET /api/v1/transactions/{tx_id}
// @Summary Get transaction by ID
// @Description Get detailed information about a transaction
// @Tags Transaction
// @Accept json
// @Produce json
// @Param tx_id path string true "Transaction ID"
// @Success 200 {object} utils.Response{data=models.Transaction} "Success"
// @Failure 404 {object} utils.Response "Transaction not found"
// @Router /transactions/{tx_id} [get]
// @Security Bearer
func (h *TransactionHandler) GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	txID := chi.URLParam(r, "tx_id")
	if txID == "" {
		utils.WriteError(w, http.StatusBadRequest, "Transaction ID is required", fmt.Errorf("Transaction ID is required"))
		return
	}

	transaction, err := h.transactionService.FindTransactionById(txID)
	if err != nil {
		if err.Error() == "transaction not found" {
			utils.WriteError(w, http.StatusNotFound, "Transaction not found", err)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch transaction", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Transaction retrieved successfully", transaction)
}

// @Router /transactions/{tx_id} [put]
func (h *TransactionHandler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	txID := chi.URLParam(r, "tx_id")
	if txID == "" {
		utils.WriteError(w, http.StatusBadRequest, "Transaction ID is required", fmt.Errorf("Transaction ID is required"))
		return
	}

	var input models.UpdateTransaction
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
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
		utils.WriteError(w, http.StatusBadRequest, "Invalid transaction status", fmt.Errorf("Invalid transaction status"))
		return
	}

	transaction, err := h.transactionService.UpdateTransaction(input)
	if err != nil {
		if err.Error() == "transaction not found" {
			utils.WriteError(w, http.StatusNotFound, "Transaction not found", err)
			return
		}
		if err.Error() == "invalid status transition" {
			utils.WriteError(w, http.StatusBadRequest, "Invalid status transition", err)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update transaction", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Transaction updated successfully", transaction)
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
