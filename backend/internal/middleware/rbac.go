package middleware

import (
	"e-commerce/backend/internal/services"
	"fmt"
	"net/http"
	"strings"
)

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
