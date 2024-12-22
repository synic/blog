package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

type ArticleRepository interface {
	FindAll(ctx context.Context) ([]*model.Article, error)
	FindByTag(ctx context.Context, tag string) ([]*model.Article, error)
	Search(ctx context.Context, query string) ([]*model.Article, error)
	FindOneBySlug(ctx context.Context, slug string) (*model.Article, error)
	TagInfo(ctx context.Context) map[string]int
}

type FSArticleRepository struct {
	tags            map[string]int
	slugLookupIndex map[string]int
	articles        []*model.Article
}

func NewFSArticleRepository(
	articles []*model.Article,
) (*FSArticleRepository, error) {
	r := &FSArticleRepository{
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
) (*FSArticleRepository, ParseResult, error) {

	res, err := parseArticles(filesystem, includeUnpublished)

	if err != nil {
		return nil, res, err
	}

	repo, err := NewFSArticleRepository(res.Articles)

	return repo, res, err
}

func (r *FSArticleRepository) TagInfo(context.Context) map[string]int {
	return r.tags
}

func (r *FSArticleRepository) Tags(context.Context) []string {
	return slices.Collect(maps.Keys(r.tags))
}

func (r *FSArticleRepository) FindAll(context.Context) ([]*model.Article, error) {
	return r.articles, nil
}

func (r *FSArticleRepository) Search(
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

func (r *FSArticleRepository) FindOneBySlug(
	_ context.Context,
	slug string,
) (*model.Article, error) {
	if i, ok := r.slugLookupIndex[strings.ToLower(slug)]; ok {
		return r.articles[i], nil
	}

	return nil, ErrNotFound
}

func (r *FSArticleRepository) FindAllTags(ctx context.Context) ([]string, error) {
	return r.Tags(ctx), nil
}

func (r *FSArticleRepository) FindByTag(
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

func (r *FSArticleRepository) Count(ctx context.Context) int {
	return len(r.articles)
}

type ParseResult struct {
	Articles           []*model.Article
	Duration           time.Duration
	Unpublished        int
	Count              int
	IncludeUnpublished bool
}

func (r ParseResult) String() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("🔖 Read %d articles in %s", r.Count, r.Duration))

	if r.IncludeUnpublished {
		b.WriteString(fmt.Sprintf(", including %d unpublished", r.Unpublished))
	}

	return b.String()
}

func parseArticles(filesystem fs.FS, includeUnpublished bool) (ParseResult, error) {
	res := ParseResult{IncludeUnpublished: includeUnpublished}
	begin := time.Now()

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
			res.Articles = append(res.Articles, &article)
			res.Count += 1
			if !article.IsPublished {
				res.Unpublished += 1
			}
		}

		return nil
	})

	if err != nil {
		return res, err
	}

	res.Duration = time.Since(begin)

	return res, nil
}
