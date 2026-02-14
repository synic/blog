package internal

import (
	"io/fs"
	"net/http"

	"github.com/synic/blog/internal/controller"
	"github.com/synic/blog/internal/middleware"
	"github.com/synic/blog/internal/static"
	"github.com/synic/blog/internal/view"
)

func RegisterRoutes(
	handler *http.ServeMux,
	assets fs.FS,
	auth middleware.Middleware,
	csrf middleware.Middleware,
	articleController controller.ArticleController,
	commentController controller.CommentController,
	authController controller.AuthController,
	leaderboardController controller.LeaderboardController,
) {
	// static files
	handler.Handle("GET /static/", static.StaticHandler(assets))

	// articles
	handler.HandleFunc("/{$}", articleController.Index)
	handler.HandleFunc("/article/create", articleController.Create)
	handler.HandleFunc("GET /article/{date}/{slug}", articleController.Article)
	handler.HandleFunc("/archive", articleController.Archive)
	handler.HandleFunc(
		"/articles/{date}/{slug}",
		func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(
				w,
				r,
				"/article/"+r.PathValue("date")+"/"+r.PathValue("slug"),
				http.StatusMovedPermanently,
			)
		},
	)
	handler.HandleFunc("/feed.xml", articleController.Feed)

	// comments
	handler.Handle(
		"GET /article/{date}/{slug}/comments",
		middleware.With(commentController.List, auth, csrf),
	)
	handler.Handle(
		"POST /article/{date}/{slug}/comments",
		middleware.With(commentController.Create, auth, csrf),
	)

	// admin
	handler.Handle(
		"GET /admin/comments/{id}/approve",
		middleware.With(commentController.Approve, auth),
	)
	handler.Handle(
		"GET /admin/comments/{id}/delete",
		middleware.With(commentController.Delete, auth),
	)

	// auth
	handler.Handle("GET /auth/login", middleware.With(authController.Login, auth))
	handler.Handle("GET /auth/callback", middleware.With(authController.Callback, auth))
	handler.Handle("POST /auth/logout", middleware.With(authController.Logout, auth, csrf))
	handler.HandleFunc("GET /unsubscribe", authController.Unsubscribe)

	// leaderboard
	handler.Handle("GET /leaderboard", middleware.With(leaderboardController.Show, auth))

	// errors
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		view.Error(w, r, nil, http.StatusNotFound, "Not Found", "That's a 404.")
	})
}
