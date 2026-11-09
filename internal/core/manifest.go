package core

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const ManifestFilename = ".go-bootstrap.yaml"

type Manifest struct {
	Project ProjectInfo              `yaml:"project"`
	Plugins map[Category]PluginEntry `yaml:"plugins"`
}

type ProjectInfo struct {
	Name   string `yaml:"name"`
	Module string `yaml:"module"`
}

type PluginEntry struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

// LoadManifest reads .go-bootstrap.yaml from the specified directory.
func LoadManifest(dir string) (*Manifest, error) {
	manifestPath := filepath.Join(dir, ManifestFilename)
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", ManifestFilename, err)
	}

	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", ManifestFilename, err)
	}

	if m.Plugins == nil {
		m.Plugins = make(map[Category]PluginEntry)
	}

	return &m, nil
}

// SaveManifest writes .go-bootstrap.yaml to the specified directory.
func (m *Manifest) Save(dir string) error {
	manifestPath := filepath.Join(dir, ManifestFilename)
	data, err := yaml.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed to serialize manifest: %w", err)
	}

	if err := os.WriteFile(manifestPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", ManifestFilename, err)
	}

	return nil
}

// ManifestExists checks if .go-bootstrap.yaml exists in the given directory.
func ManifestExists(dir string) bool {
	manifestPath := filepath.Join(dir, ManifestFilename)
	info, err := os.Stat(manifestPath)
	return err == nil && !info.IsDir()
}
