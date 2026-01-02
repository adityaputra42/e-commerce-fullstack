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

type OrderHandler struct {
	orderService services.OrderService
}

func NewOrderHandler(orderService services.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// GetAllOrders - GET /api/v1/orders
// @Summary List all orders
// @Description Get a paginated list of all orders
// @Tags Order
// @Accept json
// @Produce json
// @Param user_id query int false "Filter by User ID"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} utils.Response{data=[]models.Order} "Success"
// @Router /orders [get]
// @Security Bearer
func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {

	userIDStr := r.URL.Query().Get("user_id")
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

	var userID int64
	if userIDStr != "" {
		if uid, err := strconv.ParseInt(userIDStr, 10, 64); err == nil {
			userID = uid
		}
	}

	param := models.OrderListRequest{

		UserId: userID,
		Page:   int64(page),
		Limit:  int64(limit),
	}

	orders, err := h.orderService.FindAllOrder(param)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch orders", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "success", orders)
}

// GetOrderByID - GET /api/v1/orders/{id}
// @Summary Get order by ID
// @Description Get detailed information about an order
// @Tags Order
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} utils.Response{data=models.Order} "Success"
// @Failure 404 {object} utils.Response "Order not found"
// @Router /orders/{id} [get]
// @Security Bearer
func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "id")
	if orderID == "" {
		utils.WriteError(w, http.StatusBadRequest, "Order ID is required", fmt.Errorf("Order ID is required"))
		return
	}

	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "User not authenticated", fmt.Errorf("User not authenticated"))
		return
	}

	order, err := h.orderService.FindById(orderID, userID)
	if err != nil {
		if err.Error() == "order not found" {
			utils.WriteError(w, http.StatusNotFound, "Order not found", err)
			return
		}
		if err.Error() == "unauthorized: order does not belong to user" {
			utils.WriteError(w, http.StatusForbidden, "You don't have permission to access this order", err)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch order", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Order retrieved successfully", order)
}

// UpdateOrder - PUT /api/v1/orders/{id}
// @Summary Update order status
// @Description Update the status of an existing order
// @Tags Order
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param request body models.UpdateOrder true "Update order request"
// @Success 200 {object} utils.Response{data=models.Order} "Order updated successfully"
// @Router /orders/{id} [put]
// @Security Bearer
func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "id")
	if orderID == "" {
		utils.WriteError(w, http.StatusBadRequest, "Order ID is required", fmt.Errorf("Order ID is required"))
		return
	}

	var updateData models.UpdateOrder
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	updateData.ID = orderID

	validStatuses := []string{"pending", "confirmed", "processing", "shipped", "delivered", "cancelled"}
	isValidStatus := false
	for _, status := range validStatuses {
		if updateData.Status == status {
			isValidStatus = true
			break
		}
	}

	if !isValidStatus {
		utils.WriteError(w, http.StatusBadRequest, "Invalid order status", fmt.Errorf("Invalid order status"))
		return
	}

	updatedOrder, err := h.orderService.UpdateOrder(updateData)
	if err != nil {
		if err.Error() == "order not found" {
			utils.WriteError(w, http.StatusNotFound, "Order not found", err)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update order", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Order updated successfully", updatedOrder)
}

// @Router /orders/{id}/cancel [patch]
func (h *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "id")
	if orderID == "" {
		utils.WriteError(w, http.StatusBadRequest, "Order ID is required", fmt.Errorf("Order ID is required"))
		return
	}

	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "User not authenticated", fmt.Errorf("User not authenticated"))
		return
	}

	cancelledOrder, err := h.orderService.CancelOrder(orderID, userID)
	if err != nil {
		if err.Error() == "order not found" {
			utils.WriteError(w, http.StatusNotFound, "Order not found", err)
			return
		}
		if err.Error() == "unauthorized: order does not belong to user" {
			utils.WriteError(w, http.StatusForbidden, "You don't have permission to cancel this order", err)
			return
		}
		if err.Error() == "cannot cancel order with current status" {
			utils.WriteError(w, http.StatusBadRequest, "Cannot cancel order with current status", err)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to cancel order", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Order cancelled successfully", cancelledOrder)
}

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "id")
	if orderID == "" {
		utils.WriteError(w, http.StatusBadRequest, "Order ID is required", fmt.Errorf("Order ID is required"))
		return
	}

	err := h.orderService.DeleteOrder(orderID)
	if err != nil {
		if err.Error() == "order not found" {
			utils.WriteError(w, http.StatusNotFound, "Order not found", err)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete order", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Order deleted successfully", nil)
}
