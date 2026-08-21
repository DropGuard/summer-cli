package generator

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/dropguard/summer-cli/templates"
)

// ProjectOptions holds the parameters for generation
type ProjectOptions struct {
	Name             string
	GroupId          string
	ArtifactId       string
	Package          string
	FrameworkVersion string
}

// Generate extracts the embedded archetype fs and applies template replacements
func Generate(opts ProjectOptions) error {
	// 1. Create the project root directory
	if err := os.MkdirAll(opts.Name, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	packagePath := strings.ReplaceAll(opts.Package, ".", "/")

	// 2. Walk the embedded filesystem
	err := fs.WalkDir(templates.ScaffoldFS, "scaffold", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root 'scaffold' directory itself
		if path == "scaffold" {
			return nil
		}

		// Determine target path relative to 'scaffold/'
		relativePath := strings.TrimPrefix(path, "scaffold/")
		targetPath := relativePath

		// Handle Java source directory relocation based on package
		if strings.HasPrefix(targetPath, "src/main/java/") {
			fileRel := strings.TrimPrefix(targetPath, "src/main/java/")
			targetPath = filepath.Join("src/main/java", packagePath, fileRel)
		} else if strings.HasPrefix(targetPath, "src/test/java/") {
			fileRel := strings.TrimPrefix(targetPath, "src/test/java/")
			targetPath = filepath.Join("src/test/java", packagePath, fileRel)
		}

		// Prepend the project name (the root output dir)
		finalPath := filepath.Join(opts.Name, targetPath)

		if d.IsDir() {
			if err := os.MkdirAll(finalPath, 0755); err != nil {
				return err
			}
			return nil
		}

		// Ensure parent directory exists (needed because we relocated targetPath)
		if err := os.MkdirAll(filepath.Dir(finalPath), 0755); err != nil {
			return err
		}

		// Process file content
		content, err := fs.ReadFile(templates.ScaffoldFS, path)
		if err != nil {
			return err
		}

		strContent := string(content)
		strContent = strings.ReplaceAll(strContent, "${groupId}", opts.GroupId)
		strContent = strings.ReplaceAll(strContent, "${artifactId}", opts.ArtifactId)
		strContent = strings.ReplaceAll(strContent, "${package}", opts.Package)
		strContent = strings.ReplaceAll(strContent, "${version}", "1.0-SNAPSHOT")
		strContent = strings.ReplaceAll(strContent, "${frameworkVersion}", opts.FrameworkVersion)

		return os.WriteFile(finalPath, []byte(strContent), 0644)
	})

	if err != nil {
		return fmt.Errorf("failed to generate project: %w", err)
	}

	return nil
}
