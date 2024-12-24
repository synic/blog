package converter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/synic/blog/internal/model"
)

func TestParseMetadataValid(t *testing.T) {
	content := `---
title: Test Article
publishedAt: 2024-01-01T00:00:00Z
tags: [test, article]
summary: Test summary
---
Article content`

	expected := model.Article{
		Title:       "Test Article",
		PublishedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Tags:        []string{"test", "article"},
		Summary:     "Test summary",
		OpenGraphData: model.OpenGraphData{
			Title: "Test Article",
			Type:  "article",
		},
	}

	article, body, err := parseMetadata(content)
	assert.NoError(t, err)
	assert.Equal(t, expected, article)
	assert.Equal(t, "Article content", body)
}

func TestParseMetadataWithCustomOpenGraph(t *testing.T) {
	content := `---
title: Test Article
publishedAt: 2024-01-01T00:00:00Z
tags: [test]
openGraph:
  title: Custom OG Title
---
Content`

	expected := model.Article{
		Title:       "Test Article",
		PublishedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Tags:        []string{"test"},
		OpenGraphData: model.OpenGraphData{
			Title: "Custom OG Title",
			Type:  "article",
		},
	}

	article, body, err := parseMetadata(content)
	assert.NoError(t, err)
	assert.Equal(t, expected, article)
	assert.Equal(t, "Content", body)
}

func TestParseMetadataMissingFrontmatter(t *testing.T) {
	_, _, err := parseMetadata("No frontmatter here")
	assert.ErrorContains(t, err, "unable to parse frontmatter block")
}

func TestParseMetadataMissingTitle(t *testing.T) {
	content := `---
tags: [test]
---
Content`

	_, _, err := parseMetadata(content)
	assert.ErrorContains(t, err, "title is required")
}

func TestParseMetadataMissingTags(t *testing.T) {
	content := `---
title: Test
---
Content`

	_, _, err := parseMetadata(content)
	assert.ErrorContains(t, err, "tags are required")
}

func TestParseMetadataInvalidYAML(t *testing.T) {
	content := `---
title: [Invalid Syntax
---`

	_, _, err := parseMetadata(content)
	assert.ErrorContains(t, err, "unable to parse frontmatter block")
}

func TestParseSlugWithDate(t *testing.T) {
	slug, err := parseSlug("2024-01-01_test-article.md")
	assert.NoError(t, err)
	assert.Equal(t, "test-article", slug)
}

func TestParseSlugWithoutDate(t *testing.T) {
	slug, err := parseSlug("test-article.md")
	assert.NoError(t, err)
	assert.Equal(t, "test-article", slug)
}

func TestParseSlugNotMarkdown(t *testing.T) {
	_, err := parseSlug("test.txt")
	assert.ErrorContains(t, err, "file was not a markdown file")
}

func TestParseSlugTooManyParts(t *testing.T) {
	_, err := parseSlug("2024-01-01_part1_test.md")
	assert.ErrorContains(t, err, "invalid number of parts")
}

func TestParseArticleFromDataValid(t *testing.T) {
	content := `---
title: Test Article
publishedAt: 2024-01-01T00:00:00Z
tags: [test]
summary: Test summary
---
Article content`

	expected := model.Article{
		Slug:        "test-article",
		Title:       "Test Article",
		PublishedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Tags:        []string{"test"},
		Summary:     "<p>Test summary</p>\n",
		Body:        "<p>Article content</p>\n",
		IsPublished: true,
		OpenGraphData: model.OpenGraphData{
			Title: "Test Article",
			Type:  "article",
		},
	}

	article, err := parseArticleFromData("test-article.md", content)
	assert.NoError(t, err)
	assert.Equal(t, expected, article)
}

func TestParseArticleFromDataUnpublished(t *testing.T) {
	content := `---
title: Draft
tags: [draft]
summary: Draft summary
---
Draft content`

	article, err := parseArticleFromData("draft.md", content)
	assert.NoError(t, err)

	assert.Equal(t, "draft", article.Slug)
	assert.Equal(t, "Draft", article.Title)
	assert.Equal(t, []string{"draft"}, article.Tags)
	assert.Equal(t, "<p>Draft summary</p>\n", article.Summary)
	assert.Equal(t, "<p>Draft content</p>\n", article.Body)
	assert.False(t, article.IsPublished)
	assert.False(t, article.PublishedAt.IsZero())
	assert.Equal(t, "Draft", article.OpenGraphData.Title)
}

func TestParseArticleFromDataInvalidFrontmatter(t *testing.T) {
	_, err := parseArticleFromData("test.md", "Invalid content")
	assert.ErrorContains(t, err, "unable to parse frontmatter block")
}

func TestParseArticleFromDataInvalidFilepath(t *testing.T) {
	content := `---
title: Test
tags: [test]
---
Content`

	_, err := parseArticleFromData("test.txt", content)
	assert.ErrorContains(t, err, "unable to parse article slug")
}

func TestParseNonexistentFile(t *testing.T) {
	_, err := Parse("nonexistent.md")
	assert.ErrorContains(t, err, "no such file or directory")
}
