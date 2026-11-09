package core

import (
	"testing"
)

type dummyPlugin struct {
	meta corePluginMetaStub
}

type corePluginMetaStub struct {
	Name        string
	Category    Category
	Version     string
	Description string
}

func (d *dummyPlugin) Metadata() PluginMetadata {
	return PluginMetadata{
		Name:        d.meta.Name,
		Category:    d.meta.Category,
		Version:     d.meta.Version,
		Description: d.meta.Description,
	}
}

func (d *dummyPlugin) Category() Category {
	return d.meta.Category
}

func (d *dummyPlugin) Dependencies() []string {
	return nil
}

func (d *dummyPlugin) Validate(ctx *GenerateContext) error {
	return nil
}

func (d *dummyPlugin) Generate(ctx *GenerateContext) ([]GeneratedFile, error) {
	return nil, nil
}

func TestRegistry(t *testing.T) {
	reg := NewRegistry()

	p1 := &dummyPlugin{meta: corePluginMetaStub{Name: "pluginA", Category: CategoryDatabase, Version: "1.0.0"}}
	p2 := &dummyPlugin{meta: corePluginMetaStub{Name: "pluginB", Category: CategoryDatabase, Version: "1.1.0"}}

	if err := reg.Register(p1); err != nil {
		t.Fatalf("failed to register p1: %v", err)
	}
	if err := reg.Register(p2); err != nil {
		t.Fatalf("failed to register p2: %v", err)
	}

	// Registering duplicate should fail
	if err := reg.Register(p1); err == nil {
		t.Errorf("expected error when registering duplicate plugin")
	}

	// Get plugin
	found, err := reg.Get(CategoryDatabase, "pluginA")
	if err != nil {
		t.Fatalf("failed to get pluginA: %v", err)
	}
	if found.Metadata().Name != "pluginA" {
		t.Errorf("expected pluginA, got %s", found.Metadata().Name)
	}

	// Get non-existent
	if _, err := reg.Get(CategoryDatabase, "non-existent"); err == nil {
		t.Errorf("expected error getting non-existent plugin")
	}

	// List
	list := reg.List(CategoryDatabase)
	if len(list) != 2 {
		t.Errorf("expected 2 plugins, got %d", len(list))
	}
}
