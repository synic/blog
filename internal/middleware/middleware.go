package middleware

import (
	"net/http"
)

type Middleware = func(http.Handler) http.Handler

func Wrap(handler http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	wrapped := handler
	for _, h := range middleware {
		wrapped = h(wrapped)
	}
	return wrapped
}

func With(fn http.HandlerFunc, mws ...Middleware) http.Handler {
	var h http.Handler = fn
	for _, mw := range mws {
		h = mw(h)
	}
	return h
}
