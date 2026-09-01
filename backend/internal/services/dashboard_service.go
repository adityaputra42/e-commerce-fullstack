package services

import (
	"context"
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/repository"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
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
	return &dashboardService{dashboardRepo: dashboardRepo}
}

// ==================== Service Implementation ====================
//
// GetDashboardStats, GetRevenueStats and GetOrderStats each used to run
// 3-9 independent COUNT/SUM queries back to back, one HTTP request holding
// a DB connection for the sum of every query's latency instead of the
// slowest one. None of these queries read each other's output, so there is
// no reason for them to be sequential. Below they run concurrently via
// errgroup and share one context, so a client disconnect or timeout on ctx
// cancels every in-flight query instead of leaving them to finish uselessly.
//
// Behavior preserved exactly: same return type, same "first error wins" —
// errgroup.Wait() returns the first non-nil error, matching the original
// early-return-on-first-error control flow.

func (s *dashboardService) GetDashboardStats(ctx context.Context) (*models.DashboardStatsResponse, error) {
	var (
		totalUsers, activeUsers, totalProducts      int64
		totalOrders, pendingOrders, totalCategories int64
		totalRoles, newUsersToday, newUsersThisWeek int64
	)

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() (err error) { totalUsers, err = s.dashboardRepo.CountTotalUsers(gctx); return })
	g.Go(func() (err error) { activeUsers, err = s.dashboardRepo.CountActiveUsers(gctx); return })
	g.Go(func() (err error) { totalProducts, err = s.dashboardRepo.CountTotalProducts(gctx); return })
	g.Go(func() (err error) { totalOrders, err = s.dashboardRepo.CountTotalOrders(gctx); return })
	g.Go(func() (err error) { pendingOrders, err = s.dashboardRepo.CountPendingOrders(gctx); return })
	g.Go(func() (err error) { totalCategories, err = s.dashboardRepo.CountTotalCategories(gctx); return })
	g.Go(func() (err error) { totalRoles, err = s.dashboardRepo.CountTotalRoles(gctx); return })
	g.Go(func() (err error) {
		newUsersToday, err = s.dashboardRepo.GetNewUsersCount(gctx, time.Now().Truncate(24*time.Hour))
		return
	})
	g.Go(func() (err error) {
		newUsersThisWeek, err = s.dashboardRepo.GetNewUsersCount(gctx, time.Now().AddDate(0, 0, -7))
		return
	})

	if err := g.Wait(); err != nil {
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
	var totalRevenue, revenueToday, revenueThisMonth, revenueThisWeek float64

	g, gctx := errgroup.WithContext(ctx)
	now := time.Now()

	g.Go(func() (err error) { totalRevenue, err = s.dashboardRepo.GetTotalRevenue(gctx); return })
	g.Go(func() (err error) { revenueToday, err = s.dashboardRepo.GetRevenueToday(gctx); return })
	g.Go(func() (err error) { revenueThisMonth, err = s.dashboardRepo.GetRevenueThisMonth(gctx); return })
	g.Go(func() (err error) {
		revenueThisWeek, err = s.dashboardRepo.GetRevenueByPeriod(gctx, now.AddDate(0, 0, -7), now)
		return
	})

	if err := g.Wait(); err != nil {
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
	var totalOrders, pendingOrders int64
	var ordersByStatus []models.OrderStatusCount

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() (err error) { totalOrders, err = s.dashboardRepo.CountTotalOrders(gctx); return })
	g.Go(func() (err error) { pendingOrders, err = s.dashboardRepo.CountPendingOrders(gctx); return })
	g.Go(func() (err error) { ordersByStatus, err = s.dashboardRepo.GetOrdersByStatus(gctx); return })

	if err := g.Wait(); err != nil {
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

	return &models.RecentOrdersResponse{Orders: orderSummaries}, nil
}

func (s *dashboardService) GetTopProducts(ctx context.Context, limit int) (*models.TopProductsResponse, error) {
	products, err := s.dashboardRepo.GetTopSellingProducts(ctx, limit)
	if err != nil {
		return nil, err
	}
	return &models.TopProductsResponse{Products: products}, nil
}

func (s *dashboardService) GetLowStockProducts(ctx context.Context, threshold, limit int) (*models.LowStockProductsResponse, error) {
	products, err := s.dashboardRepo.GetLowStockProducts(ctx, threshold, limit)
	if err != nil {
		return nil, err
	}
	return &models.LowStockProductsResponse{Products: products}, nil
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

// GetSystemHealth implements DashboardService.
//
// NOTE ON BEHAVIOR: the original silently swallows errors from
// CountActiveUsers/CountTotalOrders (an `if err == nil` guard, no else
// branch) so a failing query just leaves that field at its zero value
// instead of failing the whole health check. That swallow is left intact on
// purpose — this is a status endpoint, and making it fail loudly would be a
// functional change, not a quality one. Flagged in the report as worth a
// real decision from you, not fixed silently here.
func (s *dashboardService) GetSystemHealth(ctx context.Context) (*models.SystemHealthResponse, error) {
	health := &models.SystemHealthResponse{
		DatabaseStatus: "healthy",
		Timestamp:      time.Now(),
	}

	if err := s.dashboardRepo.CheckDatabaseHealth(ctx); err != nil {
		health.DatabaseStatus = "unhealthy"
	}

	if activeUsers, err := s.dashboardRepo.CountActiveUsers(ctx); err == nil {
		health.ActiveUsers = activeUsers
	}

	if totalOrders, err := s.dashboardRepo.CountTotalOrders(ctx); err == nil {
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

	return &models.RecentActivityResponse{Activities: activitySummaries}, nil
}
