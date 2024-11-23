package view

import (
	"net/http"

	"github.com/a-h/templ"
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

	w.WriteHeader(conf.status)
	comp.Render(r.Context(), w)
}

func Error(
	w http.ResponseWriter,
	r *http.Request,
	err error,
	status int,
	title, message string,
) {
	Render(w, r, ErrorView(title, message), WithStatus(status))
}
