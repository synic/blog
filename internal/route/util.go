package route

import (
	"net/http"

	"github.com/a-h/templ"
)

type renderConfig struct {
	status int
}

func getDefaultConfig() renderConfig {
	return renderConfig{
		status: 200,
	}
}

func RenderTempl(
	w http.ResponseWriter,
	r *http.Request,
	c templ.Component,
	options ...func(*renderConfig),
) {
	conf := getDefaultConfig()

	for _, option := range options {
		option(&conf)
	}

	w.WriteHeader(conf.status)
	c.Render(r.Context(), w)
}
