package parser

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/synic/adamthings.me/internal/markdown"
	"github.com/synic/adamthings.me/internal/model"
)

var markdownRenderer = markdown.New()

func parseArticleMetadataBlock(content string) (string, error) {
	r := regexp.MustCompile(`(?s)<!-- :metadata:(.*?)-->`)
	matches := r.FindStringSubmatch(content)

	if matches == nil {
		return "", errors.New("unable to parse metadata block")
	}

	return strings.TrimSpace(matches[1]), nil
}

func parseArticleValueFromMetadata(block string, name string) (string, error) {
	r := regexp.MustCompile(fmt.Sprintf(`(?m)^%s: (.*)$`, name))
	matches := r.FindStringSubmatch(block)

	if matches == nil {
		return "", fmt.Errorf("unable to parse `%s` from article content", name)
	}

	return strings.TrimSpace(matches[1]), nil
}

func parseArticleTagsFromMetadata(block string) []string {
	rawTags, err := parseArticleValueFromMetadata(block, "tags")

	if err != nil {
		return []string{}
	}

	untrimmedTags := strings.Split(rawTags, ",")
	tags := make([]string, 0, len(untrimmedTags))

	for _, tag := range untrimmedTags {
		tags = append(tags, strings.TrimSpace(tag))
	}

	return tags
}

func parseArticlePublishedAtFromMetadata(block string) (time.Time, bool, error) {
	rawPublishedAt, err := parseArticleValueFromMetadata(block, "published")

	if err != nil {
		return time.Now(), false, nil
	}
	t, err := time.Parse("2006-01-02T15:04:05-0700", rawPublishedAt)

	if err != nil {
		return time.Time{}, false, err
	}

	return t, true, nil
}

func Parse(name string) (*model.Article, error) {
	content, err := os.ReadFile(name)

	if err != nil {
		return nil, err
	}

	data := string(content)

	slug, err := parseArticleSlugFromFileName(name)

	if err != nil {
		return nil, fmt.Errorf("unable to parse article slug: %w", err)
	}

	md, err := parseArticleMetadataBlock(data)

	if err != nil {
		return nil, err
	}

	summary, _ := parseArticleSummaryFromMetadata(md)
	title, err := parseArticleValueFromMetadata(md, "title")

	if err != nil {
		return nil, errors.New(
			"unable to parse metadata from article content, content not found",
		)
	}

	tags := parseArticleTagsFromMetadata(md)

	publishedAt, isPublished, err := parseArticlePublishedAtFromMetadata(md)

	if err != nil {
		return nil, fmt.Errorf("unable to parse publish date: %w", err)
	}

	summaryHtml, err := markdownRenderer.MarkdownToHtml(summary)

	if err != nil {
		return nil, fmt.Errorf("error converting article summary to html: %w", err)
	}

	bodyHtml, err := markdownRenderer.MarkdownToHtml(data)

	if err != nil {
		return nil, fmt.Errorf("error converting article body to html: %w", err)
	}

	return &model.Article{
		Slug:        slug,
		Title:       title,
		Summary:     summaryHtml,
		PublishedAt: publishedAt,
		Tags:        tags,
		IsPublished: isPublished,
		Body:        bodyHtml,
	}, nil
}

func parseArticleSummaryFromMetadata(metadata string) (string, error) {
	r := regexp.MustCompile(`(?s)\nsummary:\n\n(.*?)$`)

	matches := r.FindStringSubmatch(metadata)

	if matches == nil {
		return "", errors.New("summary not found")
	}

	return strings.TrimSpace(matches[1]), nil
}

func parseArticleSlugFromFileName(fn string) (string, error) {
	if filepath.Ext(fn) != ".md" {
		return "", errors.New("file was not a markdown file")
	}

	parts := strings.Split(fn, "_")

	if len(parts) < 1 || len(parts) > 2 {
		return "", fmt.Errorf("invalid number of parts in file name: %d", len(parts))
	}

	i := len(parts) - 1

	rawSlug := parts[i]
	return strings.ToLower(rawSlug[0 : len(rawSlug)-3]), nil
}
