package controller

import (
	"net/http"
	"time"

	"github.com/gorilla/feeds"
	"github.com/synic/blog/internal/store"
)

type ArticleRSSController interface {
	Feed(w http.ResponseWriter, r *http.Request)
}

type articleRssController struct {
	repo store.ArticleRepository
}

func NewArticleRSSController(repo store.ArticleRepository) ArticleRSSController {
	return &articleRssController{repo: repo}
}

func (c *articleRssController) Feed(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	articles, err := c.repo.FindAll(ctx)
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
			Link:        &feeds.Link{Href: "https://synic.dev" + article.URL()},
			Description: article.Summary,
			Created:     article.PublishedAt,
		}
		feedItems = append(feedItems, item)
	}
	feed.Items = feedItems

	w.Header().Set("Content-Type", "application/rss+xml")
	feed.WriteRss(w)
}
