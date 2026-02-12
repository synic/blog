package middleware

import (
	"context"
	"net/http"

	"github.com/synic/blog/internal/db"
	"github.com/synic/blog/internal/model"
)

type contextKey string

const userContextKey contextKey = "user"

func UserFromContext(ctx context.Context) *model.User {
	if u, ok := ctx.Value(userContextKey).(*model.User); ok {
		return u
	}
	return nil
}

func AuthMiddleware(queries *db.Queries, adminEmail string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_token")
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			session, err := queries.GetSessionByToken(r.Context(), cookie.Value)
			if err != nil {
				// Clear stale cookie
				http.SetCookie(w, &http.Cookie{
					Name:     "session_token",
					Value:    "",
					Path:     "/",
					MaxAge:   -1,
					HttpOnly: true,
				})
				next.ServeHTTP(w, r)
				return
			}

			user := &model.User{
				ID:        session.UserID,
				Username:  session.Username,
				AvatarURL: session.AvatarUrl,
				Email:     session.Email,
				IsAdmin:   session.Email != "" && session.Email == adminEmail,
			}

			ctx := context.WithValue(r.Context(), userContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
