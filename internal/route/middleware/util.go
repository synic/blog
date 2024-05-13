package middleware

import (
	"net/http"
)

func Wrap(handler http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	wrapped := handler
	for _, h := range middleware {
		wrapped = h(wrapped)
	}
	return wrapped
}
