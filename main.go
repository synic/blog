package main

import (
	"context"
	"database/sql"
	"embed"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

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

//go:embed assets/*
var embeddedAssets embed.FS

//go:embed migrations/*.sql
var embeddedMigrations embed.FS

func main() {
	ctx := context.Background()
	cfg := config.Load()
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

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	sqlDB, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	goose.SetBaseFS(embeddedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}
	if err := goose.Up(sqlDB, "migrations"); err != nil {
		log.Fatal(err)
	}

	queries := db.New(pool)

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

	mailer := mail.NewMailer(cfg)
	articleController := controller.NewArticleController(repo, commentRepo, viewRepo)
	commentController := controller.NewCommentController(commentRepo, repo, queries, mailer, cfg)
	authController := controller.NewAuthController(queries, cfg)
	leaderboardController := controller.NewLeaderboardController(viewRepo)

	mux := http.NewServeMux()

	authMW := middleware.AuthMiddleware(queries, cfg.AdminEmail)
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
