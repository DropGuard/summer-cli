package cmd

import (
	"fmt"
	"os"

	"github.com/dropguard/summer-cli/pkg/generator"
	"github.com/spf13/cobra"
)

var (
	groupId          string
	artifactId       string
	pkgName          string
	frameworkVersion string
)

var createCmd = &cobra.Command{
	Use:   "create [projectName]",
	Short: "Scaffold a new Summer Framework project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]

		if artifactId == "" {
			artifactId = projectName
		}
		if err := generator.ValidateGroupId(groupId); err != nil {
			fmt.Fprintf(os.Stderr, "\u274c %v\n", err)
			os.Exit(1)
		}
		if pkgName == "" {
			// Derived segments are sanitized ("my-app" -> "myapp") so the scaffold compiles;
			// explicit -p values are validated as-is below.
			pkgName = fmt.Sprintf("%s.%s", groupId, generator.SanitizePackageSegment(artifactId))
		}
		if err := generator.ValidatePackage(pkgName); err != nil {
			fmt.Fprintf(os.Stderr, "\u274c %v\n", err)
			os.Exit(1)
		}

		opts := generator.ProjectOptions{
			Name:             projectName,
			GroupId:          groupId,
			ArtifactId:       artifactId,
			Package:          pkgName,
			FrameworkVersion: frameworkVersion,
		}

		fmt.Printf("🌱 Creating Summer app '%s'...\n", projectName)
		if err := generator.Generate(opts); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to create project: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Project created successfully!\n\n")
		fmt.Printf("Next steps:\n")
		fmt.Printf("  cd %s\n", projectName)
		fmt.Printf("  summer dev\n")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&groupId, "group-id", "g", "com.example", "Maven groupId")
	createCmd.Flags().StringVarP(&artifactId, "artifact-id", "a", "", "Maven artifactId (defaults to projectName)")
	createCmd.Flags().StringVarP(&pkgName, "package", "p", "", "Base package (defaults to groupId.artifactId)")
	createCmd.Flags().StringVarP(&frameworkVersion, "framework-version", "f", FrameworkVersion, "Summer Framework version to use (defaults to the CLI release)")
}
