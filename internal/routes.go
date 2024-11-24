package internal

import (
	"io/fs"
	"net/http"

	"github.com/synic/adamthings.me/internal/controller"
	"github.com/synic/adamthings.me/internal/view"
)

func RegisterRoutes(
	handler *http.ServeMux,
	assets fs.FS,
	articleController controller.ArticleController,
) {
	// articles
	handler.HandleFunc("/{$}", articleController.Index)
	handler.HandleFunc("/articles/{date}/{slug}", articleController.Article)
	handler.HandleFunc("/archive", articleController.Archive)

	// static files
	handler.Handle("GET /static/", StaticHandler(assets))

	// errors
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		view.Error(w, r, nil, http.StatusNotFound, "Not Found", "That's a 404.")
	})
}
