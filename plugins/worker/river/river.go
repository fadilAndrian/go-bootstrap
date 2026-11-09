package river

import (
	_ "embed"
	"fmt"

	"github.com/fadilAndrian/go-bootstrap/internal/core"
	"github.com/fadilAndrian/go-bootstrap/internal/renderer"
)

//go:embed river.go.tmpl
var riverTmpl string

//go:embed job.go.tmpl
var jobTmpl string

type Plugin struct{}

func init() {
	core.Register(&Plugin{})
}

func (p *Plugin) Metadata() core.PluginMetadata {
	return core.PluginMetadata{
		Name:        "river",
		Category:    core.CategoryWorker,
		Version:     "1.0.0",
		Description: "River PostgreSQL-backed background job queue",
	}
}

func (p *Plugin) Category() core.Category {
	return core.CategoryWorker
}

func (p *Plugin) Dependencies() []string {
	return []string{
		"github.com/riverqueue/river@latest",
		"github.com/riverqueue/river/riverdriver/riverpgxv5@latest",
		"github.com/riverqueue/river/rivermigrate@latest",
	}
}

func (p *Plugin) Validate(ctx *core.GenerateContext) error {
	db := ctx.Plugins[core.CategoryDatabase]
	if db != "postgres" {
		return fmt.Errorf("plugin 'river' requires 'postgres' database adapter, but %q was configured", db)
	}
	return nil
}

func (p *Plugin) Generate(ctx *core.GenerateContext) ([]core.GeneratedFile, error) {
	data := map[string]any{
		"ModulePath": ctx.ModulePath,
	}

	riverContent, err := renderer.RenderString("internal/worker/river.go", riverTmpl, data)
	if err != nil {
		return nil, fmt.Errorf("failed to render river worker: %w", err)
	}

	jobContent, err := renderer.RenderString("internal/worker/job.go", jobTmpl, data)
	if err != nil {
		return nil, fmt.Errorf("failed to render river job: %w", err)
	}

	return []core.GeneratedFile{
		{
			Path:    "internal/worker/river.go",
			Content: riverContent,
		},
		{
			Path:    "internal/worker/job.go",
			Content: jobContent,
		},
	}, nil
}
