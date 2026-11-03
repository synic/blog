package internal_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/synic/blog/internal"
	"github.com/synic/blog/internal/controller"
	"github.com/synic/blog/internal/middleware"
)

func TestRegisterRoutes_StaticCacheHeaders(t *testing.T) {
	mux := http.NewServeMux()
	staticFS := fstest.MapFS{
		"css/main.css": &fstest.MapFile{
			Data: []byte("body { color: red; }"),
		},
	}

	noopMW := func(next http.Handler) http.Handler {
		return next
	}
	cacheMW := middleware.CacheControlMiddleware()

	internal.RegisterRoutes(
		mux,
		staticFS,
		noopMW,
		noopMW,
		cacheMW,
		controller.ArticleController{},
		controller.CommentController{},
		controller.AuthController{},
		controller.LeaderboardController{},
	)

	req := httptest.NewRequest(http.MethodGet, "/static/css/main.css", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	expectedCacheControl := "public, max-age=31536000, immutable"
	actual := rec.Header().Get("Cache-Control")
	if actual != expectedCacheControl {
		t.Errorf("expected Cache-Control %q, got %q", expectedCacheControl, actual)
	}
}

func TestRegisterRoutes_RobotsTxt(t *testing.T) {
	mux := http.NewServeMux()
	staticFS := fstest.MapFS{
		"robots.txt": &fstest.MapFile{
			Data: []byte("User-agent: *\n"),
		},
	}

	noopMW := func(next http.Handler) http.Handler {
		return next
	}

	internal.RegisterRoutes(
		mux,
		staticFS,
		noopMW,
		noopMW,
		noopMW,
		controller.ArticleController{},
		controller.CommentController{},
		controller.AuthController{},
		controller.LeaderboardController{},
	)

	req := httptest.NewRequest(http.MethodGet, "/robots.txt", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if body := rec.Body.String(); body != "User-agent: *\n" {
		t.Errorf("unexpected body %q", body)
	}
}
