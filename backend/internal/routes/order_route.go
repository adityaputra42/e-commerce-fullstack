package routes

import (
	"e-commerce/backend/internal/handler"
	"e-commerce/backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func OrderRoutes(r chi.Router, h *handler.OrderHandler, deps Dependencies) {

	authMiddleware := middleware.AuthMiddleware(deps.UserService, deps.JWTService)

	r.Route("/orders", func(r chi.Router) {
		r.Use(authMiddleware)

		r.Get("/", h.GetAllOrders)
		r.Get("/{id}", h.GetOrderByID)
		r.Patch("/{id}/cancel", h.CancelOrder)

		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireAdminArea(deps.RBACService))
			r.Put("/{id}", h.UpdateOrder)
			r.Delete("/{id}", h.DeleteOrder)
		})

	})
}
