package cmd

import (
	"fmt"

	"github.com/dropguard/summer-cli/pkg/runner"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check [maven-args...]",
	Short: "Validate Summer DI and AOT wiring without packaging",
	Long: `Validate a Summer application without creating an artifact.

The default check starts with a clean target so removed generated code and
bytecode cannot affect the result. Use --incremental only when target/ is
known to be current.

All other arguments are passed to Maven, for example:
  summer check -Dsummer.engine=runtime
  summer check --incremental -DskipTests
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

		if err := runner.RunCheck(options, mavenArgs...); err != nil {
			return fmt.Errorf("check failed: %w", err)
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Summer wiring is valid.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
