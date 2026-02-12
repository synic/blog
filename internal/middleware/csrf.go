package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

func CSRFMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Ensure a CSRF cookie exists on every response
			if _, err := r.Cookie("csrf_token"); err != nil {
				token := generateCSRFToken()
				http.SetCookie(w, &http.Cookie{
					Name:     "csrf_token",
					Value:    token,
					Path:     "/",
					HttpOnly: false, // JS needs to read it
					SameSite: http.SameSiteLaxMode,
				})
			}

			// Validate on mutating methods
			if r.Method == http.MethodPost || r.Method == http.MethodPut ||
				r.Method == http.MethodDelete || r.Method == http.MethodPatch {

				cookie, err := r.Cookie("csrf_token")
				if err != nil {
					http.Error(w, "CSRF token missing", http.StatusForbidden)
					return
				}

				header := r.Header.Get("X-CSRF-Token")
				if header == "" || header != cookie.Value {
					http.Error(w, "CSRF token invalid", http.StatusForbidden)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func generateCSRFToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
