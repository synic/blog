package model

import (
	"fmt"
	"time"

	"github.com/a-h/templ"
)

type Article struct {
	PublishedAt time.Time         `json:"publishedAt"`
	Extra       map[string]string `json:"extra,omitempty"`
	Body        string            `json:"body"`
	Summary     string            `json:"summary"`
	Title       string            `json:"title"`
	Slug        string            `json:"slug"`
	Tags        []string          `json:"tags"`
	IsPublished bool              `json:"isPublished"`
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

func (a *Article) OpenGraphData() OpenGraphData {
	og := OpenGraphData{Type: "article"}
	image, _ := a.Extra["ogImage"]

	og.Image = image

	title, _ := a.Extra["ogTitle"]

	if title == "" {
		title = a.Title
	}

	og.Title = title

	description, _ := a.Extra["ogDescription"]
	og.Description = description

	return og
}
