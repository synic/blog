package store

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"sort"
	"strings"

	"github.com/synic/adamthings.me/internal/model"
)

type ArticleRepository struct {
	tags            []string
	slugLookupIndex map[string]int
	articles        []*model.Article
}

func NewArticleRepository(
	articles []*model.Article,
) (*ArticleRepository, error) {
	tags := make(map[string]struct{})
	r := &ArticleRepository{
		slugLookupIndex: make(map[string]int, len(articles)),
	}

	sort.Slice(articles, func(i, j int) bool {
		return articles[i].PublishedAt.After(articles[j].PublishedAt)
	})

	for i, a := range articles {
		r.slugLookupIndex[a.Slug] = i
		for _, t := range a.Tags {
			tags[t] = struct{}{}
		}
	}

	r.articles = articles
	r.tags = slices.Collect(maps.Keys(tags))

	return r, nil
}

func (r *ArticleRepository) FindAll(_ context.Context) ([]*model.Article, error) {
	return r.articles, nil
}

func (r *ArticleRepository) Search(_ context.Context, terms string) ([]*model.Article, error) {
	articles := make([]*model.Article, 0, 10)
	terms = strings.ToLower(terms)

	for _, article := range r.articles {
		if strings.Contains(strings.ToLower(article.Title), terms) ||
			strings.Contains(strings.ToLower(article.Summary), terms) {
			articles = append(articles, article)
		}
	}

	return articles, nil
}

func (r *ArticleRepository) FindOneBySlug(
	_ context.Context,
	slug string,
) (*model.Article, error) {
	if i, ok := r.slugLookupIndex[strings.ToLower(slug)]; ok {
		return r.articles[i], nil
	}

	return nil, fmt.Errorf("article for slug `%s` not found", slug)
}

func (r *ArticleRepository) FindAllTags(context.Context) ([]string, error) {
	return r.tags, nil
}

func (r *ArticleRepository) FindByTag(
	_ context.Context,
	tag string,
) ([]*model.Article, error) {
	articles := make([]*model.Article, 0, 10)

	for _, article := range r.articles {
		for _, t := range article.Tags {
			if strings.ToLower(t) == strings.ToLower(tag) {
				articles = append(articles, article)
			}
		}
	}

	return articles, nil
}
