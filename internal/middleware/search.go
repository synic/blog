package middleware

import (
	"context"
	"net/http"
	"strings"
)

var SearchContextKey = "search"

func SearchMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(
				r.Context(),
				SearchContextKey,
				strings.TrimSpace(r.URL.Query().Get("search")),
			)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
