package store

import (
	"context"
	"log"
	"strings"

	"github.com/synic/blog/internal/db"
	"github.com/synic/blog/internal/model"
)

var botPatterns = []string{
	"bot",
	"crawler",
	"spider",
	"slurp",
	"mediapartners",
	"fetcher",
	"curl",
	"wget",
	"python-requests",
	"go-http-client",
	"scrapy",
	"headlesschrome",
	"phantomjs",
	"semrush",
	"ahrefs",
	"mj12bot",
	"dotbot",
	"rogerbot",
	"screaming frog",
	"yandexbot",
	"baiduspider",
	"duckduckbot",
	"facebookexternalhit",
	"twitterbot",
	"linkedinbot",
	"whatsapp",
	"telegrambot",
	"applebot",
	"ia_archiver",
	"archive.org_bot",
	"uptimerobot",
	"pingdom",
	"gptbot",
	"chatgpt",
	"claudebot",
	"anthropic",
	"bytespider",
	"paloaltonetworks",
}

func isBot(userAgent string) bool {
	ua := strings.ToLower(userAgent)
	for _, pattern := range botPatterns {
		if strings.Contains(ua, pattern) {
			return true
		}
	}
	return false
}

type PageViewRepository struct {
	queries  *db.Queries
	articles ArticleRepository
}

func NewPageViewRepository(queries *db.Queries, articles ArticleRepository) *PageViewRepository {
	return &PageViewRepository{
		queries:  queries,
		articles: articles,
	}
}

func (r *PageViewRepository) LogView(slug, ip, userAgent string) {
	if isBot(userAgent) {
		return
	}
	go func() {
		err := r.queries.CreatePageView(context.Background(), db.CreatePageViewParams{
			ArticleSlug: slug,
			IpAddress:   ip,
			UserAgent:   userAgent,
		})
		if err != nil {
			log.Printf("Error logging page view for %s: %v", slug, err)
		}
	}()
}

func (r *PageViewRepository) ViewCounts(ctx context.Context) ([]model.PageViewEntry, error) {
	rows, err := r.queries.CountPageViewsBySlug(ctx)
	if err != nil {
		return nil, err
	}

	entries := make([]model.PageViewEntry, 0, len(rows))
	for _, row := range rows {
		title := row.ArticleSlug
		articleURL := ""
		if article, err := r.articles.FindOneBySlug(ctx, row.ArticleSlug); err == nil {
			title = article.Title
			articleURL = article.URL()
		}
		entries = append(entries, model.PageViewEntry{
			Slug:      row.ArticleSlug,
			Title:     title,
			URL:       articleURL,
			ViewCount: int(row.ViewCount),
		})
	}

	return entries, nil
}
