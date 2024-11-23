package main

import (
	"context"
	"embed"
	"log"
	"net"
	"net/http"

	"github.com/synic/adamthings.me/internal"
	"github.com/synic/adamthings.me/internal/controller"
	"github.com/synic/adamthings.me/internal/middleware"
	"github.com/synic/adamthings.me/internal/model"
	"github.com/synic/adamthings.me/internal/store"
	"github.com/synic/adamthings.me/internal/view"
)

var (
	//go:embed assets/*
	embeddedAssets embed.FS
)

func main() {
	ctx := context.Background()
	assets := internal.MustSub(embeddedAssets, "assets")

	repo, duration, err := store.NewArticleRepositoryFromFS(
		internal.MustSub(assets, "articles"),
		internal.Debug,
	)

	if err != nil {
		log.Fatal(err)
	}

	if internal.Debug {
		log.Println("🐝 Debugging enabled, unpublished articles will be shown")
	}

	log.Printf("🔖 Read %d articles in %s", repo.Count(ctx), duration)

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

	internal.RegisterRoutes(mux, assets, controller.NewArticleController(repo))

	log.Printf("🚀 Serving on %s...", server.Addr)
	if err = server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
