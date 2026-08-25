package handler

import (
	"e-commerce/backend/internal/config"
	"e-commerce/backend/internal/middleware"
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/services"
	"e-commerce/backend/internal/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type AuthHandler struct {
	authService services.AuthService
	cfg         *config.Config
}

func NewAuthHandler(authService services.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		cfg:         cfg,
	}
}

// isSecureCookie menentukan apakah cookie refresh token butuh flag Secure
// (cuma dikirim lewat HTTPS). WAJIB true di production; boleh false untuk
// dev lokal di HTTP biasa supaya browser tidak diam-diam menolak cookie-nya.
func (h *AuthHandler) isSecureCookie() bool {
	return h.cfg.Server.Env == "production"
}

// respondWithTokens mengirim access token di JSON body (dipegang frontend
// di memory, TIDAK di localStorage) dan menaruh refresh token di httpOnly
// cookie (lihat utils/cookie.go) — bukan di JSON body seperti sebelumnya.
func (h *AuthHandler) respondWithTokens(w http.ResponseWriter, status int, message string, resp *models.TokenResponse) {
	utils.SetRefreshTokenCookie(w, resp.RefreshToken, time.Now().Add(h.cfg.JWT.RefreshTokenExpiry), h.isSecureCookie())
	// Refresh token TIDAK disertakan di response body — sudah ada di cookie.
	resp.RefreshToken = ""
	utils.WriteJSON(w, status, message, resp)
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

	h.respondWithTokens(w, http.StatusCreated, "User registered successfully", resp)
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

	h.respondWithTokens(w, http.StatusOK, "Login successful", resp)
}

// SignIn - POST /api/auth/admin/login
// @Summary SignIn Admin
// @Description Login with email and password to get access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login request"
// @Success 200 {object} utils.Response{data=models.TokenResponse} "Login successful"
// @Failure 401 {object} utils.Response "Invalid credentials"
// @Router /auth/login [post]
func (h *AuthHandler) AdminLogin(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	ipAddress := r.RemoteAddr
	userAgent := r.UserAgent()

	resp, err := h.authService.LoginAdmin(req, ipAddress, userAgent)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	h.respondWithTokens(w, http.StatusOK, "Login successful", resp)
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

	// SEBELUMNYA fungsi ini benar-benar no-op (cuma return success). Sekarang
	// setidaknya hapus cookie refresh token, supaya sesi tidak bisa dipulihkan
	// diam-diam lewat refresh setelah user logout dari device ini.
	// (JWT tetap stateless — access token yang sudah diterbitkan tetap valid
	// sampai expired; kalau butuh revoke instan, itu butuh token blacklist
	// terpisah, di luar scope perbaikan ini.)
	utils.ClearRefreshTokenCookie(w, h.isSecureCookie())

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

	// Prioritaskan cookie httpOnly (dipakai browser web). Fallback ke JSON
	// body untuk client yang tidak bisa pakai cookie (mis. aplikasi mobile),
	// supaya endpoint ini tidak cuma bisa dipakai dari browser.
	if cookieToken, err := utils.GetRefreshTokenFromCookie(r); err == nil && cookieToken != "" {
		req.RefreshToken = cookieToken
	} else if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RefreshToken == "" {
		utils.WriteError(w, http.StatusBadRequest, "Refresh token tidak ditemukan (cookie atau body)", err)
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
// @Description Kirim link reset password ke email jika akun terdaftar
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.ForgotPasswordRequest true "Forgot password request"
// @Success 200 {object} utils.Response "Jika email terdaftar, link reset telah dikirim"
// @Router /auth/forgot-password [post]
//
// PENTING (fix keamanan): handler ini SEBELUMNYA mengembalikan token reset
// password langsung di response API — artinya siapa pun yang tahu email
// seorang user bisa ambil-alih akunnya tanpa autentikasi apa pun. Sekarang
// token TIDAK PERNAH dikembalikan ke client (dikirim lewat email di dalam
// service), dan response selalu memakai pesan generik yang sama persis baik
// email-nya terdaftar maupun tidak — supaya endpoint ini juga tidak bisa
// dipakai untuk menebak email mana saja yang punya akun (user enumeration).
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req models.ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := h.authService.ForgotPassword(req); err != nil {
		// Kegagalan infrastruktur (mis. SMTP down) tetap di-log di server
		// lewat AuthService, tapi client tetap melihat pesan generik yang
		// sama supaya tidak membocorkan informasi apa pun.
		log.Printf("ForgotPassword request failed: %v", err)
	}

	utils.WriteJSON(w, http.StatusOK, "Jika email terdaftar, kami telah mengirim link reset password", nil)
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
