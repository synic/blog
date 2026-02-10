package article

import (
	"fmt"
	"strings"
	"time"

	"github.com/synic/blog/internal/text"
)

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

func CreateBlankArticleTemplate(title string, tags string, publishedAt time.Time) (string, string) {
	slug := text.Slugify(title)
	fn := ArticleFileName("./articles", slug, publishedAt)

	needsQuotes := strings.ContainsAny(title, `:"'[]{}#|>&*?!`)
	titleField := title
	if needsQuotes {
		titleField = fmt.Sprintf("%q", title)
	}

	tagList := strings.Split(tags, ",")
	for i, tag := range tagList {
		tagList[i] = strings.TrimSpace(tag)
	}
	tagsField := fmt.Sprintf("[%s]", strings.Join(tagList, ", "))

	return fn, fmt.Sprintf(`---
title: %s
slug: %s
tags: %s
publishedAt: %s
summary: |

---
`, titleField, slug, tagsField, publishedAt.Format(time.RFC3339))
}
