package middleware

import (
	"context"
	"net/http"
)

var HtmxPartialContextKey = "ao:ishtmxpartial"

func HtmxMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(
				r.Context(),
				HtmxPartialContextKey,
				r.Header.Get("HX-Request") == "true",
			)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
