package main

import (
	"context"
	"embed"
	"log"
	"net"
	"net/http"

	"github.com/synic/blog/internal"
	"github.com/synic/blog/internal/controller"
	"github.com/synic/blog/internal/middleware"
	"github.com/synic/blog/internal/model"
	"github.com/synic/blog/internal/store"
	"github.com/synic/blog/internal/view"
)

//go:embed assets/*
var embeddedAssets embed.FS

func main() {
	ctx := context.Background()
	assets := internal.MustSub(embeddedAssets, "assets")

	repo, res, err := store.NewArticleRepositoryFromFS(
		internal.MustSub(assets, "articles"),
		internal.Debug,
	)

	if err != nil {
		log.Fatal(err)
	}

	res.PrintOutput()

	bundledAssets, err := view.BundleStaticAssets(
		assets,
		"css/main.css",
		"css/syntax.min.css",
		"js/app.js",
	)

	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	server := &http.Server{
		Addr: ":3000",
		Handler: middleware.Wrap(
			mux,
			middleware.LoggerMiddleware(),
			middleware.HtmxMiddleware(),
		),
		BaseContext: func(net.Listener) context.Context {
			data := model.ContextData{
				BuildTime:           internal.BuildTime,
				BundledStaticAssets: bundledAssets,
				Debug:               internal.Debug,
			}
			return context.WithValue(ctx, "data", data)
		},
	}

	internal.RegisterRoutes(
		mux,
		assets,
		controller.NewArticleController(repo),
		controller.NewArticleRSSController(repo),
	)

	log.Printf("🚀 Serving on %s...", server.Addr)
	if err = server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
