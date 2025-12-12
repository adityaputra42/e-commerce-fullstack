package routes

import (
	"e-commerce/backend/internal/handler"
	"e-commerce/backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

// PaymentRoutes sets up routes for payment management
func PaymentRoutes(r chi.Router, h *handler.PaymentHandler, deps Dependencies) {
	authMiddleware := middleware.AuthMiddleware(deps.UserService, deps.JWTService)

	r.Route("/payments", func(r chi.Router) {
		r.Use(authMiddleware)
		r.Post("/", h.CreatePayment)
		r.Get("/", h.GetAllPayments)
		r.Get("/{id}", h.GetPaymentByID)

		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRole(deps.RBACService, "admin"))
			r.Put("/{id}", h.UpdatePayment)
			r.Delete("/{id}", h.DeletePayment)
		})

	})
}
