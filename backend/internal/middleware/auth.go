package middleware

import (
	"context"
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/services"
	"e-commerce/backend/internal/utils"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// Helper function untuk mengirim error response
func sendError(w http.ResponseWriter, statusCode int, errorType, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":   errorType,
		"message": message,
	})
}

// Context keys untuk menyimpan data user
type contextKey string

const (
	UserContextKey   contextKey = "user"
	UserIDContextKey contextKey = "user_id"
	RoleIDContextKey contextKey = "role_id"
)

// AuthMiddleware - Native Go version dengan debug logging
func AuthMiddleware(userService services.UserService, jwtService *utils.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("🔐 AUTH MIDDLEWARE - Path: %s %s", r.Method, r.URL.Path)

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Println("❌ AUTH ERROR: No Authorization header")
				sendError(w, http.StatusUnauthorized, "unauthorized", "Authorization header is required")
				return
			}

			tokenParts := strings.SplitN(authHeader, " ", 2)
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				log.Printf("❌ AUTH ERROR: Invalid header format: %s", authHeader)
				sendError(w, http.StatusUnauthorized, "unauthorized", "Invalid authorization header format")
				return
			}

			tokenString := tokenParts[1]
			log.Printf("🔑 Validating token (first 20 chars): %s...", tokenString[:min(20, len(tokenString))])

			// Validate token
			claims, err := jwtService.ValidateAccessToken(tokenString)
			if err != nil {
				log.Printf("❌ AUTH ERROR: Token validation failed: %v", err)
				sendError(w, http.StatusUnauthorized, "unauthorized", "Invalid or expired token")
				return
			}

			// 🔍 DEBUG: Log claims yang berhasil di-parse
			log.Printf("✅ Token Valid - Claims: UserID=%d, Email=%s, RoleID=%d",
				claims.UserID, claims.Email, claims.RoleID)

			// CRITICAL: Validasi UserID tidak 0
			if claims.UserID == 0 {
				log.Println("❌ AUTH ERROR: UserID is 0 in claims")
				sendError(w, http.StatusUnauthorized, "unauthorized", "Invalid user ID in token")
				return
			}

			// Get user from database
			log.Printf("🔍 Fetching user from database - UserID: %d", claims.UserID)
			user, err := userService.GetUserById(uint(claims.UserID))
			if err != nil {
				log.Printf("❌ AUTH ERROR: User not found in database - UserID: %d, Error: %v",
					claims.UserID, err)
				sendError(w, http.StatusUnauthorized, "unauthorized", "User not found")
				return
			}

			// 🔍 DEBUG: Log user yang ditemukan
			if user == nil {
				log.Printf("❌ AUTH ERROR: User is nil after query - UserID: %d", claims.UserID)
				sendError(w, http.StatusUnauthorized, "unauthorized", "User not found")
				return
			}

			log.Printf("✅ User Found - ID: %d, Email: %s, IsActive: %v",
				user.ID, user.Email, user.IsActive)

			// Validasi user.ID dari database
			if user.ID == 0 {
				log.Printf("❌ AUTH ERROR: User.ID is 0 after database query")
				sendError(w, http.StatusInternalServerError, "server_error", "Invalid user data")
				return
			}

			if !user.IsActive {
				log.Printf("⚠️  AUTH WARNING: User account is deactivated - UserID: %d", user.ID)
				sendError(w, http.StatusUnauthorized, "unauthorized", "Account is deactivated")
				return
			}

			// Set context
			ctx := context.WithValue(r.Context(), UserContextKey, user)
			ctx = context.WithValue(ctx, UserIDContextKey, user.ID)
			ctx = context.WithValue(ctx, RoleIDContextKey, user.RoleID)

			log.Printf("✅ AUTH SUCCESS - UserID: %d proceeding to handler", user.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Helper functions to get data from context
func GetUserFromContext(r *http.Request) *models.User {
	user, ok := r.Context().Value(UserContextKey).(*models.User)
	if !ok {
		log.Println("⚠️  WARNING: User not found in context")
		return nil
	}
	return user
}

func GetUserIDFromContext(r *http.Request) uint {
	userID, ok := r.Context().Value(UserIDContextKey).(uint)
	if !ok {
		log.Println("⚠️  WARNING: UserID not found in context")
		return 0
	}
	return userID
}

func GetRoleIDFromContext(r *http.Request) uint {
	roleID, ok := r.Context().Value(RoleIDContextKey).(uint)
	if !ok {
		log.Println("⚠️  WARNING: RoleID not found in context")
		return 0
	}
	return roleID
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
