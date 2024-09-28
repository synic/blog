package route

import (
	"log"
	"net/http"
	"strconv"

	"github.com/synic/adamthings.me/internal/model"
	"github.com/synic/adamthings.me/internal/store"
	"github.com/synic/adamthings.me/internal/web/component"
	"github.com/synic/adamthings.me/internal/web/pagination"
	"github.com/synic/adamthings.me/internal/web/render"
)

type articleRouterConfig struct {
	articlesPerPage    int
	maxArticlesPerPage int
}

type articleRouter struct {
	repo *store.ArticleRepository
	articleRouterConfig
}

func getDefaultArticleRouterConfig() articleRouterConfig {
	return articleRouterConfig{
		articlesPerPage:    30,
		maxArticlesPerPage: 50,
	}
}

func WithPagination(perPage, maxPerPage int) func(*articleRouterConfig) {
	return func(conf *articleRouterConfig) {
		conf.articlesPerPage = perPage
		conf.maxArticlesPerPage = maxPerPage
	}
}

func NewArticleRouter(
	repo *store.ArticleRepository,
	options ...func(*articleRouterConfig),
) articleRouter {
	conf := getDefaultArticleRouterConfig()

	for _, option := range options {
		option(&conf)
	}

	return articleRouter{repo: repo, articleRouterConfig: conf}
}

func (h articleRouter) index(w http.ResponseWriter, r *http.Request) {
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
		log.Printf("error finding articles: %s", err)
		render.Error(w, r, 404, "Not Found", "Sorry, no articles could be found.")
		return
	}

	h.renderAndPageArticles(w, r, articles, search, tag)
}

func (h articleRouter) article(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	article, err := h.repo.FindOneBySlug(r.Context(), slug)

	if err != nil {
		render.Error(w, r, 404, "Not Found", "Sorry, that article could not be found.")
		return
	}

	render.Templ(w, r, component.ArticleView(article))
}

func (h articleRouter) renderAndPageArticles(
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

	render.Templ(
		w, r,
		component.ArticlesView(
			pagination.PageData{
				Page:       page + 1,
				PerPage:    perPage,
				Items:      articles[start:end],
				Search:     search,
				Tag:        tag,
				TotalPages: len(articles) / perPage,
			},
		),
	)
}

func (h articleRouter) Mount(server *http.ServeMux) {
	server.HandleFunc("/{$}", h.index)
	server.HandleFunc("GET /articles/{date}/{slug}", h.article)
}
