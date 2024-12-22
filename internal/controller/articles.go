package controller

import (
	"errors"
	"math"
	"net/http"
	"strconv"

	"github.com/synic/blog/internal/model"
	"github.com/synic/blog/internal/store"
	"github.com/synic/blog/internal/view"
)

type articleControllerConfig struct {
	articlesPerPage    int
	maxArticlesPerPage int
}

type ArticleController struct {
	repo store.ArticleRepository
	articleControllerConfig
}

func defaultArticleControllerConfig() articleControllerConfig {
	return articleControllerConfig{
		articlesPerPage:    20,
		maxArticlesPerPage: 50,
	}
}

func WithPagination(perPage, maxPerPage int) func(*articleControllerConfig) {
	return func(conf *articleControllerConfig) {
		conf.articlesPerPage = perPage
		conf.maxArticlesPerPage = maxPerPage
	}
}

func NewArticleController(
	repo store.ArticleRepository,
	options ...func(*articleControllerConfig),
) ArticleController {
	conf := defaultArticleControllerConfig()

	for _, option := range options {
		option(&conf)
	}

	return ArticleController{repo: repo, articleControllerConfig: conf}
}

func (h ArticleController) Index(w http.ResponseWriter, r *http.Request) {
	var err error = nil
	articles := []*model.Article{}

	search := r.FormValue("search")

	if search == "" {
		search = r.URL.Query().Get("search")
	}

	tag := r.URL.Query().Get("tag")

	if tag != "" {
		articles, err = h.repo.FindByTag(r.Context(), tag)
	} else if search != "" {
		articles, err = h.repo.Search(r.Context(), search)
	} else {
		articles, err = h.repo.FindAll(r.Context())
	}

	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			view.Error(w, r, err, 404, "Not Found", "Sorry, no articles could be found.")
		} else {
			view.Error(w, r, err, 500, "Internal Server Error", "An error occurred while retrieving articles.")
		}
		return
	}

	h.renderAndPageArticles(w, r, articles, search, tag)
}

func (h ArticleController) Article(w http.ResponseWriter, r *http.Request) {
	var slug string

	slug = r.PathValue("slug")

	if slug == "" {
		view.Error(w, r, nil, 404, "Not Found", "Invalid article path")
		return
	}

	article, err := h.repo.FindOneBySlug(r.Context(), slug)

	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			view.Error(w, r, err, 404, "Not Found", "Sorry, that article could not be found.")
		} else {
			view.Error(w, r, err, 500, "Internal Server Error", "An error occurred while retrieving the article.")
		}
		return
	}

	view.Render(w, r, view.ArticleView(article))
}

func (h ArticleController) renderAndPageArticles(
	w http.ResponseWriter,
	r *http.Request,
	articles []*model.Article,
	search, tag string,
) {
	page := 0
	perPage := h.articlesPerPage

	i, err := strconv.Atoi(r.URL.Query().Get("page"))

	if err == nil {
		page = i - 1
	}

	i, err = strconv.Atoi(r.URL.Query().Get("per_page"))

	if err == nil {
		perPage = i
		if perPage > h.maxArticlesPerPage {
			perPage = h.maxArticlesPerPage
		}
	}

	start := min(max(0, page*perPage), len(articles))
	end := min(max(0, start+perPage), len(articles))

	view.Render(
		w, r,
		view.ArticlesView(
			model.PageData{
				Page:       page + 1,
				PerPage:    perPage,
				Items:      articles[start:end],
				Search:     search,
				Tag:        tag,
				TotalPages: int(math.Ceil(float64(len(articles)) / float64(perPage))),
			},
		),
	)
}

func (h ArticleController) Archive(w http.ResponseWriter, r *http.Request) {
	articles, _ := h.repo.FindAll(r.Context())
	view.Render(w, r, view.ArchiveView(len(articles), h.repo.TagInfo(r.Context())))
}
