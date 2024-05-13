package config

import (
	"log"

	"github.com/synic/adamthings.me/internal/storage"
)

type SiteConfig struct {
	Bind               string
	Repo               storage.ArticleRepository
	ArticlesPerPage    int
	MaxArticlesPerPage int
	IsDebugging        bool
}

func getDefaultConfig() SiteConfig {
	return SiteConfig{
		IsDebugging:        false,
		ArticlesPerPage:    30,
		MaxArticlesPerPage: 50,
		Bind:               ":3000",
	}
}

func New(
	repo storage.ArticleRepository,
	options ...func(*SiteConfig),
) SiteConfig {
	conf := getDefaultConfig()
	conf.Repo = repo

	for _, option := range options {
		option(&conf)
	}

	if conf.IsDebugging {
		log.Println("Starting in DEBUG mode...")
	}

	return conf
}

func WithPagination(perPage int, maxPerPage int) func(*SiteConfig) {
	return func(conf *SiteConfig) {
		conf.ArticlesPerPage = perPage
		conf.MaxArticlesPerPage = maxPerPage
	}
}

func WithIsDebugging(debug bool) func(*SiteConfig) {
	return func(conf *SiteConfig) {
		conf.IsDebugging = debug
	}
}

func WithBind(bind string) func(*SiteConfig) {
	return func(conf *SiteConfig) {
		conf.Bind = bind
	}
}
