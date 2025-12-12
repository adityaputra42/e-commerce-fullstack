package routes

import (
	"e-commerce/backend/internal/handler"
	"e-commerce/backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

// TransactionRoutes sets up routes for transaction management
func TransactionRoutes(r chi.Router, h *handler.TransactionHandler, deps Dependencies) {
	authMiddleware := middleware.AuthMiddleware(deps.UserService, deps.JWTService)

	r.Route("/transactions", func(r chi.Router) {

		r.Use(authMiddleware)

		r.Post("/", h.CreateTransaction)
		r.Get("/", h.GetAllTransactions)
		r.Get("/{tx_id}", h.GetTransactionByID)

		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRole(deps.RBACService, "admin"))
			r.Put("/{tx_id}", h.UpdateTransaction)
		})
	})
}
