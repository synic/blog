package main

import (
	"context"
	"database/sql"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"

	"github.com/synic/blog/internal"
	"github.com/synic/blog/internal/config"
	"github.com/synic/blog/internal/controller"
	"github.com/synic/blog/internal/db"
	"github.com/synic/blog/internal/mail"
	"github.com/synic/blog/internal/middleware"
	"github.com/synic/blog/internal/model"
	"github.com/synic/blog/internal/store"
	"github.com/synic/blog/internal/view"
)

func main() {
	ctx := context.Background()
	conf := config.Load()
	staticFS := os.DirFS(conf.StaticDir)

	articlesFS, err := fs.Sub(staticFS, "articles")
	if err != nil {
		log.Fatal(err)
	}

	repo, res, err := store.NewArticleRepositoryFromFS(
		articlesFS,
		conf.Debug,
	)

	if err != nil {
		log.Fatal(err)
	}

	res.PrintOutput()

	bundledAssets, err := view.BundleStaticAssets(
		staticFS,
		"css/main.min.css",
		"css/syntax.min.css",
		"js/app.min.js",
	)

	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := sql.Open("sqlite", conf.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	for _, pragma := range []string{
		"PRAGMA journal_mode=WAL",
		"PRAGMA busy_timeout=5000",
		"PRAGMA foreign_keys=ON",
		"PRAGMA synchronous=NORMAL",
	} {
		if _, err := sqlDB.Exec(pragma); err != nil {
			log.Fatal(err)
		}
	}

	goose.SetBaseFS(os.DirFS(conf.MigrationsDir))
	if err := goose.SetDialect("sqlite3"); err != nil {
		log.Fatal(err)
	}
	if err := goose.Up(sqlDB, "."); err != nil {
		log.Fatal(err)
	}

	queries := db.New(sqlDB)

	commentRepo := store.NewCommentRepository(queries)
	if err := commentRepo.LoadCounts(ctx); err != nil {
		log.Fatal(err)
	}

	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			if n, err := queries.DeleteExpiredSessions(context.Background()); err != nil {
				log.Printf("Error cleaning expired sessions: %v", err)
			} else if n > 0 {
				log.Printf("Cleaned %d expired sessions", n)
			}
		}
	}()

	viewRepo := store.NewPageViewRepository(queries, repo)

	mailer := mail.NewMailer(conf)
	articleController := controller.NewArticleController(repo, commentRepo, viewRepo)
	commentController := controller.NewCommentController(commentRepo, repo, queries, mailer, conf)
	authController := controller.NewAuthController(queries, conf)
	leaderboardController := controller.NewLeaderboardController(viewRepo)

	mux := http.NewServeMux()

	authMW := middleware.AuthMiddleware(queries, conf.AdminEmail)
	csrfMW := middleware.CSRFMiddleware()

	server := &http.Server{
		Addr: ":3000",
		Handler: middleware.Wrap(
			mux,
			middleware.LoggerMiddleware(),
			middleware.HtmxMiddleware(),
		),
		BaseContext: func(net.Listener) context.Context {
			data := model.ContextData{
				BuildTime:           conf.BuildTime,
				BundledStaticAssets: bundledAssets,
				Debug:               conf.Debug,
			}
			return context.WithValue(ctx, "data", data)
		},
	}

	internal.RegisterRoutes(
		mux,
		staticFS,
		authMW,
		csrfMW,
		articleController,
		commentController,
		authController,
		leaderboardController,
	)

	log.Printf("🚀 Serving on %s...", server.Addr)
	if err = server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
