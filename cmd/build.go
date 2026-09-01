package cmd

import (
	"fmt"

	"github.com/dropguard/summer-cli/pkg/runner"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build [maven-args...]",
	Short: "Compile and package the application into a JAR",
	Long: `Build a Summer application.

The default build starts with a clean target to prevent stale generated code
or bytecode from entering the artifact. Use --incremental only when you want
Maven to reuse target/; use --clean explicitly to document the safe mode.

All other arguments are passed to Maven, for example:
  summer build -DskipTests
  summer build --incremental -DskipTests
`,
	DisableFlagParsing: true,
	Args:               cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		options, mavenArgs, showHelp, err := parseBuildArgs(args)
		if err != nil {
			return err
		}
		if showHelp {
			return cmd.Help()
		}

		if err := runner.RunBuildWithOptions(options, mavenArgs...); err != nil {
			return fmt.Errorf("build failed: %w", err)
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Build successful!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}

// parseBuildArgs handles Summer-owned flags while leaving every other argument
// untouched for Maven. Cobra's normal parser rejects Maven's -D... flags before
// they can be forwarded, so build/dev commands deliberately use DisableFlagParsing.
func parseBuildArgs(args []string) (runner.BuildOptions, []string, bool, error) {
	var options runner.BuildOptions
	var mavenArgs []string
	cleanRequested := false
	showHelp := false

	for i, arg := range args {
		switch arg {
		case "--incremental":
			if cleanRequested {
				return runner.BuildOptions{}, nil, false,
					fmt.Errorf("--incremental cannot be combined with --clean")
			}
			options.Incremental = true
		case "--clean":
			if options.Incremental {
				return runner.BuildOptions{}, nil, false,
					fmt.Errorf("--clean cannot be combined with --incremental")
			}
			cleanRequested = true
		case "--help", "-h":
			if len(args) == 1 {
				showHelp = true
			} else {
				mavenArgs = append(mavenArgs, arg)
			}
		case "--":
			mavenArgs = append(mavenArgs, args[i+1:]...)
			return options, mavenArgs, showHelp, nil
		default:
			mavenArgs = append(mavenArgs, arg)
		}
	}

	return options, mavenArgs, showHelp, nil
}
