package converter

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/synic/blog/internal/model"
)

func TestParseMetadataValid(t *testing.T) {
	content := `---
title: Test Article
slug: test-article
publishedAt: 2024-01-01T00:00:00Z
tags: [test, article]
---
Article content`

	expected := model.Article{
		Title:       "Test Article",
		PublishedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Tags:        []string{"test", "article"},
		Slug:        "test-article",
		OpenGraphData: model.OpenGraphData{
			Title: "Test Article",
			Type:  "article",
		},
	}
	expected.Prepare()

	article, body, err := parseMetadata(content)
	assert.NoError(t, err)
	assert.Equal(t, expected, article)
	assert.Equal(t, "Article content", body)
}

func TestParseMetadataWithCustomOpenGraph(t *testing.T) {
	content := `---
title: Test Article
slug: test-article
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
		Slug:        "test-article",
		OpenGraphData: model.OpenGraphData{
			Title: "Custom OG Title",
			Type:  "article",
		},
	}
	expected.Prepare()

	article, body, err := parseMetadata(content)
	assert.NoError(t, err)
	assert.Equal(t, expected, article)
	assert.Equal(t, "Content", body)
}

func TestParseMetadataIgnoresSummaryKey(t *testing.T) {
	content := `---
title: Test Article
slug: test-article
publishedAt: 2024-01-01T00:00:00Z
tags: [test]
summary: should be ignored
---
Body`

	article, _, err := parseMetadata(content)
	assert.NoError(t, err)
	assert.Equal(t, "", article.Summary)
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
slug: test
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

func TestParseArticleFromDataValid(t *testing.T) {
	content := `---
title: Test Article
slug: test-article
publishedAt: 2024-01-01T00:00:00Z
tags: [test]
---
<!-- summary -->
Test summary
<!-- end-summary -->

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
	expected.Prepare()

	article, err := parseArticleFromData(content)
	assert.NoError(t, err)
	assert.Equal(t, expected, article)
}

func TestParseArticleFromDataUnpublished(t *testing.T) {
	content := `---
title: Draft
slug: draft
tags: [draft]
---
<!-- summary -->
Draft summary
<!-- end-summary -->

Draft content`

	article, err := parseArticleFromData(content)
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

func TestParseArticleFromDataNoSummaryBlock(t *testing.T) {
	content := `---
title: Test
slug: test
tags: [test]
publishedAt: 2024-01-01T00:00:00Z
---
Just the body.`

	article, err := parseArticleFromData(content)
	assert.NoError(t, err)
	assert.Equal(t, "", article.Summary)
	assert.Equal(t, "<p>Just the body.</p>\n", article.Body)
}

func TestParseArticleFromDataEmptySummaryBlock(t *testing.T) {
	content := `---
title: Test
slug: test
tags: [test]
publishedAt: 2024-01-01T00:00:00Z
---
<!-- summary -->
<!-- end-summary -->

Body here.`

	article, err := parseArticleFromData(content)
	assert.NoError(t, err)
	assert.Equal(t, "", article.Summary)
	assert.Equal(t, "<p>Body here.</p>\n", article.Body)
}

func TestParseArticleFromDataRenderInBody(t *testing.T) {
	content := `---
title: Test
slug: test
tags: [test]
publishedAt: 2024-01-01T00:00:00Z
---
<!-- summary render-in-body=true -->
My summary.
<!-- end-summary -->

## Heading

Trailing content.`

	article, err := parseArticleFromData(content)
	assert.NoError(t, err)
	assert.Equal(t, "<p>My summary.</p>\n", article.Summary)
	assert.Contains(t, article.Body, "<p>My summary.</p>")
	assert.Contains(t, article.Body, "<h2")
	assert.True(t, strings.HasPrefix(article.Body, "<p>My summary.</p>"))
}

func TestParseArticleFromDataNoRenderInBody(t *testing.T) {
	content := `---
title: Test
slug: test
tags: [test]
publishedAt: 2024-01-01T00:00:00Z
---
<!-- summary -->
My summary.
<!-- end-summary -->

Body content.`

	article, err := parseArticleFromData(content)
	assert.NoError(t, err)
	assert.Equal(t, "<p>My summary.</p>\n", article.Summary)
	assert.NotContains(t, article.Body, "<p>My summary.</p>")
	assert.Contains(t, article.Body, "<p>Body content.</p>")
}

func TestParseArticleFromDataMultipleSummaryBlocks(t *testing.T) {
	content := `---
title: Test
slug: test
tags: [test]
publishedAt: 2024-01-01T00:00:00Z
---
<!-- summary -->
First.
<!-- end-summary -->

<!-- summary -->
Second.
<!-- end-summary -->

Body.`

	_, err := parseArticleFromData(content)
	assert.ErrorContains(t, err, "multiple summary blocks")
}

func TestParseArticleFromDataMismatchedMarkers(t *testing.T) {
	content := `---
title: Test
slug: test
tags: [test]
publishedAt: 2024-01-01T00:00:00Z
---
<!-- summary -->
No closing marker.

Body.`

	_, err := parseArticleFromData(content)
	assert.ErrorContains(t, err, "opening summary marker without closing marker")
}

func TestParseArticleFromDataClosingWithoutOpening(t *testing.T) {
	content := `---
title: Test
slug: test
tags: [test]
publishedAt: 2024-01-01T00:00:00Z
---
<!-- end-summary -->

Body.`

	_, err := parseArticleFromData(content)
	assert.ErrorContains(t, err, "closing summary marker without opening marker")
}

func TestParseArticleFromDataMarkersInCodeBlockIgnored(t *testing.T) {
	content := "---\n" +
		"title: Test\n" +
		"slug: test\n" +
		"tags: [test]\n" +
		"publishedAt: 2024-01-01T00:00:00Z\n" +
		"---\n" +
		"Body intro.\n\n" +
		"```\n" +
		"<!-- summary -->\n" +
		"fake summary inside code\n" +
		"<!-- end-summary -->\n" +
		"```\n\n" +
		"Trailing."

	article, err := parseArticleFromData(content)
	assert.NoError(t, err)
	assert.Equal(t, "", article.Summary)
	assert.Contains(t, article.Body, "&lt;!-- summary --&gt;")
	assert.Contains(t, article.Body, "&lt;!-- end-summary --&gt;")
}

func TestParseSummaryOptions(t *testing.T) {
	opts := parseSummaryOptions("render-in-body=true,foo=bar")
	assert.Equal(t, "true", opts["render-in-body"])
	assert.Equal(t, "bar", opts["foo"])

	opts = parseSummaryOptions("  render-in-body=true , foo=bar  ")
	assert.Equal(t, "true", opts["render-in-body"])
	assert.Equal(t, "bar", opts["foo"])

	opts = parseSummaryOptions("")
	assert.Empty(t, opts)
}

func TestParseArticleFromDataRenderInBodyWithMultipleOptions(t *testing.T) {
	content := `---
title: Test
slug: test
tags: [test]
publishedAt: 2024-01-01T00:00:00Z
---
<!-- summary render-in-body=true,foo=bar -->
My summary.
<!-- end-summary -->

Body content.`

	article, err := parseArticleFromData(content)
	assert.NoError(t, err)
	assert.Equal(t, "<p>My summary.</p>\n", article.Summary)
	assert.True(t, strings.HasPrefix(article.Body, "<p>My summary.</p>"))
	assert.Contains(t, article.Body, "<p>Body content.</p>")
}

func TestParseArticleFromDataInvalidFrontmatter(t *testing.T) {
	_, err := parseArticleFromData("Invalid content")
	assert.ErrorContains(t, err, "unable to parse frontmatter block")
}

func TestParseNonexistentFile(t *testing.T) {
	_, err := Parse("nonexistent.md")
	assert.ErrorContains(t, err, "no such file or directory")
}
