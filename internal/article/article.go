package article

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/synic/blog/internal/model"
	"github.com/synic/blog/internal/text"
)

var summaryFormatRe = regexp.MustCompile(`\n(.)`)

func ArticleFileName(articlesDir string, slug string, publishedAt time.Time) string {
	return fmt.Sprintf(
		"%s/%d-%02d-%02d_%s.md",
		articlesDir,
		publishedAt.Year(),
		publishedAt.Month(),
		publishedAt.Day(),
		slug,
	)
}

func CreateBlankArticleTemplate(article model.ArticleCreatePayload) (string, string) {
	slug := text.Slugify(article.Title)
	fn := ArticleFileName("./articles", slug, article.PublishedAt)

	needsQuotes := strings.ContainsAny(article.Title, `:"'[]{}#|>&*?!`)
	titleField := article.Title
	if needsQuotes {
		titleField = fmt.Sprintf("%q", article.Title)
	}

	tagList := strings.Split(article.Tags, ",")
	for i, tag := range tagList {
		tagList[i] = strings.TrimSpace(tag)
	}
	tagsField := fmt.Sprintf("[%s]", strings.Join(tagList, ", "))

	summary := strings.TrimSpace(article.Summary)
	summary = strings.TrimSpace(
		string(summaryFormatRe.ReplaceAll([]byte(summary), []byte("\n  $1"))),
	)

	return fn, fmt.Sprintf(`---
title: %s
slug: %s
tags: %s
publishedAt: %s
summary: |
%s
---
%s
`, titleField, slug, tagsField, article.PublishedAt.Format(time.RFC3339), summary, article.Body)
}
