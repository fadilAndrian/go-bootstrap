package cli

import (
	"fmt"

	"github.com/konnco/go-bootstrap/internal/core"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available plugins and components",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Available Go Bootstrap Plugins:")
		fmt.Println("================================")

		categories := []core.Category{
			core.CategoryDatabase,
			core.CategoryHTTP,
			core.CategoryWorker,
		}

		for _, cat := range categories {
			fmt.Printf("\n[%s]\n", cat)
			plugins := core.List(cat)
			if len(plugins) == 0 {
				fmt.Println("  (none registered)")
				continue
			}
			for _, p := range plugins {
				m := p.Metadata()
				fmt.Printf("  • %-10s v%-6s %s\n", m.Name, m.Version, m.Description)
			}
		}
		fmt.Println()
	},
}
