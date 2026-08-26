package main

import (
	"github.com/dropguard/summer-cli/cmd"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestCLI_E2E builds the summer-cli binary and executes it against a temporary directory,
// verifying the creation of the project and that `summer build` runs cleanly.
func TestCLI_E2E(t *testing.T) {
	// 1. Build the binary
	binaryPath := filepath.Join(t.TempDir(), "summer_test_bin")
	buildCmd := exec.Command("go", "build", "-buildvcs=false", "-o", binaryPath)
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to compile summer binary for E2E test: %v", err)
	}

	// 2. Setup a temporary directory for the generated project
	testWorkspace := t.TempDir()

	// 3. Run `summer create`
	frameworkVersion := os.Getenv("SUMMER_E2E_FRAMEWORK_VERSION")
	if frameworkVersion == "" {
		frameworkVersion = cmd.FrameworkVersion // the version this CLI build pins by default
	}
	createCmd := exec.Command(binaryPath, "create", "e2e-app", "--group-id", "com.e2e",
		"--package", "com.e2e.app", "--framework-version", frameworkVersion)
	createCmd.Dir = testWorkspace
	createCmd.Stdout = os.Stdout
	createCmd.Stderr = os.Stderr
	if err := createCmd.Run(); err != nil {
		t.Fatalf("'summer create' command failed: %v", err)
	}

	appDir := filepath.Join(testWorkspace, "e2e-app")

	// 4. Verify basic assertions
	if _, err := os.Stat(filepath.Join(appDir, "pom.xml")); os.IsNotExist(err) {
		t.Fatalf("pom.xml was not created in %s", appDir)
	}
	if _, err := os.Stat(filepath.Join(appDir, "src/main/java/com/e2e/app/App.java")); os.IsNotExist(err) {
		t.Fatalf("App.java was not created in the correct package directory")
	}

	// 5. Run `summer build`
	// This ensures the generated code is actually valid and compiles with Maven.
	runBuildCmd := exec.Command(binaryPath, "build")
	runBuildCmd.Dir = appDir
	runBuildCmd.Stdout = os.Stdout
	runBuildCmd.Stderr = os.Stderr
	if err := runBuildCmd.Run(); err != nil {
		t.Fatalf("'summer build' command failed, meaning the generated project has a compile error: %v", err)
	}

	// 6. Verify that a JAR was created in the target directory
	targetDir := filepath.Join(appDir, "target")
	entries, err := os.ReadDir(targetDir)
	if err != nil {
		t.Fatalf("Failed to read target directory: %v", err)
	}

	foundJar := false
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) == ".jar" {
			foundJar = true
			break
		}
	}
	if !foundJar {
		t.Fatalf("Build succeeded but no JAR file was found in target/")
	}
}
