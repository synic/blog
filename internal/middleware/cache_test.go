package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/synic/blog/internal/middleware"
)

func TestCacheControlMiddleware_Default(t *testing.T) {
	mw := middleware.CacheControlMiddleware()

	called := false
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/static/css/main.css", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if !called {
		t.Fatal("expected inner handler to be called")
	}

	expected := "public, max-age=31536000, immutable"
	actual := rec.Header().Get("Cache-Control")
	if actual != expected {
		t.Errorf("expected Cache-Control %q, got %q", expected, actual)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestCacheControlMiddleware_CustomMaxAgeDuration(t *testing.T) {
	mw := middleware.CacheControlMiddleware(
		middleware.WithMaxAge(24*time.Hour),
		middleware.WithImmutable(false),
	)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/static/css/main.css", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	expected := "public, max-age=86400"
	actual := rec.Header().Get("Cache-Control")
	if actual != expected {
		t.Errorf("expected Cache-Control %q, got %q", expected, actual)
	}
}

func TestCacheControlMiddleware_CustomMaxAgeSeconds(t *testing.T) {
	mw := middleware.CacheControlMiddleware(
		middleware.WithMaxAgeSeconds(3600),
		middleware.WithImmutable(false),
	)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/static/test.js", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	expected := "public, max-age=3600"
	actual := rec.Header().Get("Cache-Control")
	if actual != expected {
		t.Errorf("expected Cache-Control %q, got %q", expected, actual)
	}
}

func TestCacheControlMiddleware_PrivateAndMustRevalidate(t *testing.T) {
	mw := middleware.CacheControlMiddleware(
		middleware.WithPrivate(true),
		middleware.WithMaxAgeSeconds(300),
		middleware.WithMustRevalidate(true),
		middleware.WithImmutable(false),
	)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/static/private.json", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	expected := "private, max-age=300, must-revalidate"
	actual := rec.Header().Get("Cache-Control")
	if actual != expected {
		t.Errorf("expected Cache-Control %q, got %q", expected, actual)
	}
}

func TestCacheControlMiddleware_NoCache(t *testing.T) {
	mw := middleware.CacheControlMiddleware(
		middleware.WithNoCache(),
	)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/static/no-cache", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	expected := "no-cache"
	actual := rec.Header().Get("Cache-Control")
	if actual != expected {
		t.Errorf("expected Cache-Control %q, got %q", expected, actual)
	}
}

func TestCacheControlMiddleware_NoStore(t *testing.T) {
	mw := middleware.CacheControlMiddleware(
		middleware.WithNoStore(),
	)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/static/no-store", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	expected := "no-store"
	actual := rec.Header().Get("Cache-Control")
	if actual != expected {
		t.Errorf("expected Cache-Control %q, got %q", expected, actual)
	}
}

func TestCacheControlMiddleware_CustomHeaderValue(t *testing.T) {
	mw := middleware.CacheControlMiddleware(
		middleware.WithHeaderValue("public, max-age=600, s-maxage=1200"),
	)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/static/custom", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	expected := "public, max-age=600, s-maxage=1200"
	actual := rec.Header().Get("Cache-Control")
	if actual != expected {
		t.Errorf("expected Cache-Control %q, got %q", expected, actual)
	}
}

func TestCacheMiddlewareAlias(t *testing.T) {
	mw := middleware.CacheMiddleware()

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/static/alias", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	expected := "public, max-age=31536000, immutable"
	actual := rec.Header().Get("Cache-Control")
	if actual != expected {
		t.Errorf("expected Cache-Control %q, got %q", expected, actual)
	}
}
