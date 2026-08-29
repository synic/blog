package middleware

import (
	"context"
	"net/http"
)

var HtmxPartialContextKey = "isHtmxPartial"

func HtmxMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			isPartial := r.Header.Get("HX-Request") == "true" &&
				r.Header.Get("HX-Request-Type") != "full"

			ctx := context.WithValue(r.Context(), HtmxPartialContextKey, isPartial)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
