package routes

import (
	"e-commerce/backend/internal/handler"
	"e-commerce/backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func AuthRoutes(r chi.Router, h *handler.AuthHandler, deps Dependencies) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", h.SignIn)
		r.Post("/admin/login", h.AdminLogin)
		r.Post("/logout", h.SignOut)
		r.Post("/register", h.SignUp)
		r.With(middleware.AuthMiddleware(deps.UserService, deps.JWTService)).Get("/profile", h.GetProfile)
		r.With(middleware.AuthMiddleware(deps.UserService, deps.JWTService)).Post("/refresh", h.Refresh)
		r.Post("/resend-verification", h.ResendVerification)
		r.Post("/reset-password", h.ResetPassword)
		r.Post("/forgot-password", h.ForgotPassword)
		r.Get("/verify-email", h.VerifyEmail)
		r.With(middleware.AuthMiddleware(deps.UserService, deps.JWTService)).Put("/change-password", h.ChangePassword)

	})
}
