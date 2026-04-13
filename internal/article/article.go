package article

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/synic/blog/internal/model"
	"github.com/synic/blog/internal/text"
)

var summaryFormatRe = regexp.MustCompile(`\n(.)`)

func ArticleBaseName(slug string, publishedAt time.Time) string {
	return fmt.Sprintf(
		"%d-%02d-%02d_%s",
		publishedAt.Year(),
		publishedAt.Month(),
		publishedAt.Day(),
		slug,
	)
}

func ArticleFileName(articlesDir string, slug string, publishedAt time.Time) string {
	return fmt.Sprintf(
		"%s/%s.md",
		articlesDir,
		ArticleBaseName(slug, publishedAt),
	)
}

func CreateArticle(
	article model.ArticleCreatePayload,
	createImageDir bool,
) (string, error) {
	slug := text.Slugify(article.Title)
	if createImageDir {
		articleBaseName := ArticleBaseName(slug, article.PublishedAt)
		log.Printf("Creating image directory `./assets/images/articles/%s`...", articleBaseName)
		err := os.MkdirAll(fmt.Sprintf("./assets/images/articles/%s", articleBaseName), 0755)
		if err != nil {
			log.Fatal(err)
			return "", fmt.Errorf(
				"Could not create image directory `./assets/images/articles/%s`.",
				articleBaseName,
			)
		}
	}

	fn, content := CreateBlankArticleTemplate(article)

	f, err := os.Create(fn)

	if err != nil {
		return fn, err
	}

	_, err = f.WriteString(content)

	if err != nil {
		return fn, err
	}

	return fn, nil
}

func CreateBlankArticleTemplate(
	article model.ArticleCreatePayload,
) (string, string) {
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

func BuildArticleURL(publishedAt time.Time, slug string) string {
	return fmt.Sprintf(
		"/article/%d-%02d-%02d/%s",
		publishedAt.Year(),
		publishedAt.Month(),
		publishedAt.Day(),
		slug,
	)
}
