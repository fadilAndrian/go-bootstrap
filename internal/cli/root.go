package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-bootstrap",
	Short: "A modular, extensible Go project generator",
	Long: `Go Bootstrap is a CLI tool that helps you create runnable, production-grade Go
applications by selecting the components you need (database, HTTP framework, background worker).

Generated projects contain idiomatic, editable Go source code with no runtime CLI lock-in.`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(serveCmd)
}
