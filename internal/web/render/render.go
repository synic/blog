package render

import (
	"net/http"

	"github.com/a-h/templ"

	"github.com/synic/adamthings.me/internal/web/view"
)

type renderTemplConfig struct {
	status int
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

func Templ(
	w http.ResponseWriter,
	r *http.Request,
	c templ.Component,
	options ...func(*renderTemplConfig),
) {
	conf := defaultRenderTemplConfig()

	for _, option := range options {
		option(&conf)
	}

	w.WriteHeader(conf.status)
	c.Render(r.Context(), w)
}

func Error(w http.ResponseWriter, r *http.Request, status int, title, message string) {
	Templ(w, r, view.ErrorView(title, message), WithStatus(status))
}
