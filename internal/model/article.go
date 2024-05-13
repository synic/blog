package model

import (
	"time"
)

type Article struct {
	PublishedAt time.Time
	Body        string
	Summary     string
	Title       string
	Url         string
	Slug        string
	Tags        []string
	IsPublished bool
}
