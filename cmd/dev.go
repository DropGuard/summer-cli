package cmd

import (
	"fmt"

	"github.com/dropguard/summer-cli/pkg/runner"
	"github.com/spf13/cobra"
)

var devCmd = &cobra.Command{
	Use:                "dev [maven-args...]",
	Short:              "Start the application in development mode with hot-reload",
	DisableFlagParsing: true,
	Args:               cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 && (args[0] == "--help" || args[0] == "-h") {
			return cmd.Help()
		}
		if len(args) > 0 && args[0] == "--" {
			args = args[1:]
		}
		if err := runner.RunDev(args...); err != nil {
			return fmt.Errorf("dev mode failed: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(devCmd)
}
