package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/synic/blog/internal/model"
	"github.com/synic/blog/internal/store"
)

type mockArticleRepository struct {
	articles []*model.Article
	tagInfo  map[string]int
	err      error
}

func (m *mockArticleRepository) FindAll(ctx context.Context) ([]*model.Article, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.articles, nil
}

func (m *mockArticleRepository) FindByTag(
	ctx context.Context,
	tag string,
) ([]*model.Article, error) {
	if m.err != nil {
		return nil, m.err
	}
	var filtered []*model.Article
	for _, article := range m.articles {
		for _, t := range article.Tags {
			if t == tag {
				filtered = append(filtered, article)
				break
			}
		}
	}
	return filtered, nil
}

func (m *mockArticleRepository) Search(
	ctx context.Context,
	query string,
) ([]*model.Article, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.articles, nil
}

func (m *mockArticleRepository) FindOneBySlug(
	ctx context.Context,
	slug string,
) (*model.Article, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, article := range m.articles {
		if article.Slug == slug {
			return article, nil
		}
	}
	return nil, store.ErrNotFound
}

func (m *mockArticleRepository) TagInfo(ctx context.Context) map[string]int {
	return m.tagInfo
}

func TestWithPagination(t *testing.T) {
	tests := []struct {
		name           string
		perPage        int
		maxPerPage     int
		expectedConfig articleControllerConfig
	}{
		{
			name:       "custom pagination",
			perPage:    10,
			maxPerPage: 30,
			expectedConfig: articleControllerConfig{
				articlesPerPage:    10,
				maxArticlesPerPage: 30,
			},
		},
		{
			name:       "zero values",
			perPage:    0,
			maxPerPage: 0,
			expectedConfig: articleControllerConfig{
				articlesPerPage:    0,
				maxArticlesPerPage: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &articleControllerConfig{}
			WithPagination(tt.perPage, tt.maxPerPage)(conf)
			assert.Equal(t, tt.expectedConfig, *conf)
		})
	}
}

func TestNewArticleController(t *testing.T) {
	repo := &mockArticleRepository{}

	tests := []struct {
		name           string
		options        []func(*articleControllerConfig)
		expectedConfig articleControllerConfig
	}{
		{
			name: "default configuration",
			expectedConfig: articleControllerConfig{
				articlesPerPage:    20,
				maxArticlesPerPage: 50,
			},
		},
		{
			name:    "custom configuration",
			options: []func(*articleControllerConfig){WithPagination(10, 30)},
			expectedConfig: articleControllerConfig{
				articlesPerPage:    10,
				maxArticlesPerPage: 30,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := NewArticleController(repo, tt.options...)
			assert.Equal(t, tt.expectedConfig, controller.articleControllerConfig)
		})
	}
}

func TestArticleController_Index(t *testing.T) {
	mockArticles := []*model.Article{
		{Title: "Test1", Slug: "test1", Tags: []string{"tag1"}},
		{Title: "Test2", Slug: "test2", Tags: []string{"tag2"}},
	}

	tests := []struct {
		name           string
		repo           *mockArticleRepository
		queryString    string
		formSearch     string
		expectedStatus int
	}{
		{
			name:           "successful retrieval",
			repo:           &mockArticleRepository{articles: mockArticles},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "repository error",
			repo:           &mockArticleRepository{err: errors.New("db error")},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "search by query string",
			repo:           &mockArticleRepository{articles: mockArticles},
			queryString:    "?search=test",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "search by form value",
			repo:           &mockArticleRepository{articles: mockArticles},
			formSearch:     "test",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "filter by tag",
			repo:           &mockArticleRepository{articles: mockArticles},
			queryString:    "?tag=tag1",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := NewArticleController(tt.repo)
			req := httptest.NewRequest(http.MethodGet, "/"+tt.queryString, nil)
			if tt.formSearch != "" {
				req.Form = make(map[string][]string)
				req.Form.Set("search", tt.formSearch)
			}
			w := httptest.NewRecorder()

			controller.Index(w, req)
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestArticleController_Article(t *testing.T) {
	mockArticles := []*model.Article{
		{Title: "Test1", Slug: "test1"},
	}

	tests := []struct {
		name           string
		repo           *mockArticleRepository
		slug           string
		expectedStatus int
	}{
		{
			name:           "article found",
			repo:           &mockArticleRepository{articles: mockArticles},
			slug:           "test1",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "article not found",
			repo:           &mockArticleRepository{articles: mockArticles},
			slug:           "nonexistent",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "repository error",
			repo:           &mockArticleRepository{err: errors.New("db error")},
			slug:           "test1",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := NewArticleController(tt.repo)
			req := httptest.NewRequest(http.MethodGet, "/article/"+tt.slug, nil)
			req.SetPathValue("slug", tt.slug)
			w := httptest.NewRecorder()

			controller.Article(w, req)
			assert.Equal(
				t,
				tt.expectedStatus,
				w.Code,
				fmt.Sprintf(
					"slug: %s, expected status: %d, got status: %d",
					tt.slug,
					tt.expectedStatus,
					w.Code,
				),
			)
		})
	}
}

func TestArticleController_Archive(t *testing.T) {
	mockArticles := []*model.Article{
		{Title: "Test1", Slug: "test1"},
		{Title: "Test2", Slug: "test2"},
	}
	mockTagInfo := map[string]int{"tag1": 1, "tag2": 2}

	tests := []struct {
		name           string
		repo           *mockArticleRepository
		expectedStatus int
	}{
		{
			name: "successful archive",
			repo: &mockArticleRepository{
				articles: mockArticles,
				tagInfo:  mockTagInfo,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "empty archive",
			repo: &mockArticleRepository{
				articles: []*model.Article{},
				tagInfo:  map[string]int{},
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := NewArticleController(tt.repo)
			req := httptest.NewRequest(http.MethodGet, "/archive", nil)
			w := httptest.NewRecorder()

			controller.Archive(w, req)
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestArticleController_renderAndPageArticles(t *testing.T) {
	articles := make([]*model.Article, 25)
	for i := 0; i < 25; i++ {
		articles[i] = &model.Article{Title: "Test" + string(rune(i))}
	}

	tests := []struct {
		name           string
		queryString    string
		articles       []*model.Article
		search         string
		tag            string
		expectedStatus int
		expectedLen    int
	}{
		{
			name:           "default pagination",
			articles:       articles,
			expectedStatus: http.StatusOK,
			expectedLen:    20, // default per page
		},
		{
			name:           "custom page",
			queryString:    "?page=2",
			articles:       articles,
			expectedStatus: http.StatusOK,
			expectedLen:    5, // remaining items
		},
		{
			name:           "custom per_page",
			queryString:    "?per_page=10",
			articles:       articles,
			expectedStatus: http.StatusOK,
			expectedLen:    10,
		},
		{
			name:           "invalid page number",
			queryString:    "?page=invalid",
			articles:       articles,
			expectedStatus: http.StatusOK,
			expectedLen:    20,
		},
		{
			name:           "exceed max per_page",
			queryString:    "?per_page=100",
			articles:       articles,
			expectedStatus: http.StatusOK,
			expectedLen:    25, // all items, capped at max
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := NewArticleController(&mockArticleRepository{})
			req := httptest.NewRequest(http.MethodGet, "/"+tt.queryString, nil)
			w := httptest.NewRecorder()

			controller.renderAndPageArticles(w, req, tt.articles, tt.search, tt.tag)
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
