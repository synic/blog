package view

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/starfederation/datastar-go/datastar"
	"github.com/synic/blog/internal/middleware"
)

type renderTemplConfig struct {
	status int
	title  string
}

func defaultRenderTemplConfig() renderTemplConfig {
	return renderTemplConfig{
		status: 200,
	}
}

func WithStatus(status int) func(*renderTemplConfig) {
	return func(conf *renderTemplConfig) {
		conf.status = status
	}
}

func WithTitle(title string) func(*renderTemplConfig) {
	return func(conf *renderTemplConfig) {
		conf.title = title
	}
}

func isDatastarRequest(r *http.Request) bool {
	if v, ok := r.Context().Value(middleware.DatastarPartialContextKey).(bool); ok {
		return v
	}
	return false
}

func Render(
	w http.ResponseWriter,
	r *http.Request,
	comp templ.Component,
	options ...func(*renderTemplConfig),
) {
	conf := defaultRenderTemplConfig()

	for _, option := range options {
		option(&conf)
	}

	if isDatastarRequest(r) {
		var buf bytes.Buffer
		if err := comp.Render(r.Context(), &buf); err != nil {
			log.Printf("Error rendering component: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		sse := datastar.NewSSE(w, r)

		sse.PatchElements(buf.String())

		if r.Method == http.MethodGet {
			sse.ExecuteScript(
				fmt.Sprintf("history.pushState({}, '', %q);", r.URL.RequestURI()),
			)
		}

		return
	} else {
		w.WriteHeader(conf.status)
		comp.Render(r.Context(), w)
	}
}

func Error(
	w http.ResponseWriter,
	r *http.Request,
	err error,
	status int,
	title, message string,
) {
	if err != nil {
		log.Printf("Error during request at: %s: %v", r.URL.Path, err)
	}
	Render(w, r, ErrorView(title, message), WithStatus(status), WithTitle(title))
}
