package nethttp

import (
	"fmt"

	"github.com/fadilAndrian/go-bootstrap/internal/core"
)

type Plugin struct{}

func init() {
	core.Register(&Plugin{})
}

func (p *Plugin) Metadata() core.PluginMetadata {
	return core.PluginMetadata{
		Name:        "nethttp",
		Category:    core.CategoryHTTP,
		Version:     "1.0.0",
		Description: "Standard library net/http ServeMux",
	}
}

func (p *Plugin) Category() core.Category {
	return core.CategoryHTTP
}

func (p *Plugin) Dependencies() []string {
	return nil
}

func (p *Plugin) Validate(ctx *core.GenerateContext) error {
	return nil
}

func (p *Plugin) Generate(ctx *core.GenerateContext) ([]core.GeneratedFile, error) {
	muxContent := []byte(`package http

import (
	"net/http"
)

func NewRouter() *http.ServeMux {
	return http.NewServeMux()
}
`)

	serverContent := []byte(fmt.Sprintf(`package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"%s/internal/config"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg config.ServerConfig, handler http.Handler) *Server {
	addr := fmt.Sprintf("%%s:%%s", cfg.Host, cfg.Port)
	return &Server{
		httpServer: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
`, ctx.ModulePath))

	return []core.GeneratedFile{
		{
			Path:    "internal/http/server.go",
			Content: serverContent,
		},
		{
			Path:    "internal/http/nethttp.go",
			Content: muxContent,
		},
	}, nil
}
