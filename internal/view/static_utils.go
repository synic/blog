package view

import (
	"context"
	"fmt"
	"io"

	"github.com/a-h/templ"
)

func cachedstaticfiles() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		if data, ok := ctx.Value("cached-static-files").(*[]byte); ok {
			_, err := w.Write(*data)
			return err
		}

		return fmt.Errorf("unable to locate cached static files")
	})
}
