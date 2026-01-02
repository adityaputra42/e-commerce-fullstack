package handler

import (
	"e-commerce/backend/internal/middleware"
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/services"
	"e-commerce/backend/internal/utils"
	"encoding/json"
	"fmt"
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

// SignUp - POST /api/auth/register
// @Summary SignUp user
// @Description Register a new user with email, username, and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Register request"
// @Success 201 {object} utils.Response{data=models.TokenResponse} "User registered successfully"
// @Failure 400 {object} utils.Response "Invalid request body or service error"
// @Router /auth/register [post]
func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	resp, err := h.authService.SignUp(req)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, "User registered successfully", resp)
}

// SignIn - POST /api/auth/login
// @Summary SignIn user
// @Description Login with email and password to get access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login request"
// @Success 200 {object} utils.Response{data=models.TokenResponse} "Login successful"
// @Failure 401 {object} utils.Response "Invalid credentials"
// @Router /auth/login [post]
func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Get request metadata
	ipAddress := r.RemoteAddr
	userAgent := r.UserAgent()

	resp, err := h.authService.SignIn(req, ipAddress, userAgent)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Login successful", resp)
}

// SignOut - POST /api/v1/auth/logout
// @Summary SignOut user
// @Description Logout from the application
// @Tags Auth
// @Produce json
// @Success 200 {object} utils.Response "Logout successful"
// @Router /auth/logout [post]
// @Security Bearer
func (h *AuthHandler) SignOut(w http.ResponseWriter, r *http.Request) {

	userID := middleware.GetUserIDFromContext(r)
	if userID == 0 {
		utils.WriteError(w, http.StatusUnauthorized, "User not authenticated", fmt.Errorf("User not authenticated"))
		return
	}

	// Optional: Implement token blacklist or session invalidation
	// For now, just return success as JWT is stateless
	utils.WriteJSON(w, http.StatusOK, "Logout successful", nil)
}

// Refresh - POST /api/v1/auth/refresh
// @Summary Refresh token
// @Description Get a new access token using refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} utils.Response{data=models.RefreshTokenResponse} "Token refreshed successfully"
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	resp, err := h.authService.RefreshToken(req)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Token refreshed successfully", resp)
}

// ForgotPassword - POST /api/v1/auth/forgot-password
// @Summary Forgot password
// @Description Request password reset token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.ForgotPasswordRequest true "Forgot password request"
// @Success 200 {object} utils.Response "Password reset token sent successfully"
// @Router /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req models.ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	token, err := h.authService.ForgotPassword(req)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Password reset token sent successfully", token)
}

// ResetPassword - POST /api/auth/reset-password
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req models.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	err := h.authService.ResetPassword(req)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Password reset successfully", nil)
}

// VerifyEmail - GET /api/auth/verify-email?token=xxx
func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		utils.WriteError(w, http.StatusBadRequest, "Verification token is required", fmt.Errorf("Verification token is required"))
		return
	}

	// Call service to verify email
	// err := h.authService.VerifyEmail(token)
	// if err != nil {
	// 	utils.WriteError(w, http.StatusBadRequest, err.Error())
	// 	return
	// }

	utils.WriteJSON(w, http.StatusOK, "Email verified successfully", nil)
}

// ResendVerification - POST /api/auth/resend-verification
func (h *AuthHandler) ResendVerification(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if req.Email == "" {
		utils.WriteError(w, http.StatusBadRequest, "Email is required", fmt.Errorf("Email is required"))
		return
	}

	// Call service to resend verification
	// err := h.authService.ResendVerification(req.Email)
	// if err != nil {
	// 	utils.WriteError(w, http.StatusBadRequest, err.Error())
	// 	return
	// }

	utils.WriteJSON(w, http.StatusOK, "Verification email sent successfully", nil)
}

// GetProfile - GET /api/v1/auth/profile
// @Summary Get current user profile
// @Description Get information about the currently logged in user
// @Tags Auth
// @Produce json
// @Success 200 {object} utils.Response{data=models.User} "Profile retrieved successfully"
// @Router /auth/profile [get]
// @Security Bearer
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user := middleware.GetUserFromContext(r)
	if user == nil {
		utils.WriteError(w, http.StatusUnauthorized, "User not authenticated", fmt.Errorf("User not authenticated"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Profile retrieved successfully", user)
}

// ChangePassword - PUT /api/auth/change-password (requires auth middleware)
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	if userID == 0 {
		utils.WriteError(w, http.StatusUnauthorized, "User not authenticated", fmt.Errorf("User not authenticated"))
		return
	}

	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		utils.WriteError(w, http.StatusBadRequest, "Old password and new password are required", fmt.Errorf("Old password and new password are required"))
		return
	}

	// Call service to change password
	// err := h.authService.ChangePassword(userID, req.OldPassword, req.NewPassword)
	// if err != nil {
	// 	utils.WriteError(w, http.StatusBadRequest, err.Error())
	// 	return
	// }

	utils.WriteJSON(w, http.StatusOK, "Password changed successfully", nil)
}
