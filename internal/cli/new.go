package cli

import (
	"fmt"
	"path/filepath"

	"github.com/fadilAndrian/go-bootstrap/internal/core"
	"github.com/spf13/cobra"
)

var (
	newDBFlag        string
	newHTTPFlag      string
	newWorkerFlag    string
	newModuleFlag    string
	newOverwriteFlag bool
)

var newCmd = &cobra.Command{
	Use:   "new <project-name>",
	Short: "Create a new Go project with selected components",
	Args:  cobra.ExactArgs(1),
	Example: `  go-bootstrap new my-project
  go-bootstrap new my-project --db=postgres --http=gin --worker=river
  go-bootstrap new my-project --db=sqlite --http=chi --worker=none`,
	RunE: func(cmd *cobra.Command, args []string) error {
		targetDir := args[0]
		projectName := filepath.Base(targetDir)

		modulePath := newModuleFlag
		if modulePath == "" {
			modulePath = projectName
		}

		fmt.Printf("🚀 Creating project %q (%s)...\n", projectName, modulePath)
		fmt.Printf("   • Database: %s\n", newDBFlag)
		fmt.Printf("   • HTTP:     %s\n", newHTTPFlag)
		fmt.Printf("   • Worker:   %s\n", newWorkerFlag)

		generator := core.NewGenerator(core.DefaultRegistry)
		err := generator.Generate(&core.GenerateOptions{
			ProjectName: projectName,
			ModulePath:  modulePath,
			TargetDir:   targetDir,
			Database:    newDBFlag,
			HTTP:        newHTTPFlag,
			Worker:      newWorkerFlag,
			Overwrite:   newOverwriteFlag,
		})
		if err != nil {
			return fmt.Errorf("project generation failed: %w", err)
		}

		fmt.Println("\n✅ Project successfully created!")
		fmt.Println("\nNext steps:")
		fmt.Printf("  cd %s\n", projectName)
		fmt.Println("  cp .env.example .env")
		fmt.Println("  go-bootstrap serve  # or 'go run ./cmd/serve'")

		return nil
	},
}

func init() {
	newCmd.Flags().StringVar(&newDBFlag, "db", "postgres", "Database adapter (postgres, mysql, sqlite, none)")
	newCmd.Flags().StringVar(&newHTTPFlag, "http", "gin", "HTTP framework (gin, chi, nethttp)")
	newCmd.Flags().StringVar(&newWorkerFlag, "worker", "river", "Worker adapter (river, none)")
	newCmd.Flags().StringVar(&newModuleFlag, "module", "", "Go module path (defaults to project name)")
	newCmd.Flags().BoolVar(&newOverwriteFlag, "overwrite", false, "Overwrite existing files")
}
