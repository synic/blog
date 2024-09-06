package cmd

import (
	"context"
	"log"
	"net/http"

	"github.com/sethvargo/go-envconfig"
	"github.com/spf13/cobra"

	"github.com/synic/adamthings.me/internal/route"
	"github.com/synic/adamthings.me/internal/route/middleware"
	"github.com/synic/adamthings.me/internal/storage"
)

func serve() {
	var conf config

	if err := envconfig.Process(context.Background(), &conf); err != nil {
		log.Fatal(err)
	}

	repo, err := storage.NewArticleRepository("./articles", isDebugBuild)

	if err != nil {
		log.Fatal(err)
	}

	server := http.NewServeMux()
	route.NewDefaultRouter("./assets").Mount(server)
	route.NewArticleRouter(repo).Mount(server)

	log.Printf("Serving on %s...", conf.Bind)

	wrapped := middleware.Wrap(server,
		middleware.LoggingMiddleware(log.Default()),
		middleware.IsHtmxPartialMiddleware(),
	)

	if err = http.ListenAndServe(conf.Bind, wrapped); err != nil {
		log.Fatal(err)
	}
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the webserver",
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
