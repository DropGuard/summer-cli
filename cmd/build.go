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

The build always starts from a clean target so stale classes cannot enter the
artifact. For incremental builds use Maven directly — the Summer Maven plugin
reconciles stale classes on every run either way.

All other arguments are passed to Maven, for example:
  summer build -DskipTests
`,
	DisableFlagParsing: true,
	Args:               cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		mavenArgs, showHelp, err := parseBuildArgs(args)
		if err != nil {
			return err
		}
		if showHelp {
			return cmd.Help()
		}

		if err := runner.RunBuild(mavenArgs...); err != nil {
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
// they can be forwarded, so build/check commands deliberately use
// DisableFlagParsing. The removed --incremental/--clean modes are rejected with
// a migration hint instead of being silently forwarded to Maven (which would
// fail with a cryptic "unknown lifecycle phase").
func parseBuildArgs(args []string) ([]string, bool, error) {
	var mavenArgs []string
	showHelp := false

	for i, arg := range args {
		switch arg {
		case "--incremental", "--clean":
			return nil, false, fmt.Errorf(
				"%s was removed: summer always builds from a clean target (stale-class reconciliation runs in the Maven plugin regardless); use plain 'mvn' for incremental builds",
				arg)
		case "--help", "-h":
			if len(args) == 1 {
				showHelp = true
			} else {
				mavenArgs = append(mavenArgs, arg)
			}
		case "--":
			mavenArgs = append(mavenArgs, args[i+1:]...)
			return mavenArgs, showHelp, nil
		default:
			mavenArgs = append(mavenArgs, arg)
		}
	}

	return mavenArgs, showHelp, nil
}
