package middleware

import (
	"context"
	"e-commerce/backend/internal/models"
	"e-commerce/backend/internal/services"
	"e-commerce/backend/internal/utils"
	"encoding/json"
	"fmt"
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

// RequirePermission -
func RequirePermission(rbacService services.RBACService, resource, action string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := GetUserIDFromContext(r)
			if userID == 0 {
				sendError(w, http.StatusUnauthorized, "unauthorized", "User not authenticated")
				return
			}

			hasPermission, err := rbacService.CheckPermission(userID, resource, action)
			if err != nil {
				sendError(w, http.StatusInternalServerError, "internal_error", "Error checking permissions")
				return
			}

			if !hasPermission {
				sendError(w, http.StatusForbidden, "forbidden", "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireRole - Native Go version
func RequireRole(rbacService services.RBACService, roleName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := GetUserIDFromContext(r)
			if userID == 0 {
				sendError(w, http.StatusUnauthorized, "unauthorized", "User not authenticated")
				return
			}

			hasRole, err := rbacService.HasRole(userID, roleName)
			if err != nil {
				sendError(w, http.StatusInternalServerError, "internal_error", "Error checking role")
				return
			}

			if !hasRole {
				sendError(w, http.StatusForbidden, "forbidden", "Insufficient role privileges")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireRoleAny - Native Go version
func RequireRoleAny(rbacService services.RBACService, roleNames ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := GetUserIDFromContext(r)
			if userID == 0 {
				sendError(w, http.StatusUnauthorized, "unauthorized", "User not authenticated")
				return
			}

			userRole, err := rbacService.GetUserRole(userID)
			if err != nil {
				sendError(w, http.StatusInternalServerError, "internal_error", "Error getting user role")
				return
			}

			for _, roleName := range roleNames {
				if userRole.Name == roleName {
					next.ServeHTTP(w, r)
					return
				}
			}

			sendError(w, http.StatusForbidden, "forbidden", "Insufficient role privileges")
		})
	}
}

// SelfOrPermission - Native Go version
func SelfOrPermission(rbacService services.RBACService, resource, action string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := GetUserIDFromContext(r)
			if userID == 0 {
				sendError(w, http.StatusUnauthorized, "unauthorized", "User not authenticated")
				return
			}

			// Extract ID from URL path
			pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
			var targetUserID string
			if len(pathParts) > 0 {
				targetUserID = pathParts[len(pathParts)-1]
			}

			if targetUserID == "" {
				sendError(w, http.StatusBadRequest, "bad_request", "User ID parameter is required")
				return
			}

			// Check if user is accessing their own resource
			if targetUserID == fmt.Sprintf("%d", userID) {
				next.ServeHTTP(w, r)
				return
			}

			// Check permission
			hasPermission, err := rbacService.CheckPermission(userID, resource, action)
			if err != nil {
				sendError(w, http.StatusInternalServerError, "internal_error", "Error checking permissions")
				return
			}

			if !hasPermission {
				sendError(w, http.StatusForbidden, "forbidden", "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Chain - Helper function to chain multiple middlewares
func Chain(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// ChainFunc - Helper function to chain multiple middlewares for HandlerFunc
func ChainFunc(handlerFunc http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) http.Handler {
	return Chain(handlerFunc, middlewares...)
}
