package routes

import (
	"e-commerce/backend/internal/handler"
	mw "e-commerce/backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

// RoleRoutes sets up routes for role management
func RoleRoutes(r chi.Router, h *handler.RoleHandler, deps Dependencies) {
	authMiddleware := mw.AuthMiddleware(deps.UserService, deps.JWTService)

	r.Route("/roles", func(r chi.Router) {
		r.Use(authMiddleware)

		r.Group(func(r chi.Router) {
			r.Get("/permissions", h.GetAllPermissions)
			r.Get("/", h.GetAllRoles)
			r.Get("/{id}", h.GetRoleByID)
			r.Get("/{id}/permissions", h.GetRolePermissions)
		})

		r.Group(func(r chi.Router) {
			r.Use(mw.RequireRole(deps.RBACService, "admin"))
			r.Post("/", h.CreateRole)
			r.Put("/{id}", h.UpdateRole)
			r.Post("/{id}/permissions", h.AssignPermissions)
			r.Delete("/{id}", h.DeleteRole)
		})

	})
}
