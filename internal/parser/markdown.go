package parser

import (
	"bytes"
	"fmt"

	formatters "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	renderer "github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/anchor"
)

type markdownRenderer struct {
	renderer goldmark.Markdown
}

func newRenderer() markdownRenderer {
	markdown := goldmark.New(
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithFormatOptions(formatters.WithClasses(true)),
			),
			extension.NewLinkify(
				extension.WithLinkifyAllowedProtocols([]string{"http:", "https:"}),
			),
			&anchor.Extender{
				Attributer: anchor.Attributes{
					"class": "header-anchor",
				},
				Texter: anchor.Text("  ¶")},
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
