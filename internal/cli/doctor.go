package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/konnco/go-bootstrap/internal/core"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check project health and component compatibility",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		fmt.Printf("🔍 Running doctor check on %s...\n\n", cwd)

		// 1. Check Go toolchain
		goVerCmd := exec.Command("go", "version")
		if out, err := goVerCmd.CombinedOutput(); err != nil {
			fmt.Printf("❌ Go toolchain: NOT FOUND (%v)\n", err)
		} else {
			fmt.Printf("✅ Go toolchain: %s", string(out))
		}

		// 2. Check manifest
		if !core.ManifestExists(cwd) {
			fmt.Printf("❌ Manifest: .go-bootstrap.yaml NOT FOUND. Not a Go Bootstrap project?\n")
			return nil
		}
		manifest, err := core.LoadManifest(cwd)
		if err != nil {
			fmt.Printf("❌ Manifest: Failed to read (%v)\n", err)
			return nil
		}
		fmt.Printf("✅ Manifest: .go-bootstrap.yaml valid (Project: %s, Module: %s)\n",
			manifest.Project.Name, manifest.Project.Module)

		// 3. Check go.mod
		if _, err := os.Stat(filepath.Join(cwd, "go.mod")); os.IsNotExist(err) {
			fmt.Println("❌ go.mod: NOT FOUND")
		} else {
			fmt.Println("✅ go.mod: Present")
		}

		// 4. Check core entry points
		requiredFiles := []string{
			"cmd/serve/main.go",
			"internal/bootstrap/app.go",
			"internal/bootstrap/lifecycle.go",
			"internal/config/config.go",
		}
		allFound := true
		for _, f := range requiredFiles {
			if _, err := os.Stat(filepath.Join(cwd, f)); os.IsNotExist(err) {
				fmt.Printf("❌ File: missing %s\n", f)
				allFound = false
			}
		}
		if allFound {
			fmt.Println("✅ Core structure: All key files present")
		}

		// 5. Check plugin compatibility from manifest
		db := manifest.Plugins[core.CategoryDatabase].Name
		worker := manifest.Plugins[core.CategoryWorker].Name
		httpPlugin := manifest.Plugins[core.CategoryHTTP].Name

		fmt.Printf("\nConfigured Components:\n")
		fmt.Printf("  • Database: %s\n", db)
		fmt.Printf("  • HTTP:     %s\n", httpPlugin)
		fmt.Printf("  • Worker:   %s\n", worker)

		if worker == "river" && db != "postgres" {
			fmt.Printf("❌ Compatibility: River worker requires PostgreSQL, but %q is configured!\n", db)
		} else {
			fmt.Println("✅ Compatibility: Selected components are mutually compatible")
		}

		fmt.Println("\nAll diagnostic checks complete.")
		return nil
	},
}
