package services

import (
	"context"
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/repository"
	"fmt"
	"time"
)

type DashboardService interface {
	GetDashboardStats(ctx context.Context) (*models.DashboardStatsResponse, error)
	GetRevenueStats(ctx context.Context) (*models.RevenueStatsResponse, error)
	GetOrderStats(ctx context.Context) (*models.OrderStatsResponse, error)
	GetRecentOrders(ctx context.Context, limit int) (*models.RecentOrdersResponse, error)
	GetTopProducts(ctx context.Context, limit int) (*models.TopProductsResponse, error)
	GetLowStockProducts(ctx context.Context, threshold, limit int) (*models.LowStockProductsResponse, error)
	GetOrderAnalytics(ctx context.Context, days int) (*models.OrderAnalyticsResponse, error)
	GetUserGrowth(ctx context.Context, days int) (*models.UserGrowthResponse, error)
	GetSystemHealth(ctx context.Context) (*models.SystemHealthResponse, error)
	GetRecentActivity(ctx context.Context, limit int) (*models.RecentActivityResponse, error)
}

type dashboardService struct {
	dashboardRepo repository.DashboardRepository
}

func NewDashboardService(dashboardRepo repository.DashboardRepository) DashboardService {
	return &dashboardService{
		dashboardRepo: dashboardRepo,
	}
}

// ==================== Response Structures ====================

// ==================== Service Implementation ====================

func (s *dashboardService) GetDashboardStats(ctx context.Context) (*models.DashboardStatsResponse, error) {
	totalUsers, err := s.dashboardRepo.CountTotalUsers(ctx)
	if err != nil {
		return nil, err
	}

	activeUsers, err := s.dashboardRepo.CountActiveUsers(ctx)
	if err != nil {
		return nil, err
	}

	totalProducts, err := s.dashboardRepo.CountTotalProducts(ctx)
	if err != nil {
		return nil, err
	}

	totalOrders, err := s.dashboardRepo.CountTotalOrders(ctx)
	if err != nil {
		return nil, err
	}

	pendingOrders, err := s.dashboardRepo.CountPendingOrders(ctx)
	if err != nil {
		return nil, err
	}

	totalCategories, err := s.dashboardRepo.CountTotalCategories(ctx)
	if err != nil {
		return nil, err
	}

	totalRoles, err := s.dashboardRepo.CountTotalRoles(ctx)
	if err != nil {
		return nil, err
	}

	today := time.Now().Truncate(24 * time.Hour)
	newUsersToday, err := s.dashboardRepo.GetNewUsersCount(ctx, today)
	if err != nil {
		return nil, err
	}

	weekAgo := time.Now().AddDate(0, 0, -7)
	newUsersThisWeek, err := s.dashboardRepo.GetNewUsersCount(ctx, weekAgo)
	if err != nil {
		return nil, err
	}

	return &models.DashboardStatsResponse{
		TotalUsers:       totalUsers,
		ActiveUsers:      activeUsers,
		TotalProducts:    totalProducts,
		TotalOrders:      totalOrders,
		PendingOrders:    pendingOrders,
		TotalCategories:  totalCategories,
		TotalRoles:       totalRoles,
		NewUsersToday:    newUsersToday,
		NewUsersThisWeek: newUsersThisWeek,
	}, nil
}

func (s *dashboardService) GetRevenueStats(ctx context.Context) (*models.RevenueStatsResponse, error) {
	totalRevenue, err := s.dashboardRepo.GetTotalRevenue(ctx)
	if err != nil {
		return nil, err
	}

	revenueToday, err := s.dashboardRepo.GetRevenueToday(ctx)
	if err != nil {
		return nil, err
	}

	revenueThisMonth, err := s.dashboardRepo.GetRevenueThisMonth(ctx)
	if err != nil {
		return nil, err
	}

	weekAgo := time.Now().AddDate(0, 0, -7)
	now := time.Now()
	revenueThisWeek, err := s.dashboardRepo.GetRevenueByPeriod(ctx, weekAgo, now)
	if err != nil {
		return nil, err
	}

	return &models.RevenueStatsResponse{
		TotalRevenue:     totalRevenue,
		RevenueToday:     revenueToday,
		RevenueThisMonth: revenueThisMonth,
		RevenueThisWeek:  revenueThisWeek,
	}, nil
}

func (s *dashboardService) GetOrderStats(ctx context.Context) (*models.OrderStatsResponse, error) {
	totalOrders, err := s.dashboardRepo.CountTotalOrders(ctx)
	if err != nil {
		return nil, err
	}

	pendingOrders, err := s.dashboardRepo.CountPendingOrders(ctx)
	if err != nil {
		return nil, err
	}

	ordersByStatus, err := s.dashboardRepo.GetOrdersByStatus(ctx)
	if err != nil {
		return nil, err
	}

	return &models.OrderStatsResponse{
		TotalOrders:    totalOrders,
		PendingOrders:  pendingOrders,
		OrdersByStatus: ordersByStatus,
	}, nil
}

func (s *dashboardService) GetRecentOrders(ctx context.Context, limit int) (*models.RecentOrdersResponse, error) {
	orders, err := s.dashboardRepo.GetRecentOrders(ctx, limit)
	if err != nil {
		return nil, err
	}

	orderSummaries := make([]models.OrderSummary, len(orders))
	for i, order := range orders {
		productName := "Unknown Product"
		if order.Product.Name != "" {
			productName = order.Product.Name
		}

		colorName := "Unknown Color"
		if order.ColorVarian.Name != "" {
			colorName = order.ColorVarian.Name
		}

		size := "Unknown Size"
		if order.SizeVarian.Size != "" {
			size = order.SizeVarian.Size
		}

		orderSummaries[i] = models.OrderSummary{
			ID:          order.ID,
			ProductName: productName,
			ColorName:   colorName,
			Size:        size,
			Quantity:    order.Quantity,
			Subtotal:    order.Subtotal,
			Status:      order.Status,
			CreatedAt:   order.CreatedAt,
		}
	}

	return &models.RecentOrdersResponse{
		Orders: orderSummaries,
	}, nil
}

func (s *dashboardService) GetTopProducts(ctx context.Context, limit int) (*models.TopProductsResponse, error) {
	products, err := s.dashboardRepo.GetTopSellingProducts(ctx, limit)
	if err != nil {
		return nil, err
	}

	return &models.TopProductsResponse{
		Products: products,
	}, nil
}

func (s *dashboardService) GetLowStockProducts(ctx context.Context, threshold, limit int) (*models.LowStockProductsResponse, error) {
	products, err := s.dashboardRepo.GetLowStockProducts(ctx, threshold, limit)
	if err != nil {
		return nil, err
	}

	return &models.LowStockProductsResponse{
		Products: products,
	}, nil
}

func (s *dashboardService) GetOrderAnalytics(ctx context.Context, days int) (*models.OrderAnalyticsResponse, error) {
	analytics, err := s.dashboardRepo.GetOrderAnalytics(ctx, days)
	if err != nil {
		return nil, err
	}

	var totalOrders int64
	var totalRevenue float64
	for _, a := range analytics {
		totalOrders += a.OrderCount
		totalRevenue += a.Revenue
	}

	averageOrder := float64(0)
	if totalOrders > 0 {
		averageOrder = totalRevenue / float64(totalOrders)
	}

	return &models.OrderAnalyticsResponse{
		Analytics: analytics,
		Summary: models.AnalyticsSummary{
			TotalOrders:  totalOrders,
			TotalRevenue: totalRevenue,
			AverageOrder: averageOrder,
		},
	}, nil
}

func (s *dashboardService) GetUserGrowth(ctx context.Context, days int) (*models.UserGrowthResponse, error) {
	growth, err := s.dashboardRepo.GetUserGrowthAnalytics(ctx, days)
	if err != nil {
		return nil, err
	}

	var totalNewUsers int64
	for _, g := range growth {
		totalNewUsers += g.UserCount
	}

	averagePerDay := float64(0)
	if days > 0 {
		averagePerDay = float64(totalNewUsers) / float64(days)
	}

	return &models.UserGrowthResponse{
		Growth: growth,
		Summary: models.UserGrowthSummary{
			TotalNewUsers: totalNewUsers,
			AveragePerDay: averagePerDay,
		},
	}, nil
}

func (s *dashboardService) GetSystemHealth(ctx context.Context) (*models.SystemHealthResponse, error) {
	health := &models.SystemHealthResponse{
		DatabaseStatus: "healthy",
		Timestamp:      time.Now(),
	}

	if err := s.dashboardRepo.CheckDatabaseHealth(ctx); err != nil {
		health.DatabaseStatus = "unhealthy"
	}

	activeUsers, err := s.dashboardRepo.CountActiveUsers(ctx)
	if err == nil {
		health.ActiveUsers = activeUsers
	}

	totalOrders, err := s.dashboardRepo.CountTotalOrders(ctx)
	if err == nil {
		health.TotalRequests = totalOrders
	}

	return health, nil
}

func (s *dashboardService) GetRecentActivity(ctx context.Context, limit int) (*models.RecentActivityResponse, error) {
	activities, err := s.dashboardRepo.GetRecentActivity(ctx, limit)
	if err != nil {
		return nil, err
	}

	activitySummaries := make([]models.ActivitySummary, len(activities))
	for i, activity := range activities {
		username := "Unknown"
		fullName := "Unknown User"

		if activity.User.Username != "" {
			username = activity.User.Username
		}

		if activity.User.FirstName != "" || activity.User.LastName != "" {
			fullName = fmt.Sprintf("%s %s", activity.User.FirstName, activity.User.LastName)
		}

		activitySummaries[i] = models.ActivitySummary{
			ID:        activity.ID,
			Username:  username,
			FullName:  fullName,
			Action:    activity.Action,
			Resource:  activity.Resource,
			Details:   activity.Details,
			IPAddress: activity.IPAddress,
			CreatedAt: activity.CreatedAt,
		}
	}

	return &models.RecentActivityResponse{
		Activities: activitySummaries,
	}, nil
}
