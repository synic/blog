package main

import (
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/synic/adamthings.me/internal/model"
	"github.com/synic/adamthings.me/internal/store"
	"github.com/synic/adamthings.me/internal/web/middleware"
	"github.com/synic/adamthings.me/internal/web/route"
)

var (
	// will maybe be set to `true` in `debug.go`
	isDebugBuild = false
	//go:embed articles/*
	articleFiles embed.FS
	//go:embed assets/*
	assetFiles embed.FS
)

func main() {
	var articles []*model.Article

	entries, err := articleFiles.ReadDir("articles")

	if err != nil {
		log.Fatal(err)
	}

	begin := time.Now()
	for _, entry := range entries {
		var article model.Article

		name := fmt.Sprintf("articles/%s", entry.Name())

		if filepath.Ext(name) != ".json" {
			continue
		}

		data, err := articleFiles.ReadFile(name)

		if err != nil {
			log.Fatalf("error reading article %s: %v", name, err)
		}

		err = json.Unmarshal(data, &article)

		if err != nil {
			log.Fatalf("error unmarshaling article %s: %v", name, err)
		}

		articles = append(articles, &article)
	}

	log.Printf("Read %d articles in %s", len(articles), time.Since(begin))

	repo, err := store.NewArticleRepository(articles)

	if err != nil {
		log.Fatal(err)
	}

	bind := os.Getenv("BIND")
	if bind == "" {
		bind = "0.0.0.0:3000"
	}

	sub, err := fs.Sub(assetFiles, "assets")

	if err != nil {
		log.Fatal(err)
	}

	server := http.NewServeMux()
	static := http.StripPrefix("/static/", http.FileServer(http.FS(sub)))
	server.Handle("GET /static/", static)
	route.NewArticleRouter(repo).Mount(server)

	log.Printf("Serving on %s...", bind)

	wrapped := middleware.Wrap(server,
		middleware.LoggingMiddleware(log.Default()),
		middleware.IsHtmxPartialMiddleware(),
		middleware.GzipMiddleware(),
	)

	if err = http.ListenAndServe(bind, wrapped); err != nil {
		log.Fatal(err)
	}
}
