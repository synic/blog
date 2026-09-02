package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type cacheControlConfig struct {
	maxAge         int
	public         bool
	private        bool
	immutable      bool
	mustRevalidate bool
	noCache        bool
	noStore        bool
	custom         string
}

func defaultCacheControlConfig() cacheControlConfig {
	return cacheControlConfig{
		maxAge:    31536000, // 1 year (in seconds)
		public:    true,
		immutable: true,
	}
}

func (c *cacheControlConfig) buildHeaderValue() string {
	if c.custom != "" {
		return c.custom
	}

	var parts []string

	if c.noStore {
		parts = append(parts, "no-store")
	}
	if c.noCache {
		parts = append(parts, "no-cache")
	}
	if c.public {
		parts = append(parts, "public")
	} else if c.private {
		parts = append(parts, "private")
	}

	if c.maxAge >= 0 {
		parts = append(parts, fmt.Sprintf("max-age=%d", c.maxAge))
	}

	if c.mustRevalidate {
		parts = append(parts, "must-revalidate")
	}

	if c.immutable {
		parts = append(parts, "immutable")
	}

	return strings.Join(parts, ", ")
}

type CacheControlOption func(*cacheControlConfig)

// WithMaxAge sets the max-age directive using a time.Duration.
func WithMaxAge(d time.Duration) CacheControlOption {
	return func(conf *cacheControlConfig) {
		conf.maxAge = int(d.Seconds())
	}
}

// WithMaxAgeSeconds sets the max-age directive using seconds.
func WithMaxAgeSeconds(seconds int) CacheControlOption {
	return func(conf *cacheControlConfig) {
		conf.maxAge = seconds
	}
}

// WithPublic sets or unsets the public directive.
func WithPublic(public bool) CacheControlOption {
	return func(conf *cacheControlConfig) {
		conf.public = public
		if public {
			conf.private = false
		}
	}
}

// WithPrivate sets or unsets the private directive.
func WithPrivate(private bool) CacheControlOption {
	return func(conf *cacheControlConfig) {
		conf.private = private
		if private {
			conf.public = false
		}
	}
}

// WithImmutable sets or unsets the immutable directive.
func WithImmutable(immutable bool) CacheControlOption {
	return func(conf *cacheControlConfig) {
		conf.immutable = immutable
	}
}

// WithMustRevalidate sets or unsets the must-revalidate directive.
func WithMustRevalidate(mustRevalidate bool) CacheControlOption {
	return func(conf *cacheControlConfig) {
		conf.mustRevalidate = mustRevalidate
	}
}

// WithNoCache configures Cache-Control: no-cache.
func WithNoCache() CacheControlOption {
	return func(conf *cacheControlConfig) {
		conf.noCache = true
		conf.public = false
		conf.immutable = false
		conf.maxAge = -1
	}
}

// WithNoStore configures Cache-Control: no-store.
func WithNoStore() CacheControlOption {
	return func(conf *cacheControlConfig) {
		conf.noStore = true
		conf.public = false
		conf.immutable = false
		conf.maxAge = -1
	}
}

// WithHeaderValue allows setting an explicit, custom Cache-Control header value.
func WithHeaderValue(value string) CacheControlOption {
	return func(conf *cacheControlConfig) {
		conf.custom = value
	}
}

// CacheControlMiddleware returns a middleware that sets the Cache-Control header.
// By default, it sets: "public, max-age=31536000, immutable" (1 year).
func CacheControlMiddleware(options ...CacheControlOption) Middleware {
	conf := defaultCacheControlConfig()
	for _, opt := range options {
		opt(&conf)
	}

	headerValue := conf.buildHeaderValue()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if headerValue != "" {
				w.Header().Set("Cache-Control", headerValue)
			}
			next.ServeHTTP(w, r)
		})
	}
}

// CacheMiddleware is an alias for CacheControlMiddleware.
var CacheMiddleware = CacheControlMiddleware
