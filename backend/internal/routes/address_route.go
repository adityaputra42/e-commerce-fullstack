package routes

import (
	"e-commerce/backend/internal/handler"
	mw "e-commerce/backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func AddressRoutes(r chi.Router, h *handler.AddressHandler, deps Dependencies) {
	r.Route("/address", func(r chi.Router) {

		r.Group(func(r chi.Router) {
			r.Use(mw.AuthMiddleware(deps.UserService, deps.JWTService))
			r.Get("/", h.GetAddresses)
			r.Get("/{id}", h.GetAddressByID)
			r.Post("/", h.CreateAddress)
			r.Put("/{id}", h.UpdateAddress)
			r.Delete("/{id}", h.DeleteAddress)
		})
	})
}
