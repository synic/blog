package route

import (
	"net/http"

	"github.com/a-h/templ"

	"github.com/synic/adamthings.me/internal/web/component"
)

type defaultRouter struct {
	assets http.Handler
}

func NewDefaultRouter(assets http.Handler) defaultRouter {
	return defaultRouter{assets: assets}
}

func (h defaultRouter) Mount(server *http.ServeMux) {
	server.Handle("GET /static/", h.assets)
	server.Handle("GET /about/{$}", templ.Handler(component.AboutView()))
}
