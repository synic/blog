package middleware

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"
)

func readStaticFiles(filesystem fs.ReadFileFS, files ...string) (string, error) {
	content := ""
	for _, file := range files {
		data, err := filesystem.ReadFile(fmt.Sprintf("assets/%s", file))
		wrap := "script"

		if strings.HasSuffix(file, ".css") {
			wrap = "style"
		}

		if err != nil {
			log.Printf("unable to read static file %s: %v", file, err)
			return "", err
		}

		content += fmt.Sprintf("<%s>/* %s */\n%s</%s>", wrap, file, data, wrap)
	}

	return content, nil
}

func CacheStaticFiles(filesystem fs.ReadFileFS, files ...string) func(http.Handler) http.Handler {
	content, err := readStaticFiles(filesystem, files...)

	if err != nil {
		log.Fatal(err)
	}

	data := []byte(content)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "cached-static-files", &data)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
