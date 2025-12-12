package routes

import (
	"e-commerce/backend/internal/handler"
	mw "e-commerce/backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

// PaymentMethodRoutes sets up routes for payment method management
func PaymentMethodRoutes(r chi.Router, h *handler.PaymentMethodHandler, deps Dependencies) {
	authMiddleware := mw.AuthMiddleware(deps.UserService, deps.JWTService)

	r.Route("/payment-methods", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)
			r.Get("/", h.GetAllPaymentMethods)
			r.Get("/{id}", h.GetPaymentMethodByID)
		})

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)
			r.Use(mw.RequireRole(deps.RBACService, "admin"))
			r.Post("/", h.CreatePaymentMethod)
			r.Put("/{id}", h.UpdatePaymentMethod)
			r.Delete("/{id}", h.DeletePaymentMethod)
		})

	})
}
