package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	// Create a temporary directory for output
	tempDir := t.TempDir()

	opts := ProjectOptions{
		Name:             filepath.Join(tempDir, "my-app"),
		GroupId:          "com.test",
		ArtifactId:       "my-app",
		Package:          "com.test.myapp",
		FrameworkVersion: "999-SNAPSHOT",
	}

	err := Generate(opts)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// 1. Verify root directory and pom.xml exist
	pomPath := filepath.Join(opts.Name, "pom.xml")
	if _, err := os.Stat(pomPath); os.IsNotExist(err) {
		t.Errorf("pom.xml was not created")
	}

	// 2. Verify placeholder replacement in pom.xml
	pomBytes, _ := os.ReadFile(pomPath)
	pomContent := string(pomBytes)
	if strings.Contains(pomContent, "${groupId}") {
		t.Errorf("pom.xml still contains ${groupId}")
	}
	if !strings.Contains(pomContent, "<groupId>com.test</groupId>") {
		t.Errorf("pom.xml missing expected groupId replacement")
	}

	// 3. Verify Java files were relocated properly
	appJavaPath := filepath.Join(opts.Name, "src/main/java/com/test/myapp/App.java")
	if _, err := os.Stat(appJavaPath); os.IsNotExist(err) {
		t.Errorf("App.java was not relocated correctly: %v", err)
	}

	// 4. Verify package statement in App.java
	appBytes, _ := os.ReadFile(appJavaPath)
	appContent := string(appBytes)
	if !strings.Contains(appContent, "package com.test.myapp;") {
		t.Errorf("App.java missing expected package statement")
	}
}
