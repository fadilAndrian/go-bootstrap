package cli

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the generated Go application (convenience command for 'go run ./cmd/serve')",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		entrypoint := filepath.Join(cwd, "cmd", "serve", "main.go")
		if _, err := os.Stat(entrypoint); os.IsNotExist(err) {
			return fmt.Errorf("entrypoint %s not found. Are you in the root of a generated project?", entrypoint)
		}

		fmt.Println("🚀 Starting application via 'go run ./cmd/serve'...")
		runCmd := exec.Command("go", "run", "./cmd/serve")
		runCmd.Stdin = os.Stdin
		runCmd.Stdout = os.Stdout
		runCmd.Stderr = os.Stderr

		// Handle OS signals and forward to child process
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		go func() {
			for sig := range sigChan {
				if runCmd.Process != nil {
					_ = runCmd.Process.Signal(sig)
				}
			}
		}()

		if err := runCmd.Run(); err != nil {
			return fmt.Errorf("application exited: %w", err)
		}

		return nil
	},
}
