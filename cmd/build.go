package cmd

import (
	"fmt"
	"os"

	"github.com/dropguard/summer-cli/pkg/runner"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Compile and package the application into a JAR",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runner.RunBuild(); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Build failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Build successful!")
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
