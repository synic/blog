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
<!-- /summary -->

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
<!-- /summary -->

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
<!-- /summary -->

Body here.`

	article, err := parseArticleFromData(content)
	assert.NoError(t, err)
	assert.Equal(t, "", article.Summary)
	assert.Equal(t, "<p>Body here.</p>\n", article.Body)
}

func TestParseArticleFromDataPlaceholder(t *testing.T) {
	content := `---
title: Test
slug: test
tags: [test]
publishedAt: 2024-01-01T00:00:00Z
---
<!-- summary -->
My summary.
<!-- /summary -->

Intro paragraph.

<!-- article-summary -->

## Heading

Trailing content.`

	article, err := parseArticleFromData(content)
	assert.NoError(t, err)
	assert.Equal(t, "<p>My summary.</p>\n", article.Summary)
	assert.Contains(t, article.Body, "<p>Intro paragraph.</p>")
	assert.Contains(t, article.Body, "<p>My summary.</p>")
	assert.Contains(t, article.Body, "<h2")
	assert.NotContains(t, article.Body, "<!-- article-summary -->")
}

func TestParseArticleFromDataMultiplePlaceholders(t *testing.T) {
	content := `---
title: Test
slug: test
tags: [test]
publishedAt: 2024-01-01T00:00:00Z
---
<!-- summary -->
Shared summary.
<!-- /summary -->

<!-- article-summary -->

Middle.

<!-- article-summary -->`

	article, err := parseArticleFromData(content)
	assert.NoError(t, err)
	assert.Equal(t, 2, strings.Count(article.Body, "<p>Shared summary.</p>"))
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
<!-- /summary -->

<!-- summary -->
Second.
<!-- /summary -->

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
<!-- /summary -->

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
		"<!-- /summary -->\n" +
		"```\n\n" +
		"Trailing."

	article, err := parseArticleFromData(content)
	assert.NoError(t, err)
	assert.Equal(t, "", article.Summary)
	assert.Contains(t, article.Body, "&lt;!-- summary --&gt;")
	assert.Contains(t, article.Body, "&lt;!-- /summary --&gt;")
}

func TestParseArticleFromDataPlaceholderInCodeBlockNotSubstituted(t *testing.T) {
	content := "---\n" +
		"title: Test\n" +
		"slug: test\n" +
		"tags: [test]\n" +
		"publishedAt: 2024-01-01T00:00:00Z\n" +
		"---\n" +
		"<!-- summary -->\n" +
		"Real summary.\n" +
		"<!-- /summary -->\n\n" +
		"Body.\n\n" +
		"```\n" +
		"<!-- article-summary -->\n" +
		"```\n"

	article, err := parseArticleFromData(content)
	assert.NoError(t, err)
	assert.Equal(t, "<p>Real summary.</p>\n", article.Summary)
	assert.Contains(t, article.Body, "&lt;!-- article-summary --&gt;")
	assert.NotContains(t, article.Body, "Real summary.</p>\n</code>")
	assert.Equal(t, 0, strings.Count(article.Body, "<p>Real summary.</p>"))
}

func TestParseArticleFromDataInvalidFrontmatter(t *testing.T) {
	_, err := parseArticleFromData("Invalid content")
	assert.ErrorContains(t, err, "unable to parse frontmatter block")
}

func TestParseNonexistentFile(t *testing.T) {
	_, err := Parse("nonexistent.md")
	assert.ErrorContains(t, err, "no such file or directory")
}
