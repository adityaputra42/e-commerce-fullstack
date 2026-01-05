package routes

import (
	"e-commerce/backend/internal/handler"

	"github.com/go-chi/chi/v5"
)

func CategoryRoutes(r chi.Router, h *handler.CategoryHandler, deps Dependencies) {
	r.Route("/categories", func(r chi.Router) {
		r.Post("/", h.CreateCategory)
		r.Get("/", h.GetAllCategories)
		r.Get("/{id}", h.GetCategoryById)
		r.Put("/{id}", h.UpdateCategory)
		r.Delete("/{id}", h.DeleteCategory)
	})
}
