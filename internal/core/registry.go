package core

import (
	"fmt"
	"sort"
	"sync"
)

type Registry struct {
	mu      sync.RWMutex
	plugins map[Category]map[string]Plugin
}

var DefaultRegistry = NewRegistry()

func NewRegistry() *Registry {
	return &Registry{
		plugins: make(map[Category]map[string]Plugin),
	}
}

func (r *Registry) Register(p Plugin) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cat := p.Category()
	meta := p.Metadata()

	if _, ok := r.plugins[cat]; !ok {
		r.plugins[cat] = make(map[string]Plugin)
	}

	if _, exists := r.plugins[cat][meta.Name]; exists {
		return fmt.Errorf("plugin %s/%s is already registered", cat, meta.Name)
	}

	r.plugins[cat][meta.Name] = p
	return nil
}

func (r *Registry) Get(category Category, name string) (Plugin, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	catMap, ok := r.plugins[category]
	if !ok {
		return nil, fmt.Errorf("no plugins registered for category %q", category)
	}

	p, ok := catMap[name]
	if !ok {
		return nil, fmt.Errorf("plugin %q not found in category %q", name, category)
	}

	return p, nil
}

func (r *Registry) List(category Category) []Plugin {
	r.mu.RLock()
	defer r.mu.RUnlock()

	catMap, ok := r.plugins[category]
	if !ok {
		return nil
	}

	list := make([]Plugin, 0, len(catMap))
	for _, p := range catMap {
		list = append(list, p)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Metadata().Name < list[j].Metadata().Name
	})

	return list
}

func (r *Registry) ListAll() map[Category][]Plugin {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[Category][]Plugin)
	for cat := range r.plugins {
		result[cat] = r.List(cat)
	}

	return result
}

func Register(p Plugin) {
	if err := DefaultRegistry.Register(p); err != nil {
		panic(err)
	}
}

func Get(category Category, name string) (Plugin, error) {
	return DefaultRegistry.Get(category, name)
}

func List(category Category) []Plugin {
	return DefaultRegistry.List(category)
}

func ListAll() map[Category][]Plugin {
	return DefaultRegistry.ListAll()
}
