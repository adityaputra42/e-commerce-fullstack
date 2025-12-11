package routes

import (
	"e-commerce/backend/internal/di"
	"e-commerce/backend/internal/services"
	"e-commerce/backend/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Dependencies struct {
	RBACService services.RBACService
	UserService services.UserService
	JWTService  *utils.JWTService
}

func SetupRoutes(handler *di.Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json", "multipart/form-data"))

	deps := buildDependencies(handler)
	r.Route("/api/v1", func(api chi.Router) {
		AuthRoutes(api, handler.AuthHandler)
		UserRoutes(api, handler.UserHandler, deps)
		ProductRoutes(api, handler.ProductHandler, deps)
		AddressRoutes(api, handler.AddressHandler, deps)
	})

	return r
}

func buildDependencies(handler *di.Handler) Dependencies {
	return Dependencies{
		RBACService: handler.RBACService,
		UserService: handler.UserService,
		JWTService:  handler.JWTService,
	}
}
