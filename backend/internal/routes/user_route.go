package routes

import (
	"e-commerce/backend/internal/handler"
	mw "e-commerce/backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func UserRoutes(r chi.Router, h *handler.UserHandler, deps Dependencies) {
	r.Route("/users", func(r chi.Router) {

		r.Group(func(r chi.Router) {
			r.Use(mw.AuthMiddleware(deps.UserService, deps.JWTService))
			r.Get("/me", h.GetCurrentUser)
			r.Put("/me/password", h.UpdateCurrentUserPassword)
			r.Group(func(r chi.Router) {
				r.Use(mw.RequireRole(deps.RBACService, "admin"))
				r.Get("/", h.GetUsers)
				r.Get("/{id}", h.GetUserById)
				r.Post("/", h.CreateUser)
				r.Delete("/{id}/delete", h.DeleteUser)
				r.Put("/{id}/activate", h.ActivateUser)
				r.Put("/{id}/deactivate", h.DeactivateUser)
				r.Post("/bulk", h.BulkUserActions)
				r.Put("/{id}/password", h.UpdatePassword)

			})
			r.Group(func(r chi.Router) {
				r.Use(mw.SelfOrPermission(deps.RBACService, "user", "update"))
				r.Put("/{userId}", h.UpdateUser)
			})
		})
	})
}
