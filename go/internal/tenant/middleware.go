package tenant

import (
	"context"
	"net/http"
	"strings"
)

// ContextKey type for context keys
type ContextKey string

const (
	// TenantIDKey is the context key for tenant ID
	TenantIDKey ContextKey = "tenantId"
	// UserIDKey is the context key for user ID
	UserIDKey ContextKey = "userId"
)

// Middleware provides tenant isolation for HTTP requests
type Middleware struct {
	manager *Manager
}

// NewMiddleware creates a new tenant middleware
func NewMiddleware(manager *Manager) *Middleware {
	return &Middleware{
		manager: manager,
	}
}

// TenantIsolation middleware extracts tenant information from requests
func (m *Middleware) TenantIsolation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract tenant ID from header, query parameter, or subdomain
		tenantID := m.extractTenantID(r)
		userID := m.extractUserID(r)

		// Validate tenant access if tenant ID is provided
		if tenantID != "" {
			if _, err := m.manager.GetTenant(tenantID); err != nil {
				http.Error(w, "Invalid tenant", http.StatusForbidden)
				return
			}

			// Validate user access if both tenant and user are provided
			if userID != "" {
				if err := m.manager.ValidateAccess(userID, tenantID); err != nil {
					http.Error(w, "Access denied", http.StatusForbidden)
					return
				}
			}
		}

		// Add tenant and user info to context
		ctx := r.Context()
		if tenantID != "" {
			ctx = context.WithValue(ctx, TenantIDKey, tenantID)
		}
		if userID != "" {
			ctx = context.WithValue(ctx, UserIDKey, userID)
		}

		// Continue with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// extractTenantID extracts tenant ID from various sources
func (m *Middleware) extractTenantID(r *http.Request) string {
	// 1. Check X-Tenant-ID header
	if tenantID := r.Header.Get("X-Tenant-ID"); tenantID != "" {
		return tenantID
	}

	// 2. Check query parameter
	if tenantID := r.URL.Query().Get("tenant"); tenantID != "" {
		return tenantID
	}

	// 3. Check subdomain (e.g., research-lab.awsresearch.com)
	host := r.Host
	if strings.Contains(host, ".") {
		parts := strings.Split(host, ".")
		if len(parts) > 2 && parts[0] != "www" {
			return parts[0]
		}
	}

	// 4. Check URL path prefix (e.g., /tenant/{tenantId}/...)
	path := r.URL.Path
	if strings.HasPrefix(path, "/tenant/") {
		parts := strings.Split(strings.TrimPrefix(path, "/tenant/"), "/")
		if len(parts) > 0 && parts[0] != "" {
			return parts[0]
		}
	}

	return ""
}

// extractUserID extracts user ID from request headers or session
func (m *Middleware) extractUserID(r *http.Request) string {
	// 1. Check X-User-ID header
	if userID := r.Header.Get("X-User-ID"); userID != "" {
		return userID
	}

	// 2. Check Authorization header (extract from JWT token)
	if auth := r.Header.Get("Authorization"); auth != "" {
		// In a real implementation, this would decode JWT token
		// For demo purposes, we'll use a simple format
		if strings.HasPrefix(auth, "Bearer ") {
			token := strings.TrimPrefix(auth, "Bearer ")
			// Extract user ID from token (simplified)
			if strings.HasPrefix(token, "user-") {
				return token
			}
		}
	}

	// 3. Check session cookie (simplified)
	if cookie, err := r.Cookie("session-user-id"); err == nil {
		return cookie.Value
	}

	return ""
}

// GetTenantID retrieves tenant ID from request context
func GetTenantID(ctx context.Context) string {
	if tenantID, ok := ctx.Value(TenantIDKey).(string); ok {
		return tenantID
	}
	return ""
}

// GetUserID retrieves user ID from request context
func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}

// RequireTenant middleware ensures a valid tenant is present
func (m *Middleware) RequireTenant(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantID := GetTenantID(r.Context())
		if tenantID == "" {
			http.Error(w, "Tenant required", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RequireUser middleware ensures a valid user is present
func (m *Middleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserID(r.Context())
		if userID == "" {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RequirePermission middleware checks if user has specific permission
func (m *Middleware) RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := GetUserID(r.Context())
			if userID == "" {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			user, err := m.manager.GetUser(userID)
			if err != nil {
				http.Error(w, "User not found", http.StatusForbidden)
				return
			}

			// Check if user has the required permission
			hasPermission := false
			for _, perm := range user.Permissions {
				if perm == permission || perm == "*" {
					hasPermission = true
					break
				}
			}

			if !hasPermission {
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}