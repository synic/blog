package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/synic/adamthings.me/internal/handler"
	"github.com/synic/adamthings.me/internal/middleware"
	"github.com/synic/adamthings.me/internal/model"
	"github.com/synic/adamthings.me/internal/store"
)

var (
	BuildTime = fmt.Sprintf("%d", time.Now().Unix())
	// will maybe be set to `true` in `debug.go`
	isDebugBuild = false
	//go:embed articles/json/*
	embeddedArticles embed.FS
	//go:embed assets/*
	embeddedAssets embed.FS
)

func main() {
	var articles []*model.Article

	entries, err := embeddedArticles.ReadDir("articles/json")

	if err != nil {
		log.Fatal(err)
	}

	begin := time.Now()
	for _, entry := range entries {
		var article model.Article

		name := path.Join("articles", "json", entry.Name())

		if filepath.Ext(name) != ".json" {
			continue
		}

		data, err := embeddedArticles.ReadFile(name)

		if err != nil {
			log.Fatalf("error reading article %s: %v", name, err)
		}

		err = json.Unmarshal(data, &article)

		if err != nil {
			log.Fatalf("error unmarshaling article %s: %v", name, err)
		}

		articles = append(articles, &article)
	}

	log.Printf("🔖 Read %d articles in %s", len(articles), time.Since(begin))

	repo, err := store.NewArticleRepository(articles)

	if err != nil {
		log.Fatal(err)
	}

	bind := os.Getenv("BIND")
	if bind == "" {
		bind = "0.0.0.0:3000"
	}

	server := http.NewServeMux()
	server.Handle("GET /static/", handler.StaticHandler(embeddedAssets))
	handler.NewArticleRouter(repo).Mount(server)

	log.Printf("🚀 Serving on %s...", bind)

	wrapped := middleware.Wrap(
		server,
		middleware.AddContextData(map[string]any{"BuildTime": BuildTime}),
		middleware.CacheStaticFiles(embeddedAssets, "css/main.css"),
		middleware.LoggingMiddleware(log.Default()),
		middleware.IsHtmxPartialMiddleware(),
		middleware.GzipMiddleware(),
	)

	if err = http.ListenAndServe(bind, wrapped); err != nil {
		log.Fatal(err)
	}
}
