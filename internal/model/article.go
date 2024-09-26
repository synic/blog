package model

import "time"

type Article struct {
	PublishedAt time.Time `json:"published_at"`
	Body        string    `json:"body"`
	Summary     string    `json:"summary"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Tags        []string  `json:"tags"`
	IsPublished bool      `json:"is_published"`
}
