package chi

import (
	_ "embed"
	"fmt"

	"github.com/fadilAndrian/go-bootstrap/internal/core"
	"github.com/fadilAndrian/go-bootstrap/internal/renderer"
)

//go:embed chi.go.tmpl
var chiTmpl string

type Plugin struct{}

func init() {
	core.Register(&Plugin{})
}

func (p *Plugin) Metadata() core.PluginMetadata {
	return core.PluginMetadata{
		Name:        "chi",
		Category:    core.CategoryHTTP,
		Version:     "1.0.0",
		Description: "Chi lightweight, idiomatic, and composable router",
	}
}

func (p *Plugin) Category() core.Category {
	return core.CategoryHTTP
}

func (p *Plugin) Dependencies() []string {
	return []string{
		"github.com/go-chi/chi/v5@v5.2.1",
	}
}

func (p *Plugin) Validate(ctx *core.GenerateContext) error {
	return nil
}

func (p *Plugin) Generate(ctx *core.GenerateContext) ([]core.GeneratedFile, error) {
	data := map[string]any{
		"ModulePath": ctx.ModulePath,
	}

	chiContent, err := renderer.RenderString("internal/http/chi.go", chiTmpl, data)
	if err != nil {
		return nil, fmt.Errorf("failed to render chi router: %w", err)
	}

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
			Path:    "internal/http/chi.go",
			Content: chiContent,
		},
	}, nil
}
