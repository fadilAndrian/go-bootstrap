package core

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/konnco/go-bootstrap/templates/base"
)

type GenerateOptions struct {
	ProjectName string
	ModulePath  string
	TargetDir   string
	Database    string // "postgres", "mysql", "sqlite", "none"
	HTTP        string // "gin", "chi", "nethttp"
	Worker      string // "river", "none"
	Overwrite   bool
	SkipTidy    bool
}

type Generator struct {
	registry *Registry
}

func NewGenerator(reg *Registry) *Generator {
	if reg == nil {
		reg = DefaultRegistry
	}
	return &Generator{registry: reg}
}

// Generate creates the application in TargetDir.
func (g *Generator) Generate(opts *GenerateOptions) error {
	if opts.TargetDir == "" {
		opts.TargetDir = "."
	}

	absTarget, err := filepath.Abs(opts.TargetDir)
	if err != nil {
		return fmt.Errorf("failed to resolve target directory: %w", err)
	}
	opts.TargetDir = absTarget

	if opts.ProjectName == "" {
		opts.ProjectName = filepath.Base(opts.TargetDir)
	}

	// Try detecting existing Go module if not explicitly provided
	if opts.ModulePath == "" {
		detectedModule, err := detectModulePath(opts.TargetDir)
		if err == nil && detectedModule != "" {
			opts.ModulePath = detectedModule
		} else {
			opts.ModulePath = opts.ProjectName
		}
	}

	// Set default components if not provided
	if opts.Database == "" {
		opts.Database = "postgres"
	}
	if opts.HTTP == "" {
		opts.HTTP = "gin"
	}
	if opts.Worker == "" {
		opts.Worker = "river"
	}

	ctx := &GenerateContext{
		ProjectName: opts.ProjectName,
		ModulePath:  opts.ModulePath,
		TargetDir:   opts.TargetDir,
		Plugins: map[Category]string{
			CategoryDatabase: opts.Database,
			CategoryHTTP:     opts.HTTP,
			CategoryWorker:   opts.Worker,
		},
		Data: make(map[string]any),
	}

	// 1. Validate compatibility
	if err := g.validateCompatibility(ctx); err != nil {
		return err
	}

	// 2. Collect base files
	hasDb := opts.Database != "" && opts.Database != "none"
	hasWorker := opts.Worker != "" && opts.Worker != "none"

	baseFiles, err := base.Render(base.Options{
		ProjectName:   opts.ProjectName,
		ModulePath:    opts.ModulePath,
		HasDatabase:   hasDb,
		DatabaseType:  opts.Database,
		HTTPFramework: opts.HTTP,
		WorkerType:    opts.Worker,
		HasWorker:     hasWorker,
	})
	if err != nil {
		return fmt.Errorf("failed to render base templates: %w", err)
	}

	var allFiles []GeneratedFile
	for _, bf := range baseFiles {
		allFiles = append(allFiles, GeneratedFile{
			Path:    bf.Path,
			Content: bf.Content,
		})
	}

	// 3. Collect plugin files and dependencies
	var dependencies []string
	var selectedPlugins []Plugin

	if hasDb {
		p, err := g.registry.Get(CategoryDatabase, opts.Database)
		if err != nil {
			return err
		}
		selectedPlugins = append(selectedPlugins, p)
	}

	if opts.HTTP != "" {
		p, err := g.registry.Get(CategoryHTTP, opts.HTTP)
		if err != nil {
			return err
		}
		selectedPlugins = append(selectedPlugins, p)
	}

	if hasWorker {
		p, err := g.registry.Get(CategoryWorker, opts.Worker)
		if err != nil {
			return err
		}
		selectedPlugins = append(selectedPlugins, p)
	}

	for _, p := range selectedPlugins {
		if err := p.Validate(ctx); err != nil {
			return fmt.Errorf("validation failed for plugin %s/%s: %w", p.Category(), p.Metadata().Name, err)
		}
		files, err := p.Generate(ctx)
		if err != nil {
			return fmt.Errorf("generation failed for plugin %s/%s: %w", p.Category(), p.Metadata().Name, err)
		}
		allFiles = append(allFiles, files...)
		dependencies = append(dependencies, p.Dependencies()...)
	}

	// 4. Check for file conflicts before writing any file
	if !opts.Overwrite {
		for _, f := range allFiles {
			dest := filepath.Join(opts.TargetDir, f.Path)
			if _, err := os.Stat(dest); err == nil {
				return fmt.Errorf("file %q already exists. Use --overwrite to overwrite existing files", f.Path)
			}
		}
	}

	// 5. Ensure target directory exists
	if err := os.MkdirAll(opts.TargetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory %s: %w", opts.TargetDir, err)
	}

	// 6. Write files
	for _, f := range allFiles {
		dest := filepath.Join(opts.TargetDir, f.Path)
		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", f.Path, err)
		}
		if err := os.WriteFile(dest, f.Content, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", f.Path, err)
		}
	}

	// 7. Save project manifest
	manifest := &Manifest{
		Project: ProjectInfo{
			Name:   opts.ProjectName,
			Module: opts.ModulePath,
		},
		Plugins: make(map[Category]PluginEntry),
	}
	for _, p := range selectedPlugins {
		manifest.Plugins[p.Category()] = PluginEntry{
			Name:    p.Metadata().Name,
			Version: p.Metadata().Version,
		}
	}
	if err := manifest.Save(opts.TargetDir); err != nil {
		return fmt.Errorf("failed to save manifest: %w", err)
	}

	// 8. Manage Go module and dependencies
	if !opts.SkipTidy {
		if err := g.setupGoDependencies(opts.TargetDir, opts.ModulePath, dependencies); err != nil {
			return fmt.Errorf("failed to setup Go dependencies: %w", err)
		}
	}

	return nil
}

func (g *Generator) validateCompatibility(ctx *GenerateContext) error {
	db := ctx.Plugins[CategoryDatabase]
	worker := ctx.Plugins[CategoryWorker]

	if worker == "river" {
		if db == "" || db == "none" {
			return fmt.Errorf("incompatible plugins: 'river' worker requires a PostgreSQL database, but no database was selected")
		}
		if db != "postgres" {
			return fmt.Errorf("incompatible plugins: 'river' worker requires 'postgres' database, but %q was selected", db)
		}
	}

	return nil
}

func (g *Generator) setupGoDependencies(dir, modulePath string, dependencies []string) error {
	goModPath := filepath.Join(dir, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		cmd := exec.Command("go", "mod", "init", modulePath)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("go mod init failed: %s: %w", string(out), err)
		}
	}

	for _, dep := range dependencies {
		if dep == "" {
			continue
		}
		cmd := exec.Command("go", "get", dep)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("go get %s failed: %s: %w", dep, string(out), err)
		}
	}

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go mod tidy failed: %s: %w", string(out), err)
	}

	return nil
}

func detectModulePath(dir string) (string, error) {
	goModPath := filepath.Join(dir, "go.mod")
	f, err := os.Open(goModPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}
	return "", fmt.Errorf("module line not found in %s", goModPath)
}
