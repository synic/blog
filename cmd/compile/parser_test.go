package main

import (
	"slices"
	"testing"
)

func TestParseArticleMetadataBlock(t *testing.T) {
	data, err := parseArticleMetadataBlock(`
<!-- :metadata:

title: this is a test
published-at: 2024-01-10T15:04:04-0800
tags: Test,Foo
summary:

This is a good summary!

-->

This is the article.
title: just trying to confuse the parser with this title
`)

	if err != nil {
		t.Errorf("unexpected parsing error: %s\n", err)
	}

	expected := `title: this is a test
published-at: 2024-01-10T15:04:04-0800
tags: Test,Foo
summary:

This is a good summary!`

	if data != expected {
		t.Errorf("parsed metadata content was different than expected: %s", data)
	}

}

func TestParseArticleSummaryFromContent(t *testing.T) {
	summary, err := parseArticleSummaryFromMetadata(`
title: this is a test
published-at: 2024-01-10T15:04:04-0800
tags: Test,Foo
summary:

This is a good summary!`)

	if err != nil {
		t.Errorf("unexpected error when parsing summary: %s\n", err)
	}

	expected := `This is a good summary!`

	if summary != expected {
		t.Errorf("parsed summary was different than expected: %s", summary)
	}
}

func TestParseArticleValueFromMetadata(t *testing.T) {
	md := `
value1: one
value2: two
value3: three
`

	v, err := parseArticleValueFromMetadata(md, "value1")

	if err != nil {
		t.Error(err)
	}

	if v != "one" {
		t.Errorf("expected one, got %s", v)
	}

	v, err = parseArticleValueFromMetadata(md, "value2")

	if err != nil {
		t.Error(err)
	}

	if v != "two" {
		t.Errorf("expected two, got %s", v)
	}

	v, err = parseArticleValueFromMetadata(md, "value3")

	if err != nil {
		t.Error(err)
	}

	if v != "three" {
		t.Errorf("expected three, got %s", v)
	}
}

func TestParseArticleTagsFromMetadata(t *testing.T) {
	md := `
title: this is a test
published-at: 2024-01-10T15:04:04-0800
tags: Test,Foo
summary:

This is a good summary!
	`

	tags := parseArticleTagsFromMetadata(md)

	if slices.Compare(tags, []string{"Test", "Foo"}) != 0 {
		t.Errorf("invalid tags received: %+v", tags)
	}
}
