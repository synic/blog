package view

import (
	"context"
	"embed"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/a-h/templ"
)

func inlinestatic(files ...string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		if fs, ok := ctx.Value("staticFS").(embed.FS); ok {
			for _, file := range files {
				data, err := fs.ReadFile(fmt.Sprintf("assets/%s", file))
				wrap := "script"

				if strings.HasSuffix(file, ".css") {
					wrap = "style"
				}

				if err != nil {
					log.Printf("unable to read static file %s: %v", file, err)
					return err
				}

				_, err = io.WriteString(
					w,
					fmt.Sprintf("<%s>/* %s */\n%s</%s>", wrap, file, data, wrap),
				)

				if err != nil {
					return err
				}
			}

			return nil
		}

		return fmt.Errorf("unable to locate static files \"%s\"", strings.Join(files, ","))
	})
}
