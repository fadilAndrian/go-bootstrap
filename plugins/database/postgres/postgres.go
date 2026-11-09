package postgres

import (
	_ "embed"
	"fmt"

	"github.com/konnco/go-bootstrap/internal/core"
	"github.com/konnco/go-bootstrap/internal/renderer"
)

//go:embed postgres.go.tmpl
var postgresTmpl string

type Plugin struct{}

func init() {
	core.Register(&Plugin{})
}

func (p *Plugin) Metadata() core.PluginMetadata {
	return core.PluginMetadata{
		Name:        "postgres",
		Category:    core.CategoryDatabase,
		Version:     "1.0.0",
		Description: "PostgreSQL adapter using pgx/v5 connection pool",
	}
}

func (p *Plugin) Category() core.Category {
	return core.CategoryDatabase
}

func (p *Plugin) Dependencies() []string {
	return []string{
		"github.com/jackc/pgx/v5@v5.7.2",
	}
}

func (p *Plugin) Validate(ctx *core.GenerateContext) error {
	return nil
}

func (p *Plugin) Generate(ctx *core.GenerateContext) ([]core.GeneratedFile, error) {
	data := map[string]any{
		"ModulePath": ctx.ModulePath,
	}

	content, err := renderer.RenderString("internal/database/postgres.go", postgresTmpl, data)
	if err != nil {
		return nil, fmt.Errorf("failed to render postgres adapter: %w", err)
	}

	return []core.GeneratedFile{
		{
			Path:    "internal/database/postgres.go",
			Content: content,
		},
	}, nil
}
