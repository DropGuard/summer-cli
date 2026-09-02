package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// Version is the CLI's own version, injected at release build time via
// -ldflags "-X github.com/dropguard/summer-cli/cmd.Version=vX.Y.Z". "dev" in source builds.
var Version = "dev"

var rootCmd = &cobra.Command{
	Use:   "summer",
	Short: "Summer CLI - The official CLI for the Summer Framework",
	Long: `Summer CLI is a fast, dependency-free tool to scaffold and manage
Summer Framework projects. It embeds the archetype and wraps Maven
commands to provide a modern, frictionless developer experience.`,
	Version: Version,
	// Cobra prints "Error: ..." itself; Execute() below prints the message once
	// more with exit-code handling. Silence cobra's copy to avoid double-printing.
	SilenceErrors: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		// Preserve Maven's own exit code so CI scripts can distinguish
		// compilation failures from test failures; -1 (killed by signal)
		// degrades to 1.
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() > 0 {
			os.Exit(exitErr.ExitCode())
		}
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
}
