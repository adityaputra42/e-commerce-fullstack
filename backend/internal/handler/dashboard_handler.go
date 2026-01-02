package handler

import (
	"e-commerce/backend/internal/services"
	"e-commerce/backend/internal/utils"
	"log"
	"net/http"
	"strconv"
)

type DashboardHandler struct {
	dashboardService services.DashboardService
}

func NewDashboardHandler(dashboardService services.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

// GetDashboardStats - GET /api/v1/dashboard/stats
// @Summary Dashboard overview stats
// @Description Get overall statistics for the dashboard
// @Tags Dashboard
// @Produce json
// @Success 200 {object} utils.Response "Dashboard stats retrieved successfully"
// @Router /dashboard/stats [get]
// @Security Bearer
func (h *DashboardHandler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("Authorization Header:", r.Header.Get("Authorization"))

	stats, err := h.dashboardService.GetDashboardStats(ctx)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get dashboard stats", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Dashboard stats retrieved successfully", stats)
}

// GetRevenueStats - GET /api/v1/dashboard/revenue
// @Summary Revenue statistics
// @Description Get revenue-related statistics
// @Tags Dashboard
// @Produce json
// @Success 200 {object} utils.Response "Revenue stats retrieved successfully"
// @Router /dashboard/revenue [get]
// @Security Bearer
func (h *DashboardHandler) GetRevenueStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("Authorization Header:", r.Header.Get("Authorization"))

	stats, err := h.dashboardService.GetRevenueStats(ctx)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get revenue stats", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Revenue stats retrieved successfully", stats)
}

// GetOrderStats - GET /api/v1/dashboard/orders
// @Summary Order statistics
// @Description Get order-related statistics
// @Tags Dashboard
// @Produce json
// @Success 200 {object} utils.Response "Order stats retrieved successfully"
// @Router /dashboard/orders [get]
// @Security Bearer
func (h *DashboardHandler) GetOrderStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	stats, err := h.dashboardService.GetOrderStats(ctx)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get order stats", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Order stats retrieved successfully", stats)
}

// GetRecentOrders - GET /api/v1/dashboard/recent-orders
// @Summary Recent orders
// @Description Get a list of the most recent orders
// @Tags Dashboard
// @Produce json
// @Param limit query int false "Limit number of orders" default(10)
// @Success 200 {object} utils.Response "Recent orders retrieved successfully"
// @Router /dashboard/recent-orders [get]
// @Security Bearer
func (h *DashboardHandler) GetRecentOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse limit from query parameter
	limitStr := r.URL.Query().Get("limit")
	limit := 10 // default
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit <= 0 {
			utils.WriteError(w, http.StatusBadRequest, "Invalid limit parameter", err)
			return
		}
		limit = parsedLimit
	}

	orders, err := h.dashboardService.GetRecentOrders(ctx, limit)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get recent orders", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Recent orders retrieved successfully", orders)
}

// GetTopProducts retrieves top selling products
func (h *DashboardHandler) GetTopProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse limit from query parameter
	limitStr := r.URL.Query().Get("limit")
	limit := 10 // default
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit <= 0 {
			utils.WriteError(w, http.StatusBadRequest, "Invalid limit parameter", err)
			return
		}
		limit = parsedLimit
	}

	products, err := h.dashboardService.GetTopProducts(ctx, limit)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get top products", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Top products retrieved successfully", products)
}

// GetLowStockProducts retrieves products with low stock
func (h *DashboardHandler) GetLowStockProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse threshold from query parameter
	thresholdStr := r.URL.Query().Get("threshold")
	threshold := 10 // default
	if thresholdStr != "" {
		parsedThreshold, err := strconv.Atoi(thresholdStr)
		if err != nil || parsedThreshold < 0 {
			utils.WriteError(w, http.StatusBadRequest, "Invalid threshold parameter", err)
			return
		}
		threshold = parsedThreshold
	}

	// Parse limit from query parameter
	limitStr := r.URL.Query().Get("limit")
	limit := 10 // default
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit <= 0 {
			utils.WriteError(w, http.StatusBadRequest, "Invalid limit parameter", err)
			return
		}
		limit = parsedLimit
	}

	products, err := h.dashboardService.GetLowStockProducts(ctx, threshold, limit)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get low stock products", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Low stock products retrieved successfully", products)
}

// GetOrderAnalytics retrieves order analytics
func (h *DashboardHandler) GetOrderAnalytics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse days from query parameter
	daysStr := r.URL.Query().Get("days")
	days := 30 // default
	if daysStr != "" {
		parsedDays, err := strconv.Atoi(daysStr)
		if err != nil || parsedDays <= 0 {
			utils.WriteError(w, http.StatusBadRequest, "Invalid days parameter", err)
			return
		}
		days = parsedDays
	}

	analytics, err := h.dashboardService.GetOrderAnalytics(ctx, days)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get order analytics", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Order analytics retrieved successfully", analytics)
}

// GetUserGrowth retrieves user growth analytics
func (h *DashboardHandler) GetUserGrowth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse days from query parameter
	daysStr := r.URL.Query().Get("days")
	days := 30 // default
	if daysStr != "" {
		parsedDays, err := strconv.Atoi(daysStr)
		if err != nil || parsedDays <= 0 {
			utils.WriteError(w, http.StatusBadRequest, "Invalid days parameter", err)
			return
		}
		days = parsedDays
	}

	growth, err := h.dashboardService.GetUserGrowth(ctx, days)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get user growth", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "User growth retrieved successfully", growth)
}

// GetSystemHealth retrieves system health status
func (h *DashboardHandler) GetSystemHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	health, err := h.dashboardService.GetSystemHealth(ctx)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get system health", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "System health retrieved successfully", health)
}

// GetRecentActivity retrieves recent activity logs
func (h *DashboardHandler) GetRecentActivity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse limit from query parameter
	limitStr := r.URL.Query().Get("limit")
	limit := 20 // default
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit <= 0 {
			utils.WriteError(w, http.StatusBadRequest, "Invalid limit parameter", err)
			return
		}
		limit = parsedLimit
	}

	activity, err := h.dashboardService.GetRecentActivity(ctx, limit)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get recent activity", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Recent activity retrieved successfully", activity)
}
