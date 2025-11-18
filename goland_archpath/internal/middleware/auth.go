package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
	RoleKey   contextKey = "role"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a new JWT token for a user
func GenerateJWT(userID uint, role string, secret string, expiry time.Duration) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateJWT validates and parses a JWT token
func ValidateJWT(tokenString string, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// AuthMiddleware extracts and validates JWT from Authorization header
func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				// No token - просто пропускаем дальше (для публичных эндпоинтов)
				next.ServeHTTP(w, r)
				return
			}

			// Ожидаем формат: "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				next.ServeHTTP(w, r)
				return
			}

			tokenString := parts[1]
			claims, err := ValidateJWT(tokenString, jwtSecret)
			if err != nil {
				// Невалидный токен - пропускаем дальше
				next.ServeHTTP(w, r)
				return
			}

			// Добавляем данные пользователя в контекст
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, RoleKey, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAuth ensures the request has a valid JWT token
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(UserIDKey)
		if userID == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequireModerator ensures the user has moderator role
func RequireModerator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(UserIDKey)
		if userID == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		role := r.Context().Value(RoleKey)
		if role == nil || role.(string) != "moderator" {
			http.Error(w, "Forbidden: moderator access required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// GetUserIDFromContext retrieves user ID from context
func GetUserIDFromContext(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value(UserIDKey).(uint)
	return userID, ok
}

// GetRoleFromContext retrieves user role from context
func GetRoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(RoleKey).(string)
	return role, ok
}