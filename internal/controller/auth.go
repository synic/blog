package controller

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/synic/blog/internal/config"
	"github.com/synic/blog/internal/db"
	"github.com/synic/blog/internal/view"
)

type AuthController struct {
	queries *db.Queries
	config  config.Config
}

func NewAuthController(queries *db.Queries, cfg config.Config) AuthController {
	return AuthController{queries: queries, config: cfg}
}

func (c AuthController) Login(w http.ResponseWriter, r *http.Request) {
	state := generateToken()

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		MaxAge:   600,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	returnTo := r.URL.Query().Get("return_to")
	if returnTo == "" {
		returnTo = r.Header.Get("Referer")
	}
	if returnTo == "" {
		returnTo = "/"
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "return_to",
		Value:    returnTo,
		Path:     "/",
		MaxAge:   600,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	u := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&state=%s&scope=%s",
		url.QueryEscape(c.config.GitHubClientID),
		url.QueryEscape(c.config.ServerAddress+"/auth/callback"),
		url.QueryEscape(state),
		url.QueryEscape("read:user user:email"),
	)

	http.Redirect(w, r, u, http.StatusFound)
}

func (c AuthController) Callback(w http.ResponseWriter, r *http.Request) {
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil || stateCookie.Value != r.URL.Query().Get("state") {
		http.Error(w, "Invalid OAuth state", http.StatusForbidden)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "oauth_state",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing code", http.StatusBadRequest)
		return
	}

	accessToken, err := c.exchangeCode(code)
	if err != nil {
		http.Error(w, "Failed to exchange code", http.StatusInternalServerError)
		return
	}

	ghUser, err := c.fetchGitHubUser(accessToken)
	if err != nil {
		http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
		return
	}

	email := c.fetchGitHubEmail(accessToken)

	dbUser, err := c.queries.UpsertUser(r.Context(), db.UpsertUserParams{
		GithubID:  ghUser.ID,
		Username:  ghUser.Login,
		AvatarUrl: ghUser.AvatarURL,
		Email:     email,
	})
	if err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	sessionToken := generateToken()
	csrfToken := generateToken()

	_, err = c.queries.CreateSession(r.Context(), db.CreateSessionParams{
		UserID:    dbUser.ID,
		Token:     sessionToken,
		CsrfToken: csrfToken,
		ExpiresAt: pgtype.Timestamptz{Time: time.Now().Add(30 * 24 * time.Hour), Valid: true},
	})
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		MaxAge:   30 * 24 * 60 * 60,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Path:     "/",
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
	})

	returnTo := "/"
	if cookie, err := r.Cookie("return_to"); err == nil {
		returnTo = cookie.Value
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "return_to",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	http.Redirect(w, r, appendShowComments(returnTo), http.StatusFound)
}

func (c AuthController) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		view.Error(w, r, nil, http.StatusBadRequest, "Invalid Link", "This unsubscribe link is invalid.")
		return
	}

	result, err := c.queries.UnsubscribeUser(r.Context(), token)
	if err != nil || result == 0 {
		view.Error(w, r, err, http.StatusNotFound, "Not Found", "This unsubscribe link is invalid or has already been used.")
		return
	}

	view.Render(w, r, view.UnsubscribeView())
}

func (c AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("session_token"); err == nil {
		c.queries.DeleteSession(r.Context(), cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	returnTo := r.Header.Get("Referer")
	if returnTo == "" {
		returnTo = "/"
	}

	returnTo = appendShowComments(returnTo)

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", returnTo)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	http.Redirect(w, r, returnTo, http.StatusFound)
}

func appendShowComments(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	q := u.Query()
	q.Set("show_comments", "1")
	u.RawQuery = q.Encode()
	return u.String()
}

func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

type githubTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func (c AuthController) exchangeCode(code string) (string, error) {
	data := url.Values{
		"client_id":     {c.config.GitHubClientID},
		"client_secret": {c.config.GitHubClientSecret},
		"code":          {code},
	}

	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", nil)
	if err != nil {
		return "", err
	}
	req.URL.RawQuery = data.Encode()
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResp githubTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}

type githubUser struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
}

func (c AuthController) fetchGitHubUser(accessToken string) (githubUser, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return githubUser{}, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return githubUser{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return githubUser{}, err
	}

	var user githubUser
	if err := json.Unmarshal(body, &user); err != nil {
		return githubUser{}, err
	}

	return user, nil
}

type githubEmail struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

func (c AuthController) fetchGitHubEmail(accessToken string) string {
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return ""
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	var emails []githubEmail
	if err := json.Unmarshal(body, &emails); err != nil {
		log.Printf("Failed to parse GitHub emails: %v", err)
		return ""
	}

	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email
		}
	}

	return ""
}
