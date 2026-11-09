package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fadilAndrian/go-bootstrap/internal/core"
	"github.com/spf13/cobra"
)

var (
	initDBFlag        string
	initHTTPFlag      string
	initWorkerFlag    string
	initOverwriteFlag bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the current directory as a Go Bootstrap project",
	Example: `  go-bootstrap init
  go-bootstrap init --db=postgres --http=gin --worker=river`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current working directory: %w", err)
		}
		projectName := filepath.Base(cwd)

		fmt.Printf("🚀 Initializing project in current directory %q...\n", cwd)
		fmt.Printf("   • Database: %s\n", initDBFlag)
		fmt.Printf("   • HTTP:     %s\n", initHTTPFlag)
		fmt.Printf("   • Worker:   %s\n", initWorkerFlag)

		generator := core.NewGenerator(core.DefaultRegistry)
		err = generator.Generate(&core.GenerateOptions{
			ProjectName: projectName,
			TargetDir:   cwd,
			Database:    initDBFlag,
			HTTP:        initHTTPFlag,
			Worker:      initWorkerFlag,
			Overwrite:   initOverwriteFlag,
		})
		if err != nil {
			return fmt.Errorf("project initialization failed: %w", err)
		}

		fmt.Println("\n✅ Project successfully initialized!")
		fmt.Println("\nNext steps:")
		fmt.Println("  cp .env.example .env")
		fmt.Println("  go-bootstrap serve  # or 'go run ./cmd/serve'")

		return nil
	},
}

func init() {
	initCmd.Flags().StringVar(&initDBFlag, "db", "postgres", "Database adapter (postgres, mysql, sqlite, none)")
	initCmd.Flags().StringVar(&initHTTPFlag, "http", "gin", "HTTP framework (gin, chi, nethttp)")
	initCmd.Flags().StringVar(&initWorkerFlag, "worker", "river", "Worker adapter (river, none)")
	initCmd.Flags().BoolVar(&initOverwriteFlag, "overwrite", false, "Overwrite existing files")
}
