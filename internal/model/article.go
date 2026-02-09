package model

import (
	"fmt"
	"time"

	"github.com/a-h/templ"
)

type Article struct {
	PublishedAt   time.Time         `json:"publishedAt"     yaml:"publishedAt"`
	Extra         map[string]string `json:"extra,omitempty" yaml:"extra,omitempty"`
	Body          string            `json:"body"            yaml:"body"`
	Summary       string            `json:"summary"         yaml:"summary"`
	Title         string            `json:"title"           yaml:"title"`
	Slug          string            `json:"slug"            yaml:"slug"`
	Tags          []string          `json:"tags"            yaml:"tags"`
	IsPublished   bool              `json:"isPublished"     yaml:"isPublished"`
	OpenGraphData OpenGraphData     `json:"openGraph"       yaml:"openGraph"`
}

func (a *Article) URL() string {
	return fmt.Sprintf(
		"/article/%d-%02d-%02d/%s",
		a.PublishedAt.Year(),
		a.PublishedAt.Month(),
		a.PublishedAt.Day(),
		a.Slug,
	)
}

func (a *Article) SafeURL() templ.SafeURL {
	return templ.URL(a.URL())
}

type ArticleListResponse struct {
	Search     string
	Tag        string
	Items      []*Article
	TotalPages int
	Page       int
	PerPage    int
}
