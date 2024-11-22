package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/synic/adamthings.me/internal"
	"github.com/synic/adamthings.me/internal/handler"
	"github.com/synic/adamthings.me/internal/middleware"
	"github.com/synic/adamthings.me/internal/store"
)

var (
	//go:embed assets/*
	embeddedAssets embed.FS
)

func mustSub(fsys fs.FS, path string) fs.FS {
	st, err := fs.Sub(fsys, path)

	if err != nil {
		panic(err)
	}

	return st
}

func main() {
	var (
		begin  = time.Now()
		assets = mustSub(embeddedAssets, "assets")
	)

	articles, err := internal.ParseArticles(mustSub(assets, "articles"))

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("🔖 Read %d articles in %s", len(articles), time.Since(begin))

	repo, err := store.NewArticleRepository(articles)

	if err != nil {
		log.Fatal(err)
	}

	staticFiles, err := internal.BundleStaticAssets(assets, "css/main.css", "css/syntax.min.css")

	if err != nil {
		log.Fatal(err)
	}

	bind := os.Getenv("BIND")
	if bind == "" {
		bind = "0.0.0.0:3000"
	}

	mux := http.NewServeMux()
	mux.Handle("GET /static/", handler.StaticHandler(assets))
	handler.NewArticleRouter(repo).Mount(mux)

	log.Printf("🚀 Serving on %s...", bind)

	server := &http.Server{
		Addr: bind,
		Handler: middleware.Wrap(
			mux,
			middleware.LoggingMiddleware(log.Default()),
			middleware.HtmxMiddleware(),
		),
		BaseContext: func(net.Listener) context.Context {
			return context.WithValue(context.Background(), "inline-static-files", staticFiles)
		},
	}

	if err = server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
