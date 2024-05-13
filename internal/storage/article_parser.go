package storage

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/synic/adamthings.me/internal/model"
	"github.com/synic/adamthings.me/internal/pkg/markdown"
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

func parseArticlePublishedAtFromMetadata(block string) (time.Time, error) {
	rawPublishedAt, err := parseArticleValueFromMetadata(block, "published")

	if err != nil {
		return time.Time{}, err
	}
	t, err := time.Parse("2006-01-02T15:04:05-0700", rawPublishedAt)

	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

func parseArticle(name string, body string) (model.Article, error) {
	slug, err := parseArticleSlugFromFileName(name)

	if err != nil {
		return model.Article{}, fmt.Errorf("unable to parse article slug: %w", err)
	}

	md, err := parseArticleMetadataBlock(body)

	if err != nil {
		return model.Article{}, err
	}

	summary, _ := parseArticleSummaryFromMetadata(md)
	title, err := parseArticleValueFromMetadata(md, "title")

	if err != nil {
		return model.Article{}, errors.New(
			"unable to parse metadata from article content, content not found",
		)
	}

	tags := parseArticleTagsFromMetadata(md)

	publishedAt, err := parseArticlePublishedAtFromMetadata(md)
	isPublished := err == nil

	summaryHtml, err := markdownRenderer.MarkdownToHtml(summary)

	if err != nil {
		return model.Article{}, fmt.Errorf("error converting article summary to html: %v", err)
	}

	bodyHtml, err := markdownRenderer.MarkdownToHtml(body)

	if err != nil {
		return model.Article{}, fmt.Errorf("error converting article body to html: %v", err)
	}

	return model.Article{
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
	if !strings.HasSuffix(fn, ".md") {
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

func readAllArticles(
	articleDirectory string,
	isDebuging bool,
) (ArticleRepository, error) {
	entries, err := os.ReadDir(articleDirectory)

	if err != nil {
		return ArticleRepository{}, fmt.Errorf("unable to open article directory: %w", err)
	}

	log.Println("loading articles, please wait...")

	articles := make([]model.Article, 0, len(entries))
	slugLookupIndex := make(map[string]int, len(entries))
	tags := make([]string, 0, 15)

	for _, entry := range entries {
		name := entry.Name()
		path := path.Join("articles", name)

		if strings.HasSuffix(path, ".md") {
			data, err := os.ReadFile(path)

			if err != nil {
				log.Printf("error reading article `%s`: %v", path, err)
				continue
			}

			article, err := parseArticle(name, string(data))

			if err != nil {
				log.Printf("error parsing article `%s`: %s", name, err)
				continue
			}

			if !article.IsPublished {
				if !isDebuging {
					log.Printf("article file `%s` does not contain a publish date, skipping.", name)
					continue
				}
				article.PublishedAt = time.Now()
			}

			for _, tag := range article.Tags {
				if !slices.Contains(tags, tag) {
					tags = append(tags, tag)
				}
			}

			article.Url = getArticleUrlPath(article)
			articles = append(articles, article)
		}
	}

	sort.Slice(articles, func(i, j int) bool {
		return articles[i].PublishedAt.After(articles[j].PublishedAt)
	})

	for i, article := range articles {
		slugLookupIndex[article.Slug] = i
	}

	if len(tags) > 0 {
		log.Printf("found tags: %+v\n", strings.Join(tags, ", "))
	} else {
		log.Println("didn't find any tags; strange")
	}

	return ArticleRepository{
		articles:        articles,
		slugLookupIndex: slugLookupIndex,
		tags:            tags,
	}, nil
}
