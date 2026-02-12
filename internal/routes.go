package internal

import (
	"io/fs"
	"net/http"

	"github.com/synic/blog/internal/controller"
	"github.com/synic/blog/internal/view"
)

func RegisterRoutes(
	handler *http.ServeMux,
	assets fs.FS,
	articleController controller.ArticleController,
	commentController controller.CommentController,
	authController controller.AuthController,
	leaderboardController controller.LeaderboardController,
) {
	// static files
	handler.Handle("GET /static/", StaticHandler(assets))

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
	handler.HandleFunc("GET /article/{date}/{slug}/comments", commentController.List)
	handler.HandleFunc("POST /article/{date}/{slug}/comments", commentController.Create)

	// admin
	handler.HandleFunc("GET /admin/comments/{id}/approve", commentController.Approve)
	handler.HandleFunc("GET /admin/comments/{id}/delete", commentController.Delete)

	// auth
	handler.HandleFunc("GET /auth/login", authController.Login)
	handler.HandleFunc("GET /auth/callback", authController.Callback)
	handler.HandleFunc("POST /auth/logout", authController.Logout)
	handler.HandleFunc("GET /unsubscribe", authController.Unsubscribe)

	// leaderboard
	handler.HandleFunc("GET /leaderboard", leaderboardController.Show)

	// errors
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		view.Error(w, r, nil, http.StatusNotFound, "Not Found", "That's a 404.")
	})
}
