package mail

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/resend/resend-go/v2"

	"github.com/synic/blog/internal/config"
)

type emailTask struct {
	to      string
	subject string
	body    string
}

type Mailer struct {
	cfg    config.Config
	client *resend.Client
	queue  chan emailTask
}

func NewMailer(cfg config.Config) *Mailer {
	var client *resend.Client
	if cfg.ResendAPIKey != "" {
		httpClient := &http.Client{
			Timeout: 15 * time.Second,
		}
		client = resend.NewCustomClient(httpClient, cfg.ResendAPIKey)
	}
	m := &Mailer{
		cfg:    cfg,
		client: client,
		queue:  make(chan emailTask, 128),
	}
	if client != nil {
		go m.worker()
	}
	return m
}

func (m *Mailer) worker() {
	for task := range m.queue {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		params := &resend.SendEmailRequest{
			From:    "Blog <noreply@synic.dev>",
			To:      []string{task.to},
			Subject: task.subject,
			Text:    task.body,
		}

		_, err := m.client.Emails.SendWithContext(ctx, params)
		cancel()
		if err != nil {
			log.Printf("Failed to send email to %s: %v", task.to, err)
		}
	}
}

func (m *Mailer) Send(to, subject, body string) {
	if m.client == nil {
		log.Printf("Resend not configured, skipping email to %s: %s", to, subject)
		return
	}

	select {
	case m.queue <- emailTask{to: to, subject: subject, body: body}:
	default:
		log.Printf("Email queue full, dropping email to %s: %s", to, subject)
	}
}

func (m *Mailer) unsubscribeFooter(unsubscribeToken string) string {
	return fmt.Sprintf(
		"\n\n---\nUnsubscribe: %s/unsubscribe?token=%s",
		m.cfg.SiteUrl, unsubscribeToken,
	)
}

func (m *Mailer) NotifyPendingComment(
	commentID int64,
	articleSlug, articleURL, username, body, unsubscribeToken string,
) {
	if m.cfg.AdminEmail == "" {
		return
	}

	subject := fmt.Sprintf("New comment pending approval on %s", articleSlug)
	text := fmt.Sprintf(
		"%s commented on %s:\n\n%s\n\nApprove: %s/admin/comments/%d/approve\nDelete: %s/admin/comments/%d/delete",
		username,
		articleSlug,
		body,
		m.cfg.SiteUrl,
		commentID,
		m.cfg.SiteUrl,
		commentID,
	)
	text += m.unsubscribeFooter(unsubscribeToken)

	m.Send(m.cfg.AdminEmail, subject, text)
}

func (m *Mailer) NotifyCommentApproved(toEmail, articleSlug, articleURL, unsubscribeToken string) {
	if toEmail == "" {
		return
	}

	subject := fmt.Sprintf("Your comment on %s has been approved", articleSlug)
	text := fmt.Sprintf(
		"Your comment on %s has been approved.\n\nView: %s",
		articleSlug, m.cfg.SiteUrl+articleURL+"?show_comments=1",
	)
	text += m.unsubscribeFooter(unsubscribeToken)

	m.Send(toEmail, subject, text)
}

func (m *Mailer) NotifyReply(
	toEmail, articleSlug, articleURL, replyUsername, replyBody, unsubscribeToken string,
) {
	if toEmail == "" {
		return
	}

	subject := fmt.Sprintf("%s replied to your comment on %s", replyUsername, articleSlug)
	text := fmt.Sprintf(
		"%s replied to your comment on %s:\n\n%s\n\nView: %s",
		replyUsername, articleSlug, replyBody, m.cfg.SiteUrl+articleURL+"?show_comments=1",
	)
	text += m.unsubscribeFooter(unsubscribeToken)

	m.Send(toEmail, subject, text)
}
