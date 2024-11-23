package internal

import (
	"io/fs"
	"net/http"

	"github.com/synic/adamthings.me/internal/controller"
)

func RegisterRoutes(
	handler *http.ServeMux,
	assets fs.FS,
	articleController controller.ArticleController,
) {
	handler.Handle("GET /static/", StaticHandler(assets))

	// articles
	handler.HandleFunc("/{$}", articleController.Index)
	handler.HandleFunc("/articles/{date}/{slug}", articleController.Article)
	handler.HandleFunc("/archive", articleController.Archive)
}
