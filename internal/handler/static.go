package handler

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
)

type staticHandlerConfig struct {
	path    string
	subtree string
	maxAge  int
}

func getDefaultStaticHandlerConfig() staticHandlerConfig {
	return staticHandlerConfig{
		maxAge:  31536000, // 1 year
		path:    "/static/",
		subtree: "assets",
	}
}

func WithCacheControlMaxAge(maxAge int) func(*staticHandlerConfig) {
	return func(conf *staticHandlerConfig) {
		conf.maxAge = maxAge
	}
}

func WithPath(path string) func(*staticHandlerConfig) {
	return func(conf *staticHandlerConfig) {
		conf.path = path
	}
}

func WithSubtree(subtree string) func(*staticHandlerConfig) {
	return func(conf *staticHandlerConfig) {
		conf.subtree = subtree
	}
}

func StaticHandler(
	filesystem fs.FS,
	options ...func(*staticHandlerConfig),
) http.Handler {
	var err error = nil

	conf := getDefaultStaticHandlerConfig()

	for _, option := range options {
		option(&conf)
	}

	if conf.subtree != "" {
		filesystem, err = fs.Sub(filesystem, conf.subtree)

		if err != nil {
			log.Fatal(err)
		}
	}

	staticHandler := http.StripPrefix(conf.path, http.FileServer(http.FS(filesystem)))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", conf.maxAge))
		staticHandler.ServeHTTP(w, r)
	})
}
