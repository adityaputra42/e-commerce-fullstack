//go:build wireinject
// +build wireinject

package di

import (
	"e-commerce/backend/internal/config"
	"e-commerce/backend/internal/database"
	"e-commerce/backend/internal/handler"
	"e-commerce/backend/internal/repository"
	"e-commerce/backend/internal/services"
	"e-commerce/backend/internal/utils"

	"github.com/google/wire"
	"gorm.io/gorm"
)

// ProvideDB provides database connection
func ProvideDB() *gorm.DB {
	return database.DB
}

// ProvideJWTService provides JWT service instance with config
func ProvideJWTService(config *config.Config) *utils.JWTService {
	return utils.NewJWTService(config)
}

// Repository Providers
var repositorySet = wire.NewSet(
	ProvideDB,
	repository.NewCategoryRepository,
	repository.NewProductRepository,
	repository.NewAddressRepository,
	repository.NewOrderRepository,
	repository.NewPasswordResetTokenRepository,
	repository.NewPaymentMethodRepository,
	repository.NewPaymentRepository,
	repository.NewPermissionRepository,
	repository.NewRBACRepository,
	repository.NewShippingRepository,
	repository.NewTransactionRepository,
	repository.NewUserReposiory,
	repository.NewActivityLogRepository,
	repository.NewRoleRepository,
	repository.NewDashboardRepository,
)

// Service Providers
var serviceSet = wire.NewSet(
	services.NewCategoryService,
	services.NewProductService,
	services.NewAddressService,
	services.NewAuthService,
	services.NewOrderService,
	services.NewPaymentMethodService,
	services.NewPaymentService,
	services.NewRBACService,
	services.NewRoleService,
	services.NewShippingService,
	services.NewTransactionService,
	services.NewUserService,
	services.NewDashboardService,
)

// Utils Providers
var utilsSet = wire.NewSet(
	ProvideJWTService,
)

// Handler Providers
var handlerSet = wire.NewSet(
	handler.NewCategoryHandler,
	handler.NewProductHandler,
	handler.NewAddressHandler,
	handler.NewUserHandler,
	handler.NewAuthHandler,
	handler.NewRoleHandler,
	handler.NewOrderHandler,
	handler.NewShippingHandler,
	handler.NewPaymentHandler,
	handler.NewPaymentMethodHandler,
	handler.NewTransactionHandler,
	handler.NewHealthHandler,
	handler.NewDashboardHandler,
)

// InitializeAllHandler initializes all handler with config
func InitializeAllHandler(config *config.Config) *Handler {
	wire.Build(
		repositorySet,
		serviceSet,
		utilsSet,
		handlerSet,
		NewHandler,
	)
	return &Handler{}
}

// Handler struct contains all handler and services
type Handler struct {
	CategoryHandler      *handler.CategoryHandler
	ProductHandler       *handler.ProductHandler
	AddressHandler       *handler.AddressHandler
	AuthHandler          *handler.AuthHandler
	UserHandler          *handler.UserHandler
	RoleHandler          *handler.RoleHandler
	OrderHandler         *handler.OrderHandler
	PaymentHandler       *handler.PaymentHandler
	PaymentMethodHandler *handler.PaymentMethodHandler
	ShippingHandler      *handler.ShippingHandler
	TransactionHandler   *handler.TransactionHandler
	HealthHandler        *handler.HealthHandler
	DashboardHandler     *handler.DashboardHandler

	// Services untuk middleware
	RBACService services.RBACService
	UserService services.UserService
	JWTService  *utils.JWTService
}

// NewHandler creates new Handler instance
func NewHandler(
	categoryHandler *handler.CategoryHandler,
	productHandler *handler.ProductHandler,

	addressHandler *handler.AddressHandler,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	roleHandler *handler.RoleHandler,
	orderHandler *handler.OrderHandler,
	paymentHandler *handler.PaymentHandler,
	paymentMethodHandler *handler.PaymentMethodHandler,
	shippingHandler *handler.ShippingHandler,
	transactionHandler *handler.TransactionHandler,
	healthHandler *handler.HealthHandler,
	dashboardHandler *handler.DashboardHandler,
	rbacService services.RBACService,
	userService services.UserService,
	jwtService *utils.JWTService,
) *Handler {
	return &Handler{
		CategoryHandler:      categoryHandler,
		ProductHandler:       productHandler,
		AddressHandler:       addressHandler,
		AuthHandler:          authHandler,
		UserHandler:          userHandler,
		RoleHandler:          roleHandler,
		OrderHandler:         orderHandler,
		PaymentHandler:       paymentHandler,
		PaymentMethodHandler: paymentMethodHandler,
		ShippingHandler:      shippingHandler,
		TransactionHandler:   transactionHandler,
		HealthHandler:        healthHandler,
		DashboardHandler:     dashboardHandler,
		RBACService:          rbacService,
		UserService:          userService,
		JWTService:           jwtService,
	}
}
