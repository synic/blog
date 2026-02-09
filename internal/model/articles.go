package model

import (
	"encoding/json"
	"fmt"
	"time"
)

type Article struct {
	PublishedAt        time.Time         `json:"publishedAt"     yaml:"publishedAt"`
	PublishedAtDisplay string            `json:"publishedAtDisplay" yaml:"-"`
	Extra              map[string]string `json:"extra,omitempty" yaml:"extra,omitempty"`
	Body               string            `json:"body"            yaml:"body"`
	Summary            string            `json:"summary"         yaml:"summary"`
	Title              string            `json:"title"           yaml:"title"`
	Slug               string            `json:"slug"            yaml:"slug"`
	Tags               []string          `json:"tags"            yaml:"tags"`
	IsPublished        bool              `json:"isPublished"     yaml:"isPublished"`
	OpenGraphData      OpenGraphData     `json:"openGraph"       yaml:"openGraph"`
	URL                string            `json:"url" yaml:"-"`
}

type ArticleList struct {
	Search     string     `json:"-"`
	Tag        string     `json:"-"`
	Items      []*Article `json:"items"`
	TotalPages int        `json:"totalPages"`
	Page       int        `json:"page"`
	PerPage    int        `json:"perPage"`
}

func UnmarshalArticle(data []byte) (Article, error) {
	var article Article
	err := json.Unmarshal(data, &article)

	if err != nil {
		return article, err
	}

	article.PublishedAtDisplay = article.PublishedAt.Format("Jan 2, 2006")
	article.URL = fmt.Sprintf(
		"/article/%d-%02d-%02d/%s",
		article.PublishedAt.Year(),
		article.PublishedAt.Month(),
		article.PublishedAt.Day(),
		article.Slug,
	)

	return article, err
}
