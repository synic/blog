package converter

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/synic/blog/internal/model"
	"gopkg.in/yaml.v3"
)

var (
	markdown         = newRenderer()
	frontmatterRegex = regexp.MustCompile(`(?s)\A---\n(.*?)\n---\n(.*)`)
)

func parseMetadata(content string) (model.Article, string, error) {
	matches := frontmatterRegex.FindStringSubmatch(content)
	if matches == nil {
		return model.Article{}, "", errors.New("unable to parse frontmatter block")
	}

	var article model.Article
	err := yaml.Unmarshal([]byte(matches[1]), &article)
	if err != nil {
		return model.Article{}, "", fmt.Errorf("error parsing frontmatter: %w", err)
	}

	// default opengraph
	if article.OpenGraphData.Title == "" {
		article.OpenGraphData.Title = article.Title
	}

	article.OpenGraphData.Type = "article"

	// Validate required fields
	if article.Title == "" {
		return model.Article{}, "", errors.New("title is required")
	}

	if article.Slug == "" {
		return model.Article{}, "", errors.New("slug is required")
	}

	if len(article.Tags) == 0 {
		return model.Article{}, "", errors.New("tags are required")
	}

	return article, strings.TrimSpace(matches[2]), nil
}

func parseArticleFromData(content string) (model.Article, error) {
	article, body, err := parseMetadata(content)

	if err != nil {
		return article, err
	}

	article.IsPublished = true

	if article.PublishedAt.IsZero() {
		article.IsPublished = false
		article.PublishedAt = time.Now()
	}

	summaryHtml, err := markdown.MarkdownToHtml(article.Summary)

	if err != nil {
		return article, fmt.Errorf("error converting article summary to html: %w", err)
	}

	bodyHtml, err := markdown.MarkdownToHtml(body)

	if err != nil {
		return article, fmt.Errorf("error converting article body to html: %w", err)
	}

	article.Summary = summaryHtml
	article.Body = bodyHtml
	return article, nil
}

func Parse(filepath string) (model.Article, error) {
	content, err := os.ReadFile(filepath)

	if err != nil {
		return model.Article{}, err
	}

	return parseArticleFromData(string(content))
}
