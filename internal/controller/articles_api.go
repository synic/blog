package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/synic/blog/internal/store"
	"github.com/synic/blog/internal/view"
)

type ArticleApiController struct {
	repo store.ArticleRepository
	articleControllerConfig
}

func NewArticleApiController(
	repo store.ArticleRepository,
	options ...func(*articleControllerConfig),
) ArticleApiController {
	conf := defaultArticleControllerConfig()

	for _, option := range options {
		option(&conf)
	}

	return ArticleApiController{repo: repo, articleControllerConfig: conf}
}

func (h ArticleApiController) Index(w http.ResponseWriter, r *http.Request) {
	res, err := buildAndPaginateArticleList(h.repo, h.articleControllerConfig, r)

	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			view.Error(w, r, err, 404, "Not Found", "Sorry, no articles could be found.")
		} else {
			view.Error(w, r, err, 500, "Internal Server Error", "An error occurred while retrieving articles.")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
