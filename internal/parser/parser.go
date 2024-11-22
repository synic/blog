package parser

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/synic/adamthings.me/internal/model"
)

var (
	markdown       = newRenderer()
	mdRegex        = regexp.MustCompile(`(?s)<!-- :metadata:(.*?)-->(.*)`)
	summaryRegex   = regexp.MustCompile(`(?s)(.*)\nsummary:\n\n(.*?)$`)
	checkLineRegex = regexp.MustCompile(`^[a-zA-Z_]*?: (.*?)$`)
	requiredFields = []string{"metadata", "Body", "Title", "Tags"}
)

type parsedData struct {
	extra       map[string]string
	metadata    string
	Body        string
	Title       string
	PublishedAt string
	Tags        string
	Summary     string
}

func capFirst(input string) string {
	var out = ""

	for k, v := range input {
		if k == 0 {
			out += strings.ToUpper(string(v))
		} else {
			out += string(v)
		}
	}

	return out
}

func checkRequired(data parsedData) error {
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)
	for _, field := range requiredFields {
		f, ok := t.FieldByName(field)

		if !ok {
			return fmt.Errorf("required field not found in metadata: %s", field)
		}

		value := v.FieldByIndex(f.Index).String()

		if value == "" {
			return fmt.Errorf("required field not found in metadata: %s", field)
		}
	}

	return nil
}

func parseMetadata(content string) (parsedData, error) {
	data := parsedData{}
	matches := mdRegex.FindStringSubmatch(content)

	if matches == nil {
		return data, errors.New("unable to parse metadata block")
	}

	data.metadata = strings.TrimSpace(matches[1])
	data.Body = strings.TrimSpace(matches[2])
	extra := make(map[string]string)

	md := data.metadata

	matches = summaryRegex.FindStringSubmatch(md)

	if matches != nil {
		md = strings.TrimSpace(matches[1])
		data.Summary = strings.TrimSpace(matches[2])
	}

	lines := strings.Split(md, "\n")
	v := reflect.ValueOf(&data).Elem()

	for _, line := range lines {
		if !checkLineRegex.MatchString(line) {
			continue
		}

		parts := strings.SplitN(line, ": ", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		if field := v.FieldByName(capFirst(key)); field.IsValid() &&
			field.CanSet() && field.String() == "" {
			field.SetString(value)
		} else {
			extra[key] = value
		}
	}

	data.extra = extra

	err := checkRequired(data)

	if err != nil {
		return data, err
	}

	return data, nil
}

func parseTags(tagString string) []string {
	untrimmedTags := strings.Split(tagString, ",")
	tags := make([]string, 0, len(untrimmedTags))

	for _, tag := range untrimmedTags {
		tags = append(tags, strings.TrimSpace(tag))
	}

	return tags
}

func parsePublishedAt(rawPublishedAt string) (time.Time, bool, error) {
	if rawPublishedAt == "" {
		return time.Now(), false, nil
	}

	t, err := time.Parse("2006-01-02T15:04:05-0700", rawPublishedAt)

	if err != nil {
		return time.Time{}, false, err
	}

	return t, true, nil
}

func parseSlug(fn string) (string, error) {
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

func parseArticleFromData(filepath, content string) (model.Article, error) {
	var article model.Article
	data, err := parseMetadata(content)

	if err != nil {
		return article, err
	}

	slug, err := parseSlug(filepath)

	if err != nil {
		return article, fmt.Errorf("unable to parse article slug: %w", err)
	}

	tags := parseTags(data.Tags)

	publishedAt, isPublished, err := parsePublishedAt(data.PublishedAt)

	if err != nil {
		return article, fmt.Errorf("unable to parse publish date: %w", err)
	}

	summaryHtml, err := markdown.MarkdownToHtml(data.Summary)

	if err != nil {
		return article, fmt.Errorf("error converting article summary to html: %w", err)
	}

	bodyHtml, err := markdown.MarkdownToHtml(data.Body)

	if err != nil {
		return article, fmt.Errorf("error converting article body to html: %w", err)
	}

	return model.Article{
		Slug:        slug,
		Title:       data.Title,
		Summary:     summaryHtml,
		Extra:       data.extra,
		PublishedAt: publishedAt,
		Tags:        tags,
		IsPublished: isPublished,
		Body:        bodyHtml,
	}, nil
}

func Parse(filepath string) (model.Article, error) {
	content, err := os.ReadFile(filepath)

	if err != nil {
		return model.Article{}, err
	}

	return parseArticleFromData(filepath, string(content))
}
