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

const articleSummaryPlaceholder = "<!-- article-summary -->"

var (
	markdown         = newRenderer()
	frontmatterRegex = regexp.MustCompile(`(?s)\A---\n(.*?)\n---\n(.*)`)
	summaryOpenRe    = regexp.MustCompile(`^\s*<!--\s*summary\s*-->\s*$`)
	summaryCloseRe   = regexp.MustCompile(`^\s*<!--\s*/summary\s*-->\s*$`)
	fenceRe          = regexp.MustCompile("^\\s{0,3}(`{3,}|~{3,})")
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

	article.Prepare()

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

	summaryMd, bodyMd, err := extractSummary(body)

	if err != nil {
		return article, fmt.Errorf("error extracting article summary: %w", err)
	}

	summaryHtml, err := markdown.MarkdownToHtml(summaryMd)

	if err != nil {
		return article, fmt.Errorf("error converting article summary to html: %w", err)
	}

	bodyHtml, err := markdown.MarkdownToHtml(bodyMd)

	if err != nil {
		return article, fmt.Errorf("error converting article body to html: %w", err)
	}

	article.Summary = summaryHtml
	article.Body = substitutePlaceholder(bodyHtml, strings.TrimRight(summaryHtml, "\n"))
	article.Prepare()
	return article, nil
}

func Parse(filepath string) (model.Article, error) {
	content, err := os.ReadFile(filepath)

	if err != nil {
		return model.Article{}, err
	}

	return parseArticleFromData(string(content))
}

func extractSummary(body string) (string, string, error) {
	lines := strings.Split(body, "\n")

	openIdx := -1
	closeIdx := -1
	inFence := false

	for i, line := range lines {
		if fenceRe.MatchString(line) {
			inFence = !inFence
			continue
		}

		if inFence {
			continue
		}

		if summaryOpenRe.MatchString(line) {
			if openIdx != -1 {
				return "", "", errors.New("multiple summary blocks found")
			}
			openIdx = i
			continue
		}

		if summaryCloseRe.MatchString(line) {
			if openIdx == -1 {
				return "", "", errors.New("closing summary marker without opening marker")
			}
			if closeIdx != -1 {
				return "", "", errors.New("multiple summary blocks found")
			}
			closeIdx = i
		}
	}

	if openIdx == -1 && closeIdx == -1 {
		return "", body, nil
	}

	if openIdx != -1 && closeIdx == -1 {
		return "", "", errors.New("opening summary marker without closing marker")
	}

	if closeIdx <= openIdx {
		return "", "", errors.New("mismatched summary markers")
	}

	summary := strings.Join(lines[openIdx+1:closeIdx], "\n")

	remaining := make([]string, 0, len(lines)-(closeIdx-openIdx+1))
	remaining = append(remaining, lines[:openIdx]...)
	remaining = append(remaining, lines[closeIdx+1:]...)

	return strings.TrimSpace(summary), strings.TrimSpace(strings.Join(remaining, "\n")), nil
}

func substitutePlaceholder(bodyHtml, summaryHtml string) string {
	return strings.ReplaceAll(bodyHtml, articleSummaryPlaceholder, summaryHtml)
}
