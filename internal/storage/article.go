package storage

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/synic/adamthings.me/internal/model"
)

type ArticleRepository struct {
	tags            []string
	slugLookupIndex map[string]int
	articles        []model.Article
}

func getArticleUrlPath(a model.Article) string {
	return fmt.Sprintf(
		"/articles/%d-%02d-%02d/%s",
		a.PublishedAt.Year(),
		a.PublishedAt.Month(),
		a.PublishedAt.Day(),
		a.Slug,
	)
}

func NewArticleRepository(
	articleDirectory string,
	isDebugging bool,
) (ArticleRepository, error) {
	repository, err := readAllArticles(articleDirectory, isDebugging)
	if err != nil {
		return ArticleRepository{}, err
	}

	log.Printf(
		"read %d articles and %d tags...\n",
		len(repository.articles),
		len(repository.tags),
	)

	return repository, nil
}

func (r ArticleRepository) FindAll(_ context.Context) ([]model.Article, error) {
	return r.articles, nil
}

func (r ArticleRepository) Search(_ context.Context, terms string) ([]model.Article, error) {
	articles := make([]model.Article, 0, 10)
	terms = strings.ToLower(terms)

	for _, article := range r.articles {
		if strings.Contains(strings.ToLower(article.Title), terms) ||
			strings.Contains(strings.ToLower(article.Summary), terms) {
			articles = append(articles, article)
		}
	}

	return articles, nil
}

func (r ArticleRepository) FindOneBySlug(
	_ context.Context,
	slug string,
) (model.Article, error) {
	if i, ok := r.slugLookupIndex[strings.ToLower(slug)]; ok {
		return r.articles[i], nil
	}

	return model.Article{}, fmt.Errorf("article for slug `%s` not found", slug)
}

func (r ArticleRepository) FindAllTags(context.Context) ([]string, error) {
	return r.tags, nil
}

func (r ArticleRepository) FindByTag(
	_ context.Context,
	tag string,
) ([]model.Article, error) {
	articles := make([]model.Article, 0, 10)

	for _, article := range r.articles {
		for _, t := range article.Tags {
			if strings.ToLower(t) == strings.ToLower(tag) {
				articles = append(articles, article)
			}
		}
	}

	return articles, nil
}
