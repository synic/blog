package view

import (
	"context"
	"embed"
	"fmt"
	"io"
	"log"

	"github.com/a-h/templ"

	"github.com/synic/adamthings.me/internal/model"
	"github.com/synic/adamthings.me/internal/web/middleware"
)

func articleURL(a *model.Article) templ.SafeURL {
	return templ.URL(fmt.Sprintf(
		"/articles/%d-%02d-%02d/%s",
		a.PublishedAt.Year(),
		a.PublishedAt.Month(),
		a.PublishedAt.Day(),
		a.Slug,
	))
}

func isPartial(ctx context.Context) bool {
	if isPartial, ok := ctx.Value(middleware.IsHtmxPartialContextKey).(bool); ok {
		return isPartial
	}
	return false
}

func inlinestatic(file, wrap string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		if fs, ok := ctx.Value("staticFS").(embed.FS); ok {
			data, err := fs.ReadFile(fmt.Sprintf("assets/%s", file))

			if err != nil {
				log.Printf("unable to read static file %s: %v", file, err)
				return err
			}

			_, err = io.WriteString(w, fmt.Sprintf("<%s>%s</%s>", wrap, data, wrap))
			return err
		}

		return fmt.Errorf("unable to locate static file %s", file)
	})
}
