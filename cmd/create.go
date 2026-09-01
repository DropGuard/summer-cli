package cmd

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dropguard/summer-cli/pkg/generator"
	"github.com/spf13/cobra"
)

var (
	groupId          string
	artifactId       string
	pkgName          string
	frameworkVersion string
)

type mavenMetadata struct {
	Versioning struct {
		Release string `xml:"release"`
		Latest  string `xml:"latest"`
	} `xml:"versioning"`
}

func fetchLatestFrameworkVersion(fallback string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://repo1.maven.org/maven2/io/github/dropguard/summer-parent/maven-metadata.xml", nil)
	if err != nil {
		return fallback
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return fallback
	}
	defer resp.Body.Close()

	var meta mavenMetadata
	if err := xml.NewDecoder(resp.Body).Decode(&meta); err != nil {
		return fallback
	}

	if meta.Versioning.Release != "" {
		return meta.Versioning.Release
	}
	if meta.Versioning.Latest != "" {
		return meta.Versioning.Latest
	}
	return fallback
}

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
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			os.Exit(1)
		}
		if pkgName == "" {
			// Derived segments are sanitized ("my-app" -> "myapp") so the scaffold compiles;
			// explicit -p values are validated as-is below.
			pkgName = fmt.Sprintf("%s.%s", groupId, generator.SanitizePackageSegment(artifactId))
		}
		if err := generator.ValidatePackage(pkgName); err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			os.Exit(1)
		}

		resolvedVersion := frameworkVersion
		if !cmd.Flags().Changed("framework-version") {
			resolvedVersion = fetchLatestFrameworkVersion(FrameworkVersion)
		}

		opts := generator.ProjectOptions{
			Name:             projectName,
			GroupId:          groupId,
			ArtifactId:       artifactId,
			Package:          pkgName,
			FrameworkVersion: resolvedVersion,
		}

		fmt.Printf("🌱 Creating Summer app '%s' (framework v%s)...\n", projectName, resolvedVersion)
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
	createCmd.Flags().StringVarP(&frameworkVersion, "framework-version", "f", FrameworkVersion, "Summer Framework version to use (defaults to latest Maven Central release)")
}
