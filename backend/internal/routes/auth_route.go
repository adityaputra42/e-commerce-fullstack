package routes

import (
	"e-commerce/backend/internal/handler"

	"github.com/go-chi/chi/v5"
)

func AuthRoutes(r chi.Router, h *handler.AuthHandler) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", h.SignIn)
		r.Post("/logout", h.SignOut)
		r.Post("/register", h.SignUp)
		r.Get("/profile", h.GetProfile)
		r.Post("/refresh", h.Refresh)
		r.Post("/resend-verification", h.ResendVerification)
		r.Post("/reset-password", h.ResetPassword)
		r.Post("/forgot-password", h.ForgotPassword)
		r.Get("/verify-email", h.VerifyEmail)
		r.Put("/change-password", h.ChangePassword)
	})
}
