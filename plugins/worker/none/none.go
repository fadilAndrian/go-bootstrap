package none

import (
	"github.com/konnco/go-bootstrap/internal/core"
)

type Plugin struct{}

func init() {
	core.Register(&Plugin{})
}

func (p *Plugin) Metadata() core.PluginMetadata {
	return core.PluginMetadata{
		Name:        "none",
		Category:    core.CategoryWorker,
		Version:     "1.0.0",
		Description: "No background worker",
	}
}

func (p *Plugin) Category() core.Category {
	return core.CategoryWorker
}

func (p *Plugin) Dependencies() []string {
	return nil
}

func (p *Plugin) Validate(ctx *core.GenerateContext) error {
	return nil
}

func (p *Plugin) Generate(ctx *core.GenerateContext) ([]core.GeneratedFile, error) {
	// Worker 'none' generates no worker files
	return nil, nil
}
