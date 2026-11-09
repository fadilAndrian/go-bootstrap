package core_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fadilAndrian/go-bootstrap/internal/core"

	// Register plugins for generator test
	_ "github.com/fadilAndrian/go-bootstrap/plugins/database/mysql"
	_ "github.com/fadilAndrian/go-bootstrap/plugins/database/postgres"
	_ "github.com/fadilAndrian/go-bootstrap/plugins/database/sqlite"
	_ "github.com/fadilAndrian/go-bootstrap/plugins/http/chi"
	_ "github.com/fadilAndrian/go-bootstrap/plugins/http/gin"
	_ "github.com/fadilAndrian/go-bootstrap/plugins/http/nethttp"
	_ "github.com/fadilAndrian/go-bootstrap/plugins/worker/none"
	_ "github.com/fadilAndrian/go-bootstrap/plugins/worker/river"
)

func TestGenerator_PostgresGinRiver(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gen-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	gen := core.NewGenerator(core.DefaultRegistry)
	err = gen.Generate(&core.GenerateOptions{
		ProjectName: "my-app",
		ModulePath:  "github.com/user/my-app",
		TargetDir:   tempDir,
		Database:    "postgres",
		HTTP:        "gin",
		Worker:      "river",
		SkipTidy:    true,
	})
	if err != nil {
		t.Fatalf("unexpected generation error: %v", err)
	}

	// Verify key files are present
	expectedFiles := []string{
		"cmd/serve/main.go",
		"internal/bootstrap/app.go",
		"internal/bootstrap/lifecycle.go",
		"internal/config/config.go",
		"internal/database/postgres.go",
		"internal/http/server.go",
		"internal/http/gin.go",
		"internal/worker/river.go",
		"internal/worker/job.go",
		"internal/routes/routes.go",
		"modules/example/handler.go",
		"configs/config.yaml",
		".env.example",
		".gitignore",
		"README.md",
		".go-bootstrap.yaml",
	}

	for _, file := range expectedFiles {
		p := filepath.Join(tempDir, file)
		if _, err := os.Stat(p); os.IsNotExist(err) {
			t.Errorf("expected file %s does not exist", file)
		}
	}

	// Verify manifest contents
	m, err := core.LoadManifest(tempDir)
	if err != nil {
		t.Fatalf("failed to load manifest: %v", err)
	}
	if m.Project.Name != "my-app" {
		t.Errorf("expected my-app, got %s", m.Project.Name)
	}
	if m.Plugins[core.CategoryDatabase].Name != "postgres" {
		t.Errorf("expected postgres db plugin, got %s", m.Plugins[core.CategoryDatabase].Name)
	}
	if m.Plugins[core.CategoryHTTP].Name != "gin" {
		t.Errorf("expected gin http plugin, got %s", m.Plugins[core.CategoryHTTP].Name)
	}
	if m.Plugins[core.CategoryWorker].Name != "river" {
		t.Errorf("expected river worker plugin, got %s", m.Plugins[core.CategoryWorker].Name)
	}
}

func TestGenerator_IncompatibleWorker(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gen-incompat-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	gen := core.NewGenerator(core.DefaultRegistry)
	err = gen.Generate(&core.GenerateOptions{
		ProjectName: "my-app",
		ModulePath:  "github.com/user/my-app",
		TargetDir:   tempDir,
		Database:    "mysql",
		HTTP:        "gin",
		Worker:      "river", // River requires postgres!
		SkipTidy:    true,
	})

	if err == nil {
		t.Fatalf("expected error for incompatible combination mysql + river, got nil")
	}
}

func TestGenerator_OverwriteProtection(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gen-overwrite-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	gen := core.NewGenerator(core.DefaultRegistry)
	opts := &core.GenerateOptions{
		ProjectName: "my-app",
		ModulePath:  "github.com/user/my-app",
		TargetDir:   tempDir,
		Database:    "sqlite",
		HTTP:        "nethttp",
		Worker:      "none",
		SkipTidy:    true,
	}

	if err := gen.Generate(opts); err != nil {
		t.Fatalf("initial generation failed: %v", err)
	}

	// Generating again without overwrite should fail
	err = gen.Generate(opts)
	if err == nil {
		t.Fatalf("expected conflict error when overwrite=false, got nil")
	}

	// Generating with overwrite should succeed
	opts.Overwrite = true
	if err := gen.Generate(opts); err != nil {
		t.Fatalf("expected success when overwrite=true, got: %v", err)
	}
}
