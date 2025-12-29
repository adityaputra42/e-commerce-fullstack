package middleware

import (
	"e-commerce/backend/internal/services"
	"net/http"
	"strconv"
	"strings"
)

// =======================
// REQUIRE PERMISSION
// =======================

func RequirePermission(
	rbac services.RBACService,
	resource string,
	action string,
) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			userID := GetUserIDFromContext(r)
			if userID == 0 {
				sendError(w, http.StatusUnauthorized, "unauthorized", "Not authenticated")
				return
			}

			ok, err := rbac.CheckPermission(userID, resource, action)
			if err != nil || !ok {
				sendError(w, http.StatusForbidden, "forbidden", "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// =======================
// REQUIRE ROLE (HIERARCHY)
// =======================

func RequireRole(
	rbac services.RBACService,
	roleName string,
) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			userID := GetUserIDFromContext(r)
			if userID == 0 {
				sendError(w, http.StatusUnauthorized, "unauthorized", "Not authenticated")
				return
			}

			ok, err := rbac.HasRole(userID, roleName)
			if err != nil || !ok {
				sendError(w, http.StatusForbidden, "forbidden", "Insufficient role privileges")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// =======================
// PERMISSION OR OWN
// =======================

func RequirePermissionOrOwn(
	rbac services.RBACService,
	resource string,
	action string,
) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			userID := GetUserIDFromContext(r)
			if userID == 0 {
				sendError(w, http.StatusUnauthorized, "unauthorized", "Not authenticated")
				return
			}

			parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
			idStr := parts[len(parts)-1]

			resourceID, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				sendError(w, http.StatusBadRequest, "bad_request", "Invalid resource ID")
				return
			}

			ok, err := rbac.CheckPermissionOrOwn(
				userID,
				resource,
				action,
				uint(resourceID),
			)

			if err != nil || !ok {
				sendError(w, http.StatusForbidden, "forbidden", "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
