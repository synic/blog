package middleware

import (
	"context"
	"net/http"
)

var DatastarPartialContextKey = "isDatastarPartial"

func DatastarMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(
				r.Context(),
				DatastarPartialContextKey,
				r.Header.Get("Datastar-Request") == "true",
			)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
