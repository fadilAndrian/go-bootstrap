package gin

import (
	_ "embed"
	"fmt"

	"github.com/fadilAndrian/go-bootstrap/internal/core"
	"github.com/fadilAndrian/go-bootstrap/internal/renderer"
)

//go:embed server.go.tmpl
var serverTmpl string

//go:embed gin.go.tmpl
var ginTmpl string

type Plugin struct{}

func init() {
	core.Register(&Plugin{})
}

func (p *Plugin) Metadata() core.PluginMetadata {
	return core.PluginMetadata{
		Name:        "gin",
		Category:    core.CategoryHTTP,
		Version:     "1.0.0",
		Description: "Gin high-performance HTTP web framework",
	}
}

func (p *Plugin) Category() core.Category {
	return core.CategoryHTTP
}

func (p *Plugin) Dependencies() []string {
	return []string{
		"github.com/gin-gonic/gin@v1.10.0",
	}
}

func (p *Plugin) Validate(ctx *core.GenerateContext) error {
	return nil
}

func (p *Plugin) Generate(ctx *core.GenerateContext) ([]core.GeneratedFile, error) {
	data := map[string]any{
		"ModulePath": ctx.ModulePath,
	}

	serverContent, err := renderer.RenderString("internal/http/server.go", serverTmpl, data)
	if err != nil {
		return nil, fmt.Errorf("failed to render http server: %w", err)
	}

	ginContent, err := renderer.RenderString("internal/http/gin.go", ginTmpl, data)
	if err != nil {
		return nil, fmt.Errorf("failed to render gin router: %w", err)
	}

	return []core.GeneratedFile{
		{
			Path:    "internal/http/server.go",
			Content: serverContent,
		},
		{
			Path:    "internal/http/gin.go",
			Content: ginContent,
		},
	}, nil
}
