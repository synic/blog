package component

import (
	"context"

	"github.com/synic/adamthings.me/internal/route/middleware"
)

func isPartial(ctx context.Context) bool {
	if isPartial, ok := ctx.Value(middleware.IsHtmxPartialContextKey).(bool); ok {
		return isPartial
	}
	return false
}
