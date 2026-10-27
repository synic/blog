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

func main() {
	log.SetOutput(os.Stdout)
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

	var (
		queries     *db.Queries
		commentRepo *store.CommentRepository
		viewRepo    *store.PageViewRepository
	)

	if conf.DatabaseURL == "" {
		log.Println("⚠️ DATABASE_URL is not set; comments, auth, and page views are disabled")
	} else {
		sqlDB := setupDatabase(conf)
		defer sqlDB.Close()

		queries = db.New(sqlDB)

		commentRepo = store.NewCommentRepository(queries)
		if err := commentRepo.LoadCounts(ctx); err != nil {
			log.Fatal(err)
		}

		go cleanExpiredSessions(queries)

		viewRepo = store.NewPageViewRepository(queries, repo)
	}

	mailer := mail.NewMailer(conf)
	articleController := controller.NewArticleController(repo, commentRepo, viewRepo)
	commentController := controller.NewCommentController(commentRepo, repo, queries, mailer, conf)
	authController := controller.NewAuthController(queries, conf)
	leaderboardController := controller.NewLeaderboardController(viewRepo)

	mux := http.NewServeMux()

	authMW := middleware.AuthMiddleware(queries, conf.AdminEmail)
	csrfMW := middleware.CSRFMiddleware()
	cacheMW := middleware.CacheControlMiddleware()

	server := &http.Server{
		Addr:              conf.HttpAddress,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		Handler: middleware.Wrap(
			mux,
			middleware.LoggerMiddleware(),
			middleware.HtmxMiddleware(),
			middleware.SearchMiddleware(),
			middleware.SecurityHeadersMiddleware(),
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
		cacheMW,
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

func setupDatabase(conf config.Config) *sql.DB {
	sqlDB, err := sql.Open("pgx", conf.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	goose.SetBaseFS(os.DirFS(conf.MigrationsDir))
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}
	if err := goose.Up(sqlDB, "."); err != nil {
		log.Fatal(err)
	}

	return sqlDB
}

func cleanExpiredSessions(queries *db.Queries) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	for range ticker.C {
		if n, err := queries.DeleteExpiredSessions(context.Background()); err != nil {
			log.Printf("Error cleaning expired sessions: %v", err)
		} else if n > 0 {
			log.Printf("Cleaned %d expired sessions", n)
		}
	}
}
