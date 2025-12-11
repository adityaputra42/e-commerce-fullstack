package middleware

import (
	"context"
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/services"
	"e-commerce/backend/internal/utils"
	"encoding/json"
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

// AuthMiddleware - Native Go version
func AuthMiddleware(userService services.UserService, jwtService *utils.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				sendError(w, http.StatusUnauthorized, "unauthorized", "Authorization header is required")
				return
			}

			tokenParts := strings.SplitN(authHeader, " ", 2)
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				sendError(w, http.StatusUnauthorized, "unauthorized", "Invalid authorization header format")
				return
			}

			tokenString := tokenParts[1]
			claims, err := jwtService.ValidateAccessToken(tokenString)
			if err != nil {
				sendError(w, http.StatusUnauthorized, "unauthorized", "Invalid or expired token")
				return
			}

			user, err := userService.GetUserById(uint(claims.UserID))
			if err != nil {
				sendError(w, http.StatusUnauthorized, "unauthorized", "User not found")
				return
			}

			if !user.IsActive {
				sendError(w, http.StatusUnauthorized, "unauthorized", "Account is deactivated")
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, user)
			ctx = context.WithValue(ctx, UserIDContextKey, user.ID)
			ctx = context.WithValue(ctx, RoleIDContextKey, user.RoleID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Helper functions to get data from context
func GetUserFromContext(r *http.Request) *models.User {
	user, ok := r.Context().Value(UserContextKey).(*models.User)
	if !ok {
		return nil
	}
	return user
}

func GetUserIDFromContext(r *http.Request) uint {
	userID, ok := r.Context().Value(UserIDContextKey).(uint)
	if !ok {
		return 0
	}
	return userID
}

func GetRoleIDFromContext(r *http.Request) uint {
	roleID, ok := r.Context().Value(RoleIDContextKey).(uint)
	if !ok {
		return 0
	}
	return roleID
}
