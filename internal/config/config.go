package config

import "os"

type Config struct {
	DatabaseURL        string
	GitHubClientID     string
	GitHubClientSecret string
	ServerAddress      string
	ResendAPIKey       string
	AdminEmail         string
}

func Load() Config {
	return Config{
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		GitHubClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		GitHubClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		ServerAddress:      os.Getenv("SERVER_ADDRESS"),
		ResendAPIKey:       os.Getenv("RESEND_API_KEY"),
		AdminEmail:         os.Getenv("ADMIN_EMAIL"),
	}
}
