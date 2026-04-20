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
	ServerAddress      string
	ResendAPIKey       string
	AdminEmail         string
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

	return Config{
		DatabaseURL:        databaseURL,
		GitHubClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		GitHubClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		ServerAddress:      os.Getenv("SERVER_ADDRESS"),
		ResendAPIKey:       os.Getenv("RESEND_API_KEY"),
		AdminEmail:         os.Getenv("ADMIN_EMAIL"),
		StaticDir:          staticDir,
		MigrationsDir:      migrationsDir,
		Debug:              Debug,
		BuildTime:          BuildTime,
	}
}
