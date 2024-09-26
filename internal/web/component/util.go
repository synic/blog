package component

import (
	"context"
	"fmt"

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
