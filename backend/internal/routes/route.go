package routes

import (
	"e-commerce/backend/internal/config"
	"e-commerce/backend/internal/di"
	"e-commerce/backend/internal/middleware"
	"e-commerce/backend/internal/services"
	"e-commerce/backend/internal/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "e-commerce/backend/docs"
)

type Dependencies struct {
	RBACService services.RBACService
	UserService services.UserService
	JWTService  *utils.JWTService
}

func SetupRoutes(handler *di.Handler, logger *logrus.Logger, cfg config.CORSConfig) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.Recovery(logger))
	r.Use(middleware.Logger(logger)) // Logger middleware
	r.Use(chimiddleware.Compress(5))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: cfg.AllowedOrigins,
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
			"X-Requested-With",
		},
		ExposedHeaders: []string{
			"Link",
		},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(chimiddleware.AllowContentType("application/json", "multipart/form-data"))

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
		DashboardRoutes(api, handler.DashboardHandler, deps)
		CategoryRoutes(api, handler.CategoryHandler, deps)
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // The url pointing to API definition
	))

	return r
}

func buildDependencies(handler *di.Handler) Dependencies {
	return Dependencies{
		RBACService: handler.RBACService,
		UserService: handler.UserService,
		JWTService:  handler.JWTService,
	}
}
