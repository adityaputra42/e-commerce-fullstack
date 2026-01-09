package routes

import (
	"e-commerce/backend/internal/handler"
	mw "e-commerce/backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

// ShippingRoutes sets up routes for shipping method management
func ShippingRoutes(r chi.Router, h *handler.ShippingHandler, deps Dependencies) {
	authMiddleware := mw.AuthMiddleware(deps.UserService, deps.JWTService)

	r.Route("/shipping", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)
			r.Get("/", h.GetAllShipping)
			r.Get("/{id}", h.GetShippingByID)
		})

		// Admin only routes
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)
			r.Use(mw.RequireAdminArea(deps.RBACService))
			r.Post("/", h.CreateShipping)
			r.Put("/{id}", h.UpdateShipping)
			r.Delete("/{id}", h.DeleteShipping)
		})

	})
}
