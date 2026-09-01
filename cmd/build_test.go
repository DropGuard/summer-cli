package cmd

import (
	"reflect"
	"testing"

	"github.com/dropguard/summer-cli/pkg/runner"
)

func TestParseBuildArgsKeepsMavenFlags(t *testing.T) {
	options, mavenArgs, showHelp, err := parseBuildArgs([]string{
		"--incremental",
		"-DskipTests",
		"-Pdev",
	})
	if err != nil {
		t.Fatalf("parseBuildArgs() error = %v", err)
	}
	if showHelp {
		t.Fatal("Maven arguments must not request CLI help")
	}
	if !options.Incremental {
		t.Fatal("--incremental must be recognized by the CLI")
	}
	want := []string{"-DskipTests", "-Pdev"}
	if !reflect.DeepEqual(mavenArgs, want) {
		t.Fatalf("Maven args = %#v, want %#v", mavenArgs, want)
	}
}

func TestParseBuildArgsRejectsConflictingModes(t *testing.T) {
	_, _, _, err := parseBuildArgs([]string{"--incremental", "--clean"})
	if err == nil {
		t.Fatal("conflicting build modes must be rejected")
	}
}

func TestParseBuildArgsDefaultsToClean(t *testing.T) {
	options, mavenArgs, showHelp, err := parseBuildArgs([]string{"-DskipTests"})
	if err != nil {
		t.Fatalf("parseBuildArgs() error = %v", err)
	}
	if showHelp || options != (runner.BuildOptions{}) {
		t.Fatalf("unexpected parse result: options=%#v help=%v", options, showHelp)
	}
	if !reflect.DeepEqual(mavenArgs, []string{"-DskipTests"}) {
		t.Fatalf("Maven args = %#v", mavenArgs)
	}
}

func TestParseBuildArgsSupportsSeparator(t *testing.T) {
	options, mavenArgs, _, err := parseBuildArgs([]string{"--incremental", "--", "-DskipTests"})
	if err != nil {
		t.Fatalf("parseBuildArgs() error = %v", err)
	}
	if !options.Incremental || !reflect.DeepEqual(mavenArgs, []string{"-DskipTests"}) {
		t.Fatalf("unexpected parse result: options=%#v args=%#v", options, mavenArgs)
	}
}
