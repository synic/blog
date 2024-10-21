package middleware

import (
	"context"
	"net/http"
)

func AddContextData(data map[string]any) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			for key, value := range data {
				ctx = context.WithValue(ctx, key, value)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
