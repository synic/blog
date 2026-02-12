package mail

import (
	"fmt"
	"log"

	"github.com/resend/resend-go/v2"
	"github.com/synic/blog/internal/config"
)

type Mailer struct {
	cfg    config.Config
	client *resend.Client
}

func NewMailer(cfg config.Config) *Mailer {
	var client *resend.Client
	if cfg.ResendAPIKey != "" {
		client = resend.NewClient(cfg.ResendAPIKey)
	}
	return &Mailer{cfg: cfg, client: client}
}

func (m *Mailer) Send(to, subject, body string) {
	if m.client == nil {
		log.Printf("Resend not configured, skipping email to %s: %s", to, subject)
		return
	}

	go func() {
		params := &resend.SendEmailRequest{
			From:    "Blog <noreply@synic.dev>",
			To:      []string{to},
			Subject: subject,
			Text:    body,
		}

		_, err := m.client.Emails.Send(params)
		if err != nil {
			log.Printf("Failed to send email to %s: %v", to, err)
		}
	}()
}

func (m *Mailer) unsubscribeFooter(unsubscribeToken string) string {
	return fmt.Sprintf(
		"\n\n---\nUnsubscribe: %s/unsubscribe?token=%s",
		m.cfg.ServerAddress, unsubscribeToken,
	)
}

func (m *Mailer) NotifyPendingComment(commentID int64, articleSlug, articleURL, username, body, unsubscribeToken string) {
	if m.cfg.AdminEmail == "" {
		return
	}

	subject := fmt.Sprintf("New comment pending approval on %s", articleSlug)
	text := fmt.Sprintf(
		"%s commented on %s:\n\n%s\n\nApprove: %s/admin/comments/%d/approve\nDelete: %s/admin/comments/%d/delete",
		username, articleSlug, body,
		m.cfg.ServerAddress, commentID,
		m.cfg.ServerAddress, commentID,
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
		articleSlug, m.cfg.ServerAddress+articleURL+"?show_comments=1",
	)
	text += m.unsubscribeFooter(unsubscribeToken)

	m.Send(toEmail, subject, text)
}

func (m *Mailer) NotifyReply(toEmail, articleSlug, articleURL, replyUsername, replyBody, unsubscribeToken string) {
	if toEmail == "" {
		return
	}

	subject := fmt.Sprintf("%s replied to your comment on %s", replyUsername, articleSlug)
	text := fmt.Sprintf(
		"%s replied to your comment on %s:\n\n%s\n\nView: %s",
		replyUsername, articleSlug, replyBody, m.cfg.ServerAddress+articleURL+"?show_comments=1",
	)
	text += m.unsubscribeFooter(unsubscribeToken)

	m.Send(toEmail, subject, text)
}
