package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "summer",
	Short: "Summer CLI - The official CLI for the Summer Framework",
	Long: `Summer CLI is a fast, dependency-free tool to scaffold and manage
Summer Framework projects. It embeds the archetype and wraps Maven
commands to provide a modern, frictionless developer experience.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
}
