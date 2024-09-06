package route

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/synic/adamthings.me/internal/model"
	"github.com/synic/adamthings.me/internal/route/component"
	"github.com/synic/adamthings.me/internal/storage"
	"github.com/synic/adamthings.me/internal/types"
)

type articleRouter struct {
	repo               storage.ArticleRepository
	articlesPerPage    int
	maxArticlesPerPage int
}

func NewArticleRouter(repo storage.ArticleRepository) articleRouter {
	return articleRouter{repo: repo, articlesPerPage: 30, maxArticlesPerPage: 50}
}

func (h articleRouter) index(w http.ResponseWriter, r *http.Request) {
	var err error = nil
	articles := []model.Article{}

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
		fmt.Printf("error finding articles: %s\n", err)
	}

	h.renderAndPageArticles(w, r, articles, search, tag)
}

func (h articleRouter) article(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	article, _ := h.repo.FindOneBySlug(r.Context(), slug)
	RenderTempl(w, r, component.ArticleView(article))
}

func (h articleRouter) renderAndPageArticles(
	w http.ResponseWriter,
	r *http.Request,
	articles []model.Article,
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

	RenderTempl(
		w, r,
		component.ArticlesView(
			types.PageData{
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
	server.HandleFunc("/{$}", h.index) // GET and POST
	server.HandleFunc("GET /articles/{date}/{slug}", h.article)
}
