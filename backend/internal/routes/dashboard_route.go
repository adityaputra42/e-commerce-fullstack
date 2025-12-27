package routes

import (
	"e-commerce/backend/internal/handler"
	"e-commerce/backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func DashboardRoutes(r chi.Router, dashboardHandler *handler.DashboardHandler, deps Dependencies) {

	authMiddleware := middleware.AuthMiddleware(deps.UserService, deps.JWTService)

	r.Route("/dashboard", func(r chi.Router) {
		r.Use(authMiddleware)
		r.Use(middleware.RequireRole(deps.RBACService, "admin"))

		r.Get("/stats", dashboardHandler.GetDashboardStats)
		r.Get("/revenue", dashboardHandler.GetRevenueStats)

		r.Route("/orders", func(r chi.Router) {
			r.Get("/stats", dashboardHandler.GetOrderStats)
			r.Get("/recent", dashboardHandler.GetRecentOrders)
		})

		r.Route("/products", func(r chi.Router) {
			r.Get("/top", dashboardHandler.GetTopProducts)
			r.Get("/low-stock", dashboardHandler.GetLowStockProducts)
		})

		r.Route("/analytics", func(r chi.Router) {
			r.Get("/orders", dashboardHandler.GetOrderAnalytics)
			r.Get("/users", dashboardHandler.GetUserGrowth)
		})

		r.Get("/health", dashboardHandler.GetSystemHealth)
		r.Get("/activity", dashboardHandler.GetRecentActivity)
	})
}
