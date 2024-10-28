package main

import (
	"bytes"
	"fmt"

	formatters "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-highlighting/v2"
	renderer "github.com/yuin/goldmark/renderer/html"
)

type markdownRenderer struct {
	renderer goldmark.Markdown
}

func newMarkdownRenderer() markdownRenderer {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithFormatOptions(formatters.WithClasses(true)),
			),
		),
		goldmark.WithRendererOptions(renderer.WithUnsafe()),
	)

	return markdownRenderer{renderer: markdown}
}

func (r *markdownRenderer) MarkdownToHtml(md string) (string, error) {
	var buf bytes.Buffer

	if err := r.renderer.Convert([]byte(md), &buf); err != nil {
		return "", fmt.Errorf("error converting markdown to html: %v", err)
	}

	return buf.String(), nil
}
