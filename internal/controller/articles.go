package controller

import (
	"errors"
	"github.com/gorilla/feeds"
	"math"
	"net/http"
	"strconv"
	"time"

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

func buildAndPaginateArticleList(repo store.ArticleRepository, conf articleControllerConfig, r *http.Request) (model.ArticleList, error) {
	var err error = nil
	articles := []*model.Article{}
	search := r.FormValue("search")

	if search == "" {
		search = r.URL.Query().Get("search")
	}

	tag := r.URL.Query().Get("tag")

	if tag != "" {
		articles, err = repo.FindByTag(r.Context(), tag)
	} else if search != "" {
		articles, err = repo.Search(r.Context(), search)
	} else {
		articles, err = repo.FindAll(r.Context())
	}

	if err != nil {
		return model.ArticleList{}, err
	}

	page := 0
	perPage := conf.articlesPerPage

	i, err := strconv.Atoi(r.URL.Query().Get("page"))

	if err == nil {
		page = i - 1
	}

	i, err = strconv.Atoi(r.URL.Query().Get("per_page"))

	if err == nil {
		perPage = i
		if perPage > conf.maxArticlesPerPage {
			perPage = conf.maxArticlesPerPage
		}
	}

	start := min(max(0, page*perPage), len(articles))
	end := min(max(0, start+perPage), len(articles))

	return model.ArticleList{
		Search:     search,
		Tag:        tag,
		Items:      articles[start:end],
		TotalPages: int(math.Ceil(float64(len(articles)) / float64(perPage))),
		Page:       page,
		PerPage:    perPage,
	}, nil
}

func (h ArticleController) Index(w http.ResponseWriter, r *http.Request) {
	res, err := buildAndPaginateArticleList(h.repo, h.articleControllerConfig, r)

	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			view.Error(w, r, err, 404, "Not Found", "Sorry, no articles could be found.")
		} else {
			view.Error(w, r, err, 500, "Internal Server Error", "An error occurred while retrieving articles.")
		}
		return
	}

	view.Render(w, r, view.ArticlesView(res))
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

func (h ArticleController) Feed(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	articles, err := h.repo.FindAll(ctx)
	if err != nil {
		http.Error(w, "Failed to generate feed", http.StatusInternalServerError)
		return
	}

	feed := &feeds.Feed{
		Title:       "Adam's Blog",
		Link:        &feeds.Link{Href: "https://synic.dev"},
		Description: "Programming, Vim, Photography, and more!",
		Created:     time.Now(),
	}

	var feedItems []*feeds.Item
	for _, article := range articles {
		item := &feeds.Item{
			Title:       article.Title,
			Link:        &feeds.Link{Href: "https://synic.dev" + article.URL},
			Description: article.Summary,
			Created:     article.PublishedAt,
		}
		feedItems = append(feedItems, item)
	}
	feed.Items = feedItems

	w.Header().Set("Content-Type", "application/rss+xml")
	feed.WriteRss(w)
}

func (h ArticleController) Archive(w http.ResponseWriter, r *http.Request) {
	articles, _ := h.repo.FindAll(r.Context())
	view.Render(w, r, view.ArchiveView(len(articles), h.repo.TagInfo(r.Context())))
}
