package store

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/synic/blog/internal/db"
	"github.com/synic/blog/internal/model"
)

type CommentRepository struct {
	queries *db.Queries
	counts  map[string]int
	mu      sync.RWMutex
}

func NewCommentRepository(queries *db.Queries) *CommentRepository {
	return &CommentRepository{
		queries: queries,
		counts:  make(map[string]int),
	}
}

func (r *CommentRepository) LoadCounts(ctx context.Context) error {
	rows, err := r.queries.CountCommentsBySlug(ctx)
	if err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, row := range rows {
		r.counts[row.ArticleSlug] = int(row.CommentCount)
	}

	return nil
}

func (r *CommentRepository) CommentCount(slug string) int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.counts[slug]
}

func (r *CommentRepository) Create(
	ctx context.Context,
	slug string,
	userID int64,
	body string,
	parentID *int64,
	approved bool,
) (db.Comment, error) {
	params := db.CreateCommentParams{
		ArticleSlug: slug,
		UserID:      userID,
		Body:        body,
		Approved:    approved,
	}
	if parentID != nil {
		params.ParentID = pgtype.Int8{Int64: *parentID, Valid: true}
	}

	comment, err := r.queries.CreateComment(ctx, params)
	if err != nil {
		return comment, err
	}

	if approved {
		r.mu.Lock()
		r.counts[slug]++
		r.mu.Unlock()
	}

	return comment, nil
}

func (r *CommentRepository) Approve(ctx context.Context, id int64) (db.Comment, error) {
	comment, err := r.queries.ApproveComment(ctx, id)
	if err != nil {
		return comment, err
	}

	r.mu.Lock()
	r.counts[comment.ArticleSlug]++
	r.mu.Unlock()

	return comment, nil
}

func (r *CommentRepository) Delete(ctx context.Context, id int64) error {
	comment, err := r.queries.GetCommentWithUser(ctx, id)
	if err != nil {
		return err
	}

	if err := r.queries.DeleteComment(ctx, id); err != nil {
		return err
	}

	if comment.Approved {
		r.mu.Lock()
		r.counts[comment.ArticleSlug]--
		r.mu.Unlock()
	}

	return nil
}

func (r *CommentRepository) ListBySlug(
	ctx context.Context,
	slug string,
) ([]model.Comment, error) {
	rows, err := r.queries.ListCommentsBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	comments := make([]model.Comment, len(rows))
	for i, row := range rows {
		var parentID *int64
		if row.ParentID.Valid {
			parentID = &row.ParentID.Int64
		}
		comments[i] = model.Comment{
			ID:        row.ID,
			Body:      row.Body,
			Username:  row.Username,
			AvatarURL: row.AvatarUrl,
			CreatedAt: row.CreatedAt.Time,
			ParentID:  parentID,
		}
	}

	return comments, nil
}

func (r *CommentRepository) GetCommentWithUser(
	ctx context.Context,
	id int64,
) (db.GetCommentWithUserRow, error) {
	return r.queries.GetCommentWithUser(ctx, id)
}
