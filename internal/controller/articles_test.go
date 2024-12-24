package controller

import (
	"context"
	"errors"
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

func TestWithPaginationCustom(t *testing.T) {
	conf := &articleControllerConfig{}
	WithPagination(10, 30)(conf)
	assert.Equal(t, articleControllerConfig{
		articlesPerPage:    10,
		maxArticlesPerPage: 30,
	}, *conf)
}

func TestWithPaginationZeroValues(t *testing.T) {
	conf := &articleControllerConfig{}
	WithPagination(0, 0)(conf)
	assert.Equal(t, articleControllerConfig{
		articlesPerPage:    0,
		maxArticlesPerPage: 0,
	}, *conf)
}

func TestNewArticleControllerDefault(t *testing.T) {
	repo := &mockArticleRepository{}
	controller := NewArticleController(repo)
	assert.Equal(t, articleControllerConfig{
		articlesPerPage:    20,
		maxArticlesPerPage: 50,
	}, controller.articleControllerConfig)
}

func TestNewArticleControllerCustomConfig(t *testing.T) {
	repo := &mockArticleRepository{}
	controller := NewArticleController(repo, WithPagination(10, 30))
	assert.Equal(t, articleControllerConfig{
		articlesPerPage:    10,
		maxArticlesPerPage: 30,
	}, controller.articleControllerConfig)
}

func TestArticleController_IndexSuccessful(t *testing.T) {
	mockArticles := []*model.Article{
		{Title: "Test1", Slug: "test1", Tags: []string{"tag1"}},
		{Title: "Test2", Slug: "test2", Tags: []string{"tag2"}},
	}
	repo := &mockArticleRepository{articles: mockArticles}
	controller := NewArticleController(repo)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	controller.Index(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestArticleController_IndexWithError(t *testing.T) {
	repo := &mockArticleRepository{err: errors.New("db error")}
	controller := NewArticleController(repo)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	controller.Index(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestArticleController_IndexWithQuerySearch(t *testing.T) {
	mockArticles := []*model.Article{
		{Title: "Test1", Slug: "test1", Tags: []string{"tag1"}},
		{Title: "Test2", Slug: "test2", Tags: []string{"tag2"}},
	}
	repo := &mockArticleRepository{articles: mockArticles}
	controller := NewArticleController(repo)

	req := httptest.NewRequest(http.MethodGet, "/?search=test", nil)
	w := httptest.NewRecorder()

	controller.Index(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestArticleController_IndexWithFormSearch(t *testing.T) {
	mockArticles := []*model.Article{
		{Title: "Test1", Slug: "test1", Tags: []string{"tag1"}},
		{Title: "Test2", Slug: "test2", Tags: []string{"tag2"}},
	}
	repo := &mockArticleRepository{articles: mockArticles}
	controller := NewArticleController(repo)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Form = make(map[string][]string)
	req.Form.Set("search", "test")
	w := httptest.NewRecorder()

	controller.Index(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestArticleController_IndexWithTagFilter(t *testing.T) {
	mockArticles := []*model.Article{
		{Title: "Test1", Slug: "test1", Tags: []string{"tag1"}},
		{Title: "Test2", Slug: "test2", Tags: []string{"tag2"}},
	}
	repo := &mockArticleRepository{articles: mockArticles}
	controller := NewArticleController(repo)

	req := httptest.NewRequest(http.MethodGet, "/?tag=tag1", nil)
	w := httptest.NewRecorder()

	controller.Index(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestArticleController_ArticleFound(t *testing.T) {
	mockArticles := []*model.Article{
		{Title: "Test1", Slug: "test1"},
	}
	repo := &mockArticleRepository{articles: mockArticles}
	controller := NewArticleController(repo)

	req := httptest.NewRequest(http.MethodGet, "/article/test1", nil)
	req.SetPathValue("slug", "test1")
	w := httptest.NewRecorder()

	controller.Article(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestArticleController_ArticleNotFound(t *testing.T) {
	mockArticles := []*model.Article{
		{Title: "Test1", Slug: "test1"},
	}
	repo := &mockArticleRepository{articles: mockArticles}
	controller := NewArticleController(repo)

	req := httptest.NewRequest(http.MethodGet, "/article/nonexistent", nil)
	req.SetPathValue("slug", "nonexistent")
	w := httptest.NewRecorder()

	controller.Article(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestArticleController_ArticleError(t *testing.T) {
	repo := &mockArticleRepository{err: errors.New("db error")}
	controller := NewArticleController(repo)

	req := httptest.NewRequest(http.MethodGet, "/article/test1", nil)
	req.SetPathValue("slug", "test1")
	w := httptest.NewRecorder()

	controller.Article(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestArticleController_ArchiveSuccessful(t *testing.T) {
	mockArticles := []*model.Article{
		{Title: "Test1", Slug: "test1"},
		{Title: "Test2", Slug: "test2"},
	}
	mockTagInfo := map[string]int{"tag1": 1, "tag2": 2}

	repo := &mockArticleRepository{
		articles: mockArticles,
		tagInfo:  mockTagInfo,
	}
	controller := NewArticleController(repo)

	req := httptest.NewRequest(http.MethodGet, "/archive", nil)
	w := httptest.NewRecorder()

	controller.Archive(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestArticleController_ArchiveEmpty(t *testing.T) {
	repo := &mockArticleRepository{
		articles: []*model.Article{},
		tagInfo:  map[string]int{},
	}
	controller := NewArticleController(repo)

	req := httptest.NewRequest(http.MethodGet, "/archive", nil)
	w := httptest.NewRecorder()

	controller.Archive(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestArticleController_RenderAndPageArticlesDefault(t *testing.T) {
	articles := make([]*model.Article, 25)
	for i := 0; i < 25; i++ {
		articles[i] = &model.Article{Title: "Test" + string(rune(i))}
	}

	controller := NewArticleController(&mockArticleRepository{})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	controller.renderAndPageArticles(w, req, articles, "", "")
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestArticleController_RenderAndPageArticlesSecondPage(t *testing.T) {
	articles := make([]*model.Article, 25)
	for i := 0; i < 25; i++ {
		articles[i] = &model.Article{Title: "Test" + string(rune(i))}
	}

	controller := NewArticleController(&mockArticleRepository{})
	req := httptest.NewRequest(http.MethodGet, "/?page=2", nil)
	w := httptest.NewRecorder()

	controller.renderAndPageArticles(w, req, articles, "", "")
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestArticleController_RenderAndPageArticlesCustomPerPage(t *testing.T) {
	articles := make([]*model.Article, 25)
	for i := 0; i < 25; i++ {
		articles[i] = &model.Article{Title: "Test" + string(rune(i))}
	}

	controller := NewArticleController(&mockArticleRepository{})
	req := httptest.NewRequest(http.MethodGet, "/?per_page=10", nil)
	w := httptest.NewRecorder()

	controller.renderAndPageArticles(w, req, articles, "", "")
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestArticleController_RenderAndPageArticlesInvalidPage(t *testing.T) {
	articles := make([]*model.Article, 25)
	for i := 0; i < 25; i++ {
		articles[i] = &model.Article{Title: "Test" + string(rune(i))}
	}

	controller := NewArticleController(&mockArticleRepository{})
	req := httptest.NewRequest(http.MethodGet, "/?page=invalid", nil)
	w := httptest.NewRecorder()

	controller.renderAndPageArticles(w, req, articles, "", "")
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestArticleController_RenderAndPageArticlesExceedMaxPerPage(t *testing.T) {
	articles := make([]*model.Article, 25)
	for i := 0; i < 25; i++ {
		articles[i] = &model.Article{Title: "Test" + string(rune(i))}
	}

	controller := NewArticleController(&mockArticleRepository{})
	req := httptest.NewRequest(http.MethodGet, "/?per_page=100", nil)
	w := httptest.NewRecorder()

	controller.renderAndPageArticles(w, req, articles, "", "")
	assert.Equal(t, http.StatusOK, w.Code)
}
