package validator

import (
	"testing"

	"github.com/fadilAndrian/go-bootstrap/internal/core"
)

func TestValidateCompatibility(t *testing.T) {
	tests := []struct {
		name      string
		db        string
		http      string
		worker    string
		expectErr bool
	}{
		{
			name:      "Valid Postgres + Gin + River",
			db:        "postgres",
			http:      "gin",
			worker:    "river",
			expectErr: false,
		},
		{
			name:      "Valid SQLite + Chi + None",
			db:        "sqlite",
			http:      "chi",
			worker:    "none",
			expectErr: false,
		},
		{
			name:      "Invalid MySQL + River",
			db:        "mysql",
			http:      "gin",
			worker:    "river",
			expectErr: true,
		},
		{
			name:      "Invalid SQLite + River",
			db:        "sqlite",
			http:      "gin",
			worker:    "river",
			expectErr: true,
		},
		{
			name:      "Invalid No DB + River",
			db:        "none",
			http:      "gin",
			worker:    "river",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &core.GenerateContext{
				ProjectName: "test-app",
				ModulePath:  "github.com/test/app",
				Plugins: map[core.Category]string{
					core.CategoryDatabase: tt.db,
					core.CategoryHTTP:     tt.http,
					core.CategoryWorker:   tt.worker,
				},
			}

			err := ValidateCompatibility(ctx)
			if (err != nil) != tt.expectErr {
				t.Errorf("ValidateCompatibility() error = %v, expectErr = %v", err, tt.expectErr)
			}
		})
	}
}
