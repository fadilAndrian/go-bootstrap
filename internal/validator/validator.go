package validator

import (
	"fmt"

	"github.com/fadilAndrian/go-bootstrap/internal/core"
)

// ValidateCompatibility checks if the chosen plugins are mutually compatible.
func ValidateCompatibility(ctx *core.GenerateContext) error {
	db := ctx.Plugins[core.CategoryDatabase]
	worker := ctx.Plugins[core.CategoryWorker]

	// River worker requires PostgreSQL
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
