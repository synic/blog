package controller

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/synic/blog/internal/config"
	"github.com/synic/blog/internal/db"
	"github.com/synic/blog/internal/mail"
	"github.com/synic/blog/internal/middleware"
	"github.com/synic/blog/internal/model"
	"github.com/synic/blog/internal/store"
	"github.com/synic/blog/internal/view"
)

type CommentController struct {
	comments *store.CommentRepository
	articles store.ArticleRepository
	queries  *db.Queries
	mailer   *mail.Mailer
	config   config.Config
}

func NewCommentController(
	comments *store.CommentRepository,
	articles store.ArticleRepository,
	queries *db.Queries,
	mailer *mail.Mailer,
	cfg config.Config,
) CommentController {
	return CommentController{
		comments: comments,
		articles: articles,
		queries:  queries,
		mailer:   mailer,
		config:   cfg,
	}
}

func (c CommentController) List(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	if _, err := c.articles.FindOneBySlug(r.Context(), slug); err != nil {
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}

	comments, err := c.comments.ListBySlug(r.Context(), slug)
	if err != nil {
		http.Error(w, "Failed to load comments", http.StatusInternalServerError)
		return
	}

	user := middleware.UserFromContext(r.Context())
	articleURL := "/article/" + r.PathValue("date") + "/" + slug
	threads := model.OrganizeComments(comments)
	view.Render(w, r, view.CommentList(articleURL, threads, user))
}

func (c CommentController) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	slug := r.PathValue("slug")

	article, err := c.articles.FindOneBySlug(r.Context(), slug)
	if err != nil {
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}

	body := strings.TrimSpace(r.FormValue("body"))
	if body == "" {
		http.Error(w, "Comment body is required", http.StatusBadRequest)
		return
	}
	if len(body) > 5000 {
		http.Error(w, "Comment is too long (max 5000 characters)", http.StatusBadRequest)
		return
	}

	var parentID *int64
	if pidStr := r.FormValue("parent_id"); pidStr != "" {
		pid, err := strconv.ParseInt(pidStr, 10, 64)
		if err == nil {
			parentID = &pid
		}
	}

	approved := user.IsAdmin
	comment, err := c.comments.Create(r.Context(), slug, user.ID, body, parentID, approved)
	if err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	articleURL := "/article/" + r.PathValue("date") + "/" + slug

	if approved {
		if parentID != nil {
			parentComment, err := c.comments.GetCommentWithUser(r.Context(), *parentID)
			if err == nil && parentComment.Email != "" && !parentComment.Unsubscribed {
				c.mailer.NotifyReply(
					parentComment.Email,
					article.Slug,
					articleURL,
					user.Username,
					body,
					parentComment.UnsubscribeToken,
				)
			} else if err != nil {
				log.Printf(
					"Failed to look up parent comment %d for reply notification: %v",
					*parentID,
					err,
				)
			}
		}

		comments, err := c.comments.ListBySlug(r.Context(), slug)
		if err != nil {
			http.Error(w, "Failed to load comments", http.StatusInternalServerError)
			return
		}
		threads := model.OrganizeComments(comments)
		view.Render(w, r, view.CommentList(articleURL, threads, user))
	} else {
		adminUser, err := c.queries.GetUserByEmail(r.Context(), c.config.AdminEmail)
		if err == nil && !adminUser.Unsubscribed {
			c.mailer.NotifyPendingComment(
				comment.ID,
				article.Slug,
				articleURL,
				user.Username,
				body,
				adminUser.UnsubscribeToken,
			)
		} else if err != nil {
			log.Printf("Failed to look up admin user for pending notification: %v", err)
		}

		comments, err := c.comments.ListBySlug(r.Context(), slug)
		if err != nil {
			http.Error(w, "Failed to load comments", http.StatusInternalServerError)
			return
		}
		threads := model.OrganizeComments(comments)
		view.Render(w, r, view.CommentListWithPending(articleURL, threads, user))
	}
}

func (c CommentController) Approve(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/auth/login?return_to="+url.QueryEscape(r.URL.Path), http.StatusFound)
		return
	}
	if !user.IsAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	existing, err := c.comments.GetCommentWithUser(r.Context(), id)
	if err != nil {
		view.Error(w, r, nil, http.StatusNotFound, "Not Found", "This comment no longer exists.")
		return
	}
	if existing.Approved {
		view.Error(w, r, nil, http.StatusConflict, "Already Approved", "This comment has already been approved.")
		return
	}

	approvedComment, err := c.comments.Approve(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to approve comment", http.StatusInternalServerError)
		return
	}

	commentWithUser, err := c.comments.GetCommentWithUser(r.Context(), id)
	if err != nil {
		log.Printf("Failed to look up approved comment %d: %v", id, err)
	} else {
		article, err := c.articles.FindOneBySlug(r.Context(), approvedComment.ArticleSlug)
		if err != nil {
			log.Printf("Failed to look up article %s: %v", approvedComment.ArticleSlug, err)
		} else {
			articleURL := "/article/" + article.PublishedAt.Format(
				"2006-01-02",
			) + "/" + article.Slug

			if commentWithUser.Email != "" && !commentWithUser.Unsubscribed {
				c.mailer.NotifyCommentApproved(
					commentWithUser.Email,
					article.Slug,
					articleURL,
					commentWithUser.UnsubscribeToken,
				)
			}

			if approvedComment.ParentID.Valid {
				parentComment, err := c.comments.GetCommentWithUser(
					r.Context(),
					approvedComment.ParentID.Int64,
				)
				if err == nil && parentComment.Email != "" && !parentComment.Unsubscribed {
					c.mailer.NotifyReply(
						parentComment.Email,
						article.Slug,
						articleURL,
						commentWithUser.Username,
						commentWithUser.Body,
						parentComment.UnsubscribeToken,
					)
				} else if err != nil {
					log.Printf(
						"Failed to look up parent comment %d for reply notification: %v",
						approvedComment.ParentID.Int64,
						err,
					)
				}
			}

			http.Redirect(
				w,
				r,
				c.config.ServerAddress+articleURL+"?show_comments=1",
				http.StatusFound,
			)
			return
		}
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (c CommentController) Delete(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/auth/login?return_to="+url.QueryEscape(r.URL.Path), http.StatusFound)
		return
	}
	if !user.IsAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	if _, err := c.comments.GetCommentWithUser(r.Context(), id); err != nil {
		view.Error(w, r, nil, http.StatusNotFound, "Not Found", "This comment has already been deleted.")
		return
	}

	if err := c.comments.Delete(r.Context(), id); err != nil {
		http.Error(w, "Failed to delete comment", http.StatusInternalServerError)
		return
	}

	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/"
	}
	http.Redirect(w, r, referer, http.StatusFound)
}
