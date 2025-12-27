package repository

import (
	"context"
	"e-commerce/backend/internal/models"
	"time"

	"gorm.io/gorm"
)

type DashboardRepository interface {
	// Stats
	CountTotalUsers(ctx context.Context) (int64, error)
	CountActiveUsers(ctx context.Context) (int64, error)
	CountTotalProducts(ctx context.Context) (int64, error)
	CountTotalOrders(ctx context.Context) (int64, error)
	CountPendingOrders(ctx context.Context) (int64, error)
	CountTotalCategories(ctx context.Context) (int64, error)
	CountTotalRoles(ctx context.Context) (int64, error)

	// Revenue
	GetTotalRevenue(ctx context.Context) (float64, error)
	GetRevenueByPeriod(ctx context.Context, startDate, endDate time.Time) (float64, error)
	GetRevenueToday(ctx context.Context) (float64, error)
	GetRevenueThisMonth(ctx context.Context) (float64, error)

	// Orders
	GetRecentOrders(ctx context.Context, limit int) ([]models.Order, error)
	GetOrdersByStatus(ctx context.Context) ([]models.OrderStatusCount, error)
	GetOrderAnalytics(ctx context.Context, days int) ([]models.OrderAnalytics, error)

	// Products
	GetTopSellingProducts(ctx context.Context, limit int) ([]models.TopProduct, error)
	GetLowStockProducts(ctx context.Context, threshold int, limit int) ([]models.LowStockProduct, error)

	// Users
	GetNewUsersCount(ctx context.Context, startDate time.Time) (int64, error)
	GetUserGrowthAnalytics(ctx context.Context, days int) ([]models.UserGrowthAnalytics, error)

	// System
	GetRecentActivity(ctx context.Context, limit int) ([]models.ActivityLog, error)
	CheckDatabaseHealth(ctx context.Context) error
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

// ==================== Stats Methods ====================

func (r *dashboardRepository) CountTotalUsers(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.User{}).Count(&count).Error
	return count, err
}

func (r *dashboardRepository) CountActiveUsers(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("is_active = ?", true).
		Count(&count).Error
	return count, err
}

func (r *dashboardRepository) CountTotalProducts(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Product{}).Count(&count).Error
	return count, err
}

func (r *dashboardRepository) CountTotalOrders(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Order{}).Count(&count).Error
	return count, err
}

func (r *dashboardRepository) CountPendingOrders(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Order{}).
		Where("status = ?", "pending").
		Count(&count).Error
	return count, err
}

func (r *dashboardRepository) CountTotalCategories(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Category{}).Count(&count).Error
	return count, err
}

func (r *dashboardRepository) CountTotalRoles(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Role{}).Count(&count).Error
	return count, err
}

// ==================== Revenue Methods ====================

func (r *dashboardRepository) GetTotalRevenue(ctx context.Context) (float64, error) {
	var total float64
	err := r.db.WithContext(ctx).
		Model(&models.Order{}).
		Where("status IN ?", []string{"completed", "paid"}).
		Select("COALESCE(SUM(subtotal), 0)").
		Scan(&total).Error
	return total, err
}

func (r *dashboardRepository) GetRevenueByPeriod(ctx context.Context, startDate, endDate time.Time) (float64, error) {
	var total float64
	err := r.db.WithContext(ctx).
		Model(&models.Order{}).
		Where("status IN ? AND created_at BETWEEN ? AND ?",
			[]string{"completed", "paid"}, startDate, endDate).
		Select("COALESCE(SUM(subtotal), 0)").
		Scan(&total).Error
	return total, err
}

func (r *dashboardRepository) GetRevenueToday(ctx context.Context) (float64, error) {
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	return r.GetRevenueByPeriod(ctx, today, tomorrow)
}

func (r *dashboardRepository) GetRevenueThisMonth(ctx context.Context) (float64, error) {
	now := time.Now()
	firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	lastDay := firstDay.AddDate(0, 1, 0)
	return r.GetRevenueByPeriod(ctx, firstDay, lastDay)
}

// ==================== Orders Methods ====================

func (r *dashboardRepository) GetRecentOrders(ctx context.Context, limit int) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("ColorVarian").
		Preload("SizeVarian").
		Order("created_at DESC").
		Limit(limit).
		Find(&orders).Error
	return orders, err
}

func (r *dashboardRepository) GetOrdersByStatus(ctx context.Context) ([]models.OrderStatusCount, error) {
	var results []models.OrderStatusCount
	err := r.db.WithContext(ctx).
		Model(&models.Order{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&results).Error
	return results, err
}

func (r *dashboardRepository) GetOrderAnalytics(ctx context.Context, days int) ([]models.OrderAnalytics, error) {
	var analytics []models.OrderAnalytics
	startDate := time.Now().AddDate(0, 0, -days)

	err := r.db.WithContext(ctx).
		Model(&models.Order{}).
		Select("DATE(created_at) as date, COUNT(*) as order_count, COALESCE(SUM(subtotal), 0) as revenue").
		Where("created_at >= ? AND status IN ?", startDate, []string{"completed", "paid"}).
		Group("DATE(created_at)").
		Order("DATE(created_at)").
		Scan(&analytics).Error

	return analytics, err
}

// ==================== Products Methods ====================

func (r *dashboardRepository) GetTopSellingProducts(ctx context.Context, limit int) ([]models.TopProduct, error) {
	var products []models.TopProduct
	err := r.db.WithContext(ctx).
		Model(&models.Order{}).
		Select("products.id as product_id, products.name as product_name, SUM(orders.quantity) as total_sold, SUM(orders.subtotal) as revenue").
		Joins("JOIN products ON products.id = orders.product_id").
		Where("orders.status IN ?", []string{"completed", "paid"}).
		Group("products.id, products.name").
		Order("total_sold DESC").
		Limit(limit).
		Scan(&products).Error
	return products, err
}

func (r *dashboardRepository) GetLowStockProducts(ctx context.Context, threshold int, limit int) ([]models.LowStockProduct, error) {
	var products []models.LowStockProduct
	err := r.db.WithContext(ctx).
		Model(&models.SizeVarian{}).
		Select(`
			products.id as product_id,
			products.name as product_name,
			color_varian.name as color_name,
			size_varian.size as size,
			size_varian.stock as stock,
			products.price as price
		`).
		Joins("JOIN color_varian ON color_varian.id = size_varian.color_varian_id").
		Joins("JOIN products ON products.id = color_varian.product_id").
		Where("size_varian.stock <= ?", threshold).
		Order("size_varian.stock ASC").
		Limit(limit).
		Scan(&products).Error
	return products, err
}

// ==================== Users Methods ====================

func (r *dashboardRepository) GetNewUsersCount(ctx context.Context, startDate time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("created_at >= ?", startDate).
		Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetUserGrowthAnalytics(ctx context.Context, days int) ([]models.UserGrowthAnalytics, error) {
	var analytics []models.UserGrowthAnalytics
	startDate := time.Now().AddDate(0, 0, -days)

	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Select("DATE(created_at) as date, COUNT(*) as user_count").
		Where("created_at >= ?", startDate).
		Group("DATE(created_at)").
		Order("DATE(created_at)").
		Scan(&analytics).Error

	return analytics, err
}

// ==================== System Methods ====================

func (r *dashboardRepository) GetRecentActivity(ctx context.Context, limit int) ([]models.ActivityLog, error) {
	var activities []models.ActivityLog
	err := r.db.WithContext(ctx).
		Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Find(&activities).Error
	return activities, err
}

func (r *dashboardRepository) CheckDatabaseHealth(ctx context.Context) error {
	return r.db.WithContext(ctx).Raw("SELECT 1").Error
}
