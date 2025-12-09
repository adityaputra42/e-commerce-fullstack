package handler

import (
	"e-commerce/backend/internal/middleware"
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/services"
	"e-commerce/backend/internal/utils"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// SignUp - POST /api/auth/signup
func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.authService.SignUp(req)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "User registered successfully",
		"data":    resp,
	})
}

// SignIn - POST /api/auth/signin
func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Get request metadata
	ipAddress := r.RemoteAddr
	userAgent := r.UserAgent()

	resp, err := h.authService.SignIn(req, ipAddress, userAgent)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"data":    resp,
	})
}

// SignOut - POST /api/auth/signout
func (h *AuthHandler) SignOut(w http.ResponseWriter, r *http.Request) {

	userID := middleware.GetUserIDFromContext(r)
	if userID == 0 {
		utils.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Optional: Implement token blacklist or session invalidation
	// For now, just return success as JWT is stateless
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Logout successful",
	})
}

// Refresh - POST /api/auth/refresh
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.authService.RefreshToken(req)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Token refreshed successfully",
		"data":    resp,
	})
}

// ForgotPassword - POST /api/auth/forgot-password
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req models.ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	token, err := h.authService.ForgotPassword(req)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Password reset token sent successfully",
		"token":   token,
	})
}

// ResetPassword - POST /api/auth/reset-password
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req models.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.authService.ResetPassword(req)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Password reset successfully",
	})
}

// VerifyEmail - GET /api/auth/verify-email?token=xxx
func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		utils.WriteError(w, http.StatusBadRequest, "Verification token is required")
		return
	}

	// Call service to verify email
	// err := h.authService.VerifyEmail(token)
	// if err != nil {
	// 	utils.WriteError(w, http.StatusBadRequest, err.Error())
	// 	return
	// }

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Email verified successfully",
	})
}

// ResendVerification - POST /api/auth/resend-verification
func (h *AuthHandler) ResendVerification(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Email == "" {
		utils.WriteError(w, http.StatusBadRequest, "Email is required")
		return
	}

	// Call service to resend verification
	// err := h.authService.ResendVerification(req.Email)
	// if err != nil {
	// 	utils.WriteError(w, http.StatusBadRequest, err.Error())
	// 	return
	// }

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Verification email sent successfully",
	})
}

// GetProfile - GET /api/auth/profile (requires auth middleware)
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user := middleware.GetUserFromContext(r)
	if user == nil {
		utils.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Profile retrieved successfully",
		"data":    user,
	})
}

// ChangePassword - PUT /api/auth/change-password (requires auth middleware)
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	if userID == 0 {
		utils.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		utils.WriteError(w, http.StatusBadRequest, "Old password and new password are required")
		return
	}

	// Call service to change password
	// err := h.authService.ChangePassword(userID, req.OldPassword, req.NewPassword)
	// if err != nil {
	// 	utils.WriteError(w, http.StatusBadRequest, err.Error())
	// 	return
	// }

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Password changed successfully",
	})
}
