package middleware

import (
	"bookify/pkg/auth"
	"context"
	"net/http"
	"strings"
)


func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid config header format", http.StatusUnauthorized)
			return
		}

		
		claims, err := auth.VerifyToken(parts[1])
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "role", claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}


func RoleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			
			userRole, ok := r.Context().Value("role").(string)
			if !ok {
				http.Error(w, "Unauthorized: Role not found", http.StatusUnauthorized)
				return
			}

			
			authorized := false
			for _, role := range allowedRoles {
				if userRole == role {
					authorized = true
					break
				}
			}

			if !authorized {
				http.Error(w, "Forbidden: You don't have permission", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}