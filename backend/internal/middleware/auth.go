package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"incus-manager/internal/service"
)

type ContextKey string

const UserIDContextKey ContextKey = "userID"
const UsernameContextKey ContextKey = "username"
const RoleContextKey ContextKey = "role"

func Authenticate(authService *service.AuthService) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := authService.ValidateToken(tokenString)
			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			userID := uint(claims["user_id"].(float64))
			username := claims["username"].(string)
			role := claims["role"].(string)

			ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
			ctx = context.WithValue(ctx, UsernameContextKey, username)
			ctx = context.WithValue(ctx, RoleContextKey, role)

			next(w, r.WithContext(ctx))
		}
	}
}

func GetUserID(r *http.Request) uint {
	if userID, ok := r.Context().Value(UserIDContextKey).(uint); ok {
		return userID
	}
	return 0
}

func GetRole(r *http.Request) string {
	if role, ok := r.Context().Value(RoleContextKey).(string); ok {
		return role
	}
	return ""
}
