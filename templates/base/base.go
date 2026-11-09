package base

import (
	_ "embed"
	"fmt"

	"github.com/konnco/go-bootstrap/internal/renderer"
)

//go:embed cmd_serve_main.go.tmpl
var cmdServeMainTmpl string

//go:embed app.go.tmpl
var appTmpl string

//go:embed lifecycle.go.tmpl
var lifecycleTmpl string

//go:embed config.go.tmpl
var configTmpl string

//go:embed routes.go.tmpl
var routesTmpl string

//go:embed example_handler.go.tmpl
var exampleHandlerTmpl string

//go:embed config.yaml.tmpl
var configYamlTmpl string

//go:embed env.example.tmpl
var envExampleTmpl string

//go:embed gitignore.tmpl
var gitignoreTmpl string

//go:embed readme.tmpl
var readmeTmpl string

type Options struct {
	ProjectName   string
	ModulePath    string
	HasDatabase   bool
	DatabaseType  string
	HTTPFramework string
	WorkerType    string
	HasWorker     bool
}

type File struct {
	Path    string
	Content []byte
}

// Render renders all base project files.
func Render(opts Options) ([]File, error) {
	templates := []struct {
		destPath string
		tmplStr  string
	}{
		{"cmd/serve/main.go", cmdServeMainTmpl},
		{"internal/bootstrap/app.go", appTmpl},
		{"internal/bootstrap/lifecycle.go", lifecycleTmpl},
		{"internal/config/config.go", configTmpl},
		{"internal/routes/routes.go", routesTmpl},
		{"modules/example/handler.go", exampleHandlerTmpl},
		{"configs/config.yaml", configYamlTmpl},
		{".env.example", envExampleTmpl},
		{".gitignore", gitignoreTmpl},
		{"README.md", readmeTmpl},
	}

	var result []File
	for _, t := range templates {
		content, err := renderer.RenderString(t.destPath, t.tmplStr, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to render %s: %w", t.destPath, err)
		}
		result = append(result, File{
			Path:    t.destPath,
			Content: content,
		})
	}

	return result, nil
}
