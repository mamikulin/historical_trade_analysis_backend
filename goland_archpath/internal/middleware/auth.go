package middleware

import (
	"archpath/internal/app/session"
	"context"
	"net/http"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
	RoleKey   contextKey = "role"
)

func AuthMiddleware(sessionManager *session.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_id")
			if err != nil {
				// Нет cookie - просто пропускаем дальше (для публичных эндпоинтов)
				next.ServeHTTP(w, r)
				return
			}

			sessionData, err := sessionManager.GetSession(r.Context(), cookie.Value)
			if err != nil {
				// Невалидная сессия - пропускаем дальше
				next.ServeHTTP(w, r)
				return
			}

			// Добавляем данные пользователя в контекст
			ctx := context.WithValue(r.Context(), UserIDKey, sessionData.UserID)
			ctx = context.WithValue(ctx, RoleKey, sessionData.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

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

func GetUserIDFromContext(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value(UserIDKey).(uint)
	return userID, ok
}

func GetRoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(RoleKey).(string)
	return role, ok
}