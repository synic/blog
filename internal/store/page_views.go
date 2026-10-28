package store

import (
	"context"
	"log"
	"strings"
	"time"

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

const maxPageViewBatch = 500

type PageViewRepository struct {
	queries *db.Queries
	repo    ArticleRepository
	queue   chan db.CreatePageViewsParams
}

func NewPageViewRepository(queries *db.Queries, repo ArticleRepository) *PageViewRepository {
	r := &PageViewRepository{
		queries: queries,
		repo:    repo,
		queue:   make(chan db.CreatePageViewsParams, 1024),
	}
	go r.worker()
	return r
}

func (r *PageViewRepository) worker() {
	batch := make([]db.CreatePageViewsParams, 0, maxPageViewBatch)
	for view := range r.queue {
		batch = append(batch[:0], view)
	drain:
		for len(batch) < maxPageViewBatch {
			select {
			case view := <-r.queue:
				batch = append(batch, view)
			default:
				break drain
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if _, err := r.queries.CreatePageViews(ctx, batch); err != nil {
			log.Printf("Error logging %d page views: %v", len(batch), err)
		}
		cancel()
	}
}

func (r *PageViewRepository) LogView(slug, ip, userAgent string) {
	if r == nil || isBot(userAgent) {
		return
	}
	select {
	case r.queue <- db.CreatePageViewsParams{
		ArticleSlug: slug,
		IpAddress:   ip,
		UserAgent:   userAgent,
	}:
	default:
		log.Printf("Page view queue full, dropping view for %s", slug)
	}
}

func (r *PageViewRepository) ViewCounts(ctx context.Context) ([]model.PageViewEntry, error) {
	if r == nil {
		return []model.PageViewEntry{}, nil
	}

	rows, err := r.queries.CountPageViewsBySlug(ctx)
	if err != nil {
		return nil, err
	}

	entries := make([]model.PageViewEntry, 0, len(rows))
	for _, row := range rows {
		title := row.ArticleSlug
		articleURL := ""
		if article, err := r.repo.FindOneBySlug(ctx, row.ArticleSlug); err == nil {
			title = article.Title
			articleURL = article.URL
		} else {
			continue
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
