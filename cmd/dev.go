package cmd

import (
	"fmt"
	"os"

	"github.com/dropguard/summer-cli/pkg/runner"
	"github.com/spf13/cobra"
)

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Start the application in development mode with hot-reload",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runner.RunDev(); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Dev mode failed: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(devCmd)
}
