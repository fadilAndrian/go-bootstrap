package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestManifestSaveAndLoad(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "manifest-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	m := &Manifest{
		Project: ProjectInfo{
			Name:   "test-app",
			Module: "github.com/test/app",
		},
		Plugins: map[Category]PluginEntry{
			CategoryDatabase: {Name: "postgres", Version: "1.0.0"},
			CategoryHTTP:     {Name: "gin", Version: "1.0.0"},
			CategoryWorker:   {Name: "river", Version: "1.0.0"},
		},
	}

	if err := m.Save(tempDir); err != nil {
		t.Fatalf("failed to save manifest: %v", err)
	}

	if !ManifestExists(tempDir) {
		t.Fatalf("expected manifest to exist in %s", tempDir)
	}

	loaded, err := LoadManifest(tempDir)
	if err != nil {
		t.Fatalf("failed to load manifest: %v", err)
	}

	if loaded.Project.Name != m.Project.Name {
		t.Errorf("expected project name %q, got %q", m.Project.Name, loaded.Project.Name)
	}

	if loaded.Project.Module != m.Project.Module {
		t.Errorf("expected module %q, got %q", m.Project.Module, loaded.Project.Module)
	}

	if len(loaded.Plugins) != 3 {
		t.Errorf("expected 3 plugins, got %d", len(loaded.Plugins))
	}

	if loaded.Plugins[CategoryDatabase].Name != "postgres" {
		t.Errorf("expected postgres, got %s", loaded.Plugins[CategoryDatabase].Name)
	}
}

func TestManifestNotExists(t *testing.T) {
	nonExistentDir := filepath.Join(os.TempDir(), "random-non-existent-dir-12345")
	if ManifestExists(nonExistentDir) {
		t.Errorf("expected ManifestExists to return false for non-existent dir")
	}
}
