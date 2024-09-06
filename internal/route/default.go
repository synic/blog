package route

import (
	"net/http"

	"github.com/a-h/templ"

	"github.com/synic/adamthings.me/internal/route/component"
)

type defaultRouter struct {
	assetsDir string
}

func NewDefaultRouter(assetsDir string) defaultRouter {
	return defaultRouter{assetsDir: assetsDir}
}

func (h defaultRouter) Mount(server *http.ServeMux) {
	fs := http.FileServer(http.Dir(h.assetsDir))
	server.Handle("GET /static/", http.StripPrefix("/static/", fs))
	server.Handle("GET /about/{$}", templ.Handler(component.AboutView()))
}
