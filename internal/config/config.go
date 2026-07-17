package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL        string
	GitHubClientID     string
	GitHubClientSecret string
	SiteUrl            string
	ResendAPIKey       string
	AdminEmail         string
	HttpAddress        string
	StaticDir          string
	MigrationsDir      string
	Debug              bool
	BuildTime          string
}

var (
	// build time (set during build)
	BuildTime string
	// set by the build process (so has to be exported), but use `Debug` instead since
	// it's a boolean
	DebugFlag string
	Debug     bool = false
)

func Load() Config {
	if DebugFlag == "true" {
		godotenv.Load()
		log.Println("🐝 Debugging enabled!")
		Debug = true
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "./data/db.sqlite"
	}

	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		staticDir = "./static"
	}

	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		migrationsDir = "./migrations"
	}

	httpAddress := os.Getenv("HTTP_ADDRESS")

	if httpAddress == "" {
		httpAddress = ":3000"
	}

	return Config{
		DatabaseURL:        databaseURL,
		GitHubClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		GitHubClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		SiteUrl:            os.Getenv("SITE_URL"),
		ResendAPIKey:       os.Getenv("RESEND_API_KEY"),
		AdminEmail:         os.Getenv("ADMIN_EMAIL"),
		HttpAddress:        httpAddress,
		StaticDir:          staticDir,
		MigrationsDir:      migrationsDir,
		Debug:              Debug,
		BuildTime:          BuildTime,
	}
}
