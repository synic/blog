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
	articleApiController controller.ArticleApiController,
) {
	// static files
	handler.Handle("GET /static/", StaticHandler(assets))

	// articles
	handler.HandleFunc("/{$}", articleController.Index)
	handler.HandleFunc("/article/{date}/{slug}", articleController.Article)
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

	// api
	handler.HandleFunc("/api/v1/articles", articleApiController.Index)

	// feed
	handler.HandleFunc("/feed.xml", articleController.Feed)

	// errors
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		view.Error(w, r, nil, http.StatusNotFound, "Not Found", "That's a 404.")
	})
}
