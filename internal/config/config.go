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

	return Config{
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		GitHubClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		GitHubClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		ServerAddress:      os.Getenv("SERVER_ADDRESS"),
		ResendAPIKey:       os.Getenv("RESEND_API_KEY"),
		AdminEmail:         os.Getenv("ADMIN_EMAIL"),
		Debug:              Debug,
		BuildTime:          BuildTime,
	}
}
