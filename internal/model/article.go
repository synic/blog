package model

import (
	"fmt"
	"time"

	"github.com/a-h/templ"
)

type Article struct {
	PublishedAt time.Time `json:"published_at"`
	Body        string    `json:"body"`
	Summary     string    `json:"summary"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Tags        []string  `json:"tags"`
	IsPublished bool      `json:"is_published"`
}

func (a *Article) URL() string {
	return fmt.Sprintf(
		"/articles/%d-%02d-%02d/%s",
		a.PublishedAt.Year(),
		a.PublishedAt.Month(),
		a.PublishedAt.Day(),
		a.Slug,
	)
}

func (a *Article) SafeURL() templ.SafeURL {
	return templ.URL(a.URL())
}
