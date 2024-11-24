package store

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"maps"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/synic/blog/internal/model"
)

var (
	ErrNotFound = errors.New("article not found")
)

type ArticleRepository struct {
	tags            map[string]int
	slugLookupIndex map[string]int
	articles        []*model.Article
}

func NewArticleRepository(
	articles []*model.Article,
) (*ArticleRepository, error) {
	r := &ArticleRepository{
		slugLookupIndex: make(map[string]int, len(articles)),
		tags:            make(map[string]int, len(articles)*2),
	}

	sort.Slice(articles, func(i, j int) bool {
		return articles[i].PublishedAt.After(articles[j].PublishedAt)
	})

	for i, a := range articles {
		r.slugLookupIndex[a.Slug] = i
		for _, t := range a.Tags {
			count, _ := r.tags[t]
			count += 1
			r.tags[t] = count
		}
	}

	r.articles = articles

	return r, nil
}

func NewArticleRepositoryFromFS(
	filesystem fs.FS,
	includeUnpublished bool,
) (*ArticleRepository, time.Duration, error) {
	begin := time.Now()

	articles, err := parseArticles(filesystem, includeUnpublished)

	if err != nil {
		return nil, time.Since(begin), err
	}

	repo, err := NewArticleRepository(articles)

	return repo, time.Since(begin), err
}

func (r *ArticleRepository) TagInfo(context.Context) map[string]int {
	return r.tags
}

func (r *ArticleRepository) Tags(context.Context) []string {
	return slices.Collect(maps.Keys(r.tags))
}

func (r *ArticleRepository) FindAll(context.Context) ([]*model.Article, error) {
	return r.articles, nil
}

func (r *ArticleRepository) Search(
	_ context.Context,
	terms string,
) ([]*model.Article, error) {
	articles := make([]*model.Article, 0, 10)
	terms = strings.ToLower(terms)

	for _, article := range r.articles {
		if strings.Contains(strings.ToLower(article.Title), terms) ||
			strings.Contains(strings.ToLower(article.Summary), terms) {
			articles = append(articles, article)
		} else {
			for _, t := range article.Tags {
				if strings.Contains(strings.ToLower(t), terms) {
					articles = append(articles, article)
				}
			}
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

	return nil, ErrNotFound
}

func (r *ArticleRepository) FindAllTags(ctx context.Context) ([]string, error) {
	return r.Tags(ctx), nil
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

func (r *ArticleRepository) Count(ctx context.Context) int {
	return len(r.articles)
}

func parseArticles(filesystem fs.FS, includeUnpublished bool) ([]*model.Article, error) {
	var articles []*model.Article

	err := fs.WalkDir(filesystem, ".", func(name string, _ fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(name) != ".json" {
			return nil
		}

		var article model.Article
		data, err := fs.ReadFile(filesystem, name)

		if err != nil {
			return err
		}

		err = json.Unmarshal(data, &article)

		if err != nil {
			return err
		}

		if includeUnpublished || article.IsPublished {
			articles = append(articles, &article)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return articles, nil
}
