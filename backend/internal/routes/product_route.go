package routes

import (
	"e-commerce/backend/internal/handler"
	mw "e-commerce/backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func ProductRoutes(r chi.Router, h *handler.ProductHandler, deps Dependencies) {
	r.Route("/products", func(r chi.Router) {

		r.Get("/", h.GetAllProducts)
		r.Get("/{id}", h.GetProductByID)

		r.Group(func(r chi.Router) {
			r.Use(mw.AuthMiddleware(deps.UserService, deps.JWTService))
			r.Use(mw.RequireAdminArea(deps.RBACService))
			r.Post("/", h.CreateProduct)
			r.Post("/{id}/color", h.AddColorVariant)
			r.Put("/{id}", h.UpdateProduct)
			r.Delete("/{id}", h.DeleteProduct)
		})
	})
}
