package main

import (
	"log"
	"net/http"

	"github.com/synic/adamthings.me/internal/config"
	"github.com/synic/adamthings.me/internal/pkg/env"
	"github.com/synic/adamthings.me/internal/route"
	"github.com/synic/adamthings.me/internal/route/middleware"
	"github.com/synic/adamthings.me/internal/storage"
)

func main() {
	e := env.New()

	repo, err := storage.NewArticleRepository("./articles", e.IsDebugBuild)

	if err != nil {
		log.Fatal(err)
	}

	conf := config.New(repo,
		config.WithIsDebugging(e.IsDebugBuild),
		config.WithBind(e.GetString("BIND", ":3000")),
	)

	server := http.NewServeMux()
	route.NewDefaultRouter(conf, "./assets").Mount(server)
	route.NewArticleRouter(conf).Mount(server)

	log.Printf("Serving on %s...", conf.Bind)

	wrapped := middleware.Wrap(server,
		middleware.LoggingMiddleware(log.Default()),
		middleware.IsHtmxPartialMiddleware(),
	)

	if err = http.ListenAndServe(conf.Bind, wrapped); err != nil {
		log.Fatal(err)
	}
}
