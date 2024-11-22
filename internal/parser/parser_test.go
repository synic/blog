package parser

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCapFirst(t *testing.T) {
	assert.Equal(t, "PublishedAt", capFirst("publishedAt"))
}

func TestParseMetadata(t *testing.T) {
	md := `
<!-- :metadata:

title: this is a test
publishedAt: 2024-01-10T15:04:04-0800
tags: Test,Foo
randomField: woot
summary:

This is a good summary!

-->

This is the article.
title: just trying to confuse the parser with this title
`
	data, err := parseMetadata(md)

	assert.Nil(t, err)
	assert.Equal(t, "this is a test", data.Title)
	assert.Equal(t, "2024-01-10T15:04:04-0800", data.PublishedAt)
	assert.Equal(t, "Test,Foo", data.Tags)
	assert.Equal(t, "This is a good summary!", data.Summary)
	assert.Equal(
		t,
		"This is the article.\ntitle: just trying to confuse the parser with this title",
		data.Body,
	)
	assert.Equal(t, `title: this is a test
publishedAt: 2024-01-10T15:04:04-0800
tags: Test,Foo
randomField: woot
summary:

This is a good summary!`, data.metadata)
	assert.Len(t, data.extra, 1)
	f, _ := data.extra["randomField"]
	assert.Equal(t, "woot", f)
}

func TestParseMetadataMissingField(t *testing.T) {
	md := `
<!-- :metadata:

publishedAt: 2024-01-10T15:04:04-0800
tags: Test, Foo
randomField: woot
summary: This is a good summary!
-->

This is the article.
`
	_, err := parseMetadata(md)

	assert.ErrorContains(t, err, "required field not found in metadata: Title")
}

func TestParseMetadataShortSummary(t *testing.T) {
	md := `
<!-- :metadata:

title: woot
publishedAt: 2024-01-10T15:04:04-0800
tags: Test, Foo
random_field: woot
summary: This is a good summary!
-->

This is the article.
`
	data, err := parseMetadata(md)

	assert.Nil(t, err)
	assert.Equal(t, "This is a good summary!", data.Summary)
}

func TestParsePublishedAt(t *testing.T) {
	c, _ := time.Parse("2006-01-02T15:04:05-0700", "2006-01-02T15:04:05-0800")
	publishedAt, isPublished, err := parsePublishedAt("2024-01-10T15:04:04-0800")

	assert.Nil(t, err)
	assert.True(t, isPublished)
	assert.Equal(
		t,
		time.Time(time.Date(2024, time.January, 10, 15, 4, 4, 0, c.Location())),
		publishedAt,
	)
}

func TestParsePublishedAtNotPublished(t *testing.T) {
	now := time.Now()

	publishedAt, isPublished, err := parsePublishedAt("")

	assert.Nil(t, err)
	assert.True(t, publishedAt.After(now))
	assert.False(t, isPublished)
}

func TestParsePublishedAtInvalid(t *testing.T) {
	_, _, err := parsePublishedAt("woot")

	assert.Error(t, err)
}

func TestParseTags(t *testing.T) {
	tags := parseTags("Foo,Bar")
	assert.Len(t, tags, 2)
	assert.Contains(t, tags, "Foo")
	assert.Contains(t, tags, "Bar")
}

func TestParseTagsExtraSpacing(t *testing.T) {
	tags := parseTags("Foo, Bar  ")
	assert.Len(t, tags, 2)
	assert.Contains(t, tags, "Foo")
	assert.Contains(t, tags, "Bar")
}

func TestParseArticleFromData(t *testing.T) {
	c, _ := time.Parse("2006-01-02T15:04:05-0700", "2006-01-02T15:04:05-0800")
	filepath := "2006-01-01_this-is-the-best-article.md"
	md := `
<!-- :metadata:

title: this is a test
publishedAt: 2024-01-10T15:04:04-0800
tags: Test,Foo
randomField: woot
summary:

This is a good summary!
# Hello!

-->

This is the article.
title: just trying to confuse the parser with this title
`

	article, err := parseArticleFromData(filepath, md)

	assert.Nil(t, err)
	assert.True(t, article.IsPublished)
	assert.Equal(t,
		time.Time(time.Date(2024, time.January, 10, 15, 4, 4, 0, c.Location())),
		article.PublishedAt,
	)
	assert.Equal(t, "this is a test", article.Title)
	assert.Equal(
		t,
		"<p>This is the article.\ntitle: just trying to confuse the parser with this title</p>\n",
		article.Body,
	)
	assert.Equal(
		t,
		"<p>This is a good summary!</p>\n<h1 id=\"hello\">Hello! <a class=\"header-anchor\" href=\"#hello\">  ¶</a></h1>\n",
		article.Summary,
	)
	assert.Len(t, article.Tags, 2)

	assert.Contains(t, article.Tags, "Foo")
	assert.Contains(t, article.Tags, "Test")
	assert.Equal(t, "this-is-the-best-article", article.Slug)
	assert.Len(t, article.Extra, 1)
	assert.Equal(t, map[string]string{"randomField": "woot"}, article.Extra)
}
