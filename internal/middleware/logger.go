package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

type loggerMiddlwareConfig struct {
	logger *log.Logger
}

func defaultLoggerMiddlewareConfig() loggerMiddlwareConfig {
	return loggerMiddlwareConfig{logger: log.Default()}
}

func WithLogger(logger *log.Logger) func(*loggerMiddlwareConfig) {
	return func(conf *loggerMiddlwareConfig) {
		conf.logger = logger
	}
}

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true

	return
}

func (rw *responseWriter) Flush() {
	if f, ok := rw.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (rw *responseWriter) Unwrap() http.ResponseWriter {
	return rw.ResponseWriter
}

func LoggerMiddleware(options ...func(*loggerMiddlwareConfig)) func(http.Handler) http.Handler {
	conf := defaultLoggerMiddlewareConfig()

	for _, option := range options {
		option(&conf)
	}

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					conf.logger.Printf("[ERROR] - %s %s", err, debug.Stack())
				}
			}()

			start := time.Now()
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)

			conf.logger.Printf(
				"[REQ] \"%s %s\" - \"%s\" - %d %s",
				r.Method,
				r.URL.EscapedPath(),
				r.Header.Get("User-Agent"),
				wrapped.status,
				time.Since(start),
			)
		}

		return http.HandlerFunc(fn)
	}
}
