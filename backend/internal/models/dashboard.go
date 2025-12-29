package models

import "time"

type OrderStatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

type OrderAnalytics struct {
	Date       string  `json:"date"`
	OrderCount int64   `json:"order_count"`
	Revenue    float64 `json:"revenue"`
}
type TopProduct struct {
	ProductID   int64   `json:"product_id" gorm:"column:product_id"`
	ProductName string  `json:"product_name" gorm:"column:product_name"`
	TotalSold   int64   `json:"total_sold" gorm:"column:total_sold"`
	Revenue     float64 `json:"revenue" gorm:"column:revenue"`
}

type LowStockProduct struct {
	ProductID   int64   `json:"product_id" gorm:"column:product_id"`
	ProductName string  `json:"product_name" gorm:"column:product_name"`
	ColorName   string  `json:"color_name" gorm:"column:color_name"`
	Size        string  `json:"size" gorm:"column:size"`
	Stock       int64   `json:"stock" gorm:"column:stock"`
	Price       float64 `json:"price" gorm:"column:price"`
}
type UserGrowthAnalytics struct {
	Date      string `json:"date"`
	UserCount int64  `json:"user_count"`
}
type DashboardStatsResponse struct {
	TotalUsers       int64 `json:"total_users"`
	ActiveUsers      int64 `json:"active_users"`
	TotalProducts    int64 `json:"total_products"`
	TotalOrders      int64 `json:"total_orders"`
	PendingOrders    int64 `json:"pending_orders"`
	TotalCategories  int64 `json:"total_categories"`
	TotalRoles       int64 `json:"total_roles"`
	NewUsersToday    int64 `json:"new_users_today"`
	NewUsersThisWeek int64 `json:"new_users_this_week"`
}

type RevenueStatsResponse struct {
	TotalRevenue     float64 `json:"total_revenue"`
	RevenueToday     float64 `json:"revenue_today"`
	RevenueThisMonth float64 `json:"revenue_this_month"`
	RevenueThisWeek  float64 `json:"revenue_this_week"`
}

type OrderStatsResponse struct {
	TotalOrders    int64              `json:"total_orders"`
	PendingOrders  int64              `json:"pending_orders"`
	OrdersByStatus []OrderStatusCount `json:"orders_by_status"`
}

type RecentOrdersResponse struct {
	Orders []OrderSummary `json:"orders"`
}

type OrderSummary struct {
	ID          string    `json:"id"`
	ProductName string    `json:"product_name"`
	ColorName   string    `json:"color_name"`
	Size        string    `json:"size"`
	Quantity    int64     `json:"quantity"`
	Subtotal    float64   `json:"subtotal"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type TopProductsResponse struct {
	Products []TopProduct `json:"products"`
}

type LowStockProductsResponse struct {
	Products []LowStockProduct `json:"products"`
}

type OrderAnalyticsResponse struct {
	Analytics []OrderAnalytics `json:"analytics"`
	Summary   AnalyticsSummary `json:"summary"`
}

type AnalyticsSummary struct {
	TotalOrders  int64   `json:"total_orders"`
	TotalRevenue float64 `json:"total_revenue"`
	AverageOrder float64 `json:"average_order"`
}

type UserGrowthResponse struct {
	Growth  []UserGrowthAnalytics `json:"growth"`
	Summary UserGrowthSummary     `json:"summary"`
}

type UserGrowthSummary struct {
	TotalNewUsers int64   `json:"total_new_users"`
	AveragePerDay float64 `json:"average_per_day"`
}

type SystemHealthResponse struct {
	DatabaseStatus string    `json:"database_status"`
	TotalRequests  int64     `json:"total_requests"`
	ActiveUsers    int64     `json:"active_users"`
	Timestamp      time.Time `json:"timestamp"`
}

type RecentActivityResponse struct {
	Activities []ActivitySummary `json:"activities"`
}

type ActivitySummary struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	FullName  string    `json:"full_name"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	Details   string    `json:"details"`
	IPAddress string    `json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
}
