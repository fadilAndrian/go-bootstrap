package core

type Category string

const (
	CategoryDatabase Category = "database"
	CategoryHTTP     Category = "http"
	CategoryWorker   Category = "worker"
)

type PluginMetadata struct {
	Name        string   `json:"name" yaml:"name"`
	Category    Category `json:"category" yaml:"category"`
	Version     string   `json:"version" yaml:"version"`
	Description string   `json:"description" yaml:"description"`
}

type GeneratedFile struct {
	Path    string // Relative to target directory
	Content []byte
}

type GenerateContext struct {
	ProjectName string
	ModulePath  string
	TargetDir   string
	Plugins     map[Category]string
	Data        map[string]any
}

type Plugin interface {
	Metadata() PluginMetadata
	Category() Category
	Dependencies() []string // Go module dependencies e.g. "github.com/gin-gonic/gin@v1.10.0"
	Validate(ctx *GenerateContext) error
	Generate(ctx *GenerateContext) ([]GeneratedFile, error)
}
