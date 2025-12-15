package routes

import (
	"e-commerce/backend/internal/di"
	"e-commerce/backend/internal/services"
	"e-commerce/backend/internal/utils"
	"log"
	"net/http"
	"time"

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
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.AllowContentType("application/json", "multipart/form-data"))
	r.Get("/", handler.HealthHandler.Root)
	r.Get("/health", handler.HealthHandler.HealthCheck)
	r.Head("/health", handler.HealthHandler.HealthCheck) // T

	deps := buildDependencies(handler)
	r.Route("/api/v1", func(api chi.Router) {
		api.Get("/health", handler.HealthHandler.HealthCheck)
		api.Head("/health", handler.HealthHandler.HealthCheck) // T
		AuthRoutes(api, handler.AuthHandler)
		UserRoutes(api, handler.UserHandler, deps)
		ProductRoutes(api, handler.ProductHandler, deps)
		AddressRoutes(api, handler.AddressHandler, deps)
		RoleRoutes(api, handler.RoleHandler, deps)
		ShippingRoutes(api, handler.ShippingHandler, deps)
		OrderRoutes(api, handler.OrderHandler, deps)
		TransactionRoutes(api, handler.TransactionHandler, deps)
		PaymentMethodRoutes(api, handler.PaymentMethodHandler, deps)
		PaymentRoutes(api, handler.PaymentHandler, deps)
	})

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("📌 Route registered: %s %s", method, route)
		return nil
	}
	if err := chi.Walk(r, walkFunc); err != nil {
		log.Printf("⚠️  Error walking routes: %s", err.Error())
	}

	return r
}

func buildDependencies(handler *di.Handler) Dependencies {
	return Dependencies{
		RBACService: handler.RBACService,
		UserService: handler.UserService,
		JWTService:  handler.JWTService,
	}
}
