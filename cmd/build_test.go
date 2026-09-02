package cmd

import (
	"reflect"
	"testing"
)

func TestParseBuildArgsKeepsMavenFlags(t *testing.T) {
	mavenArgs, showHelp, err := parseBuildArgs([]string{"-DskipTests", "-Pdev"})
	if err != nil {
		t.Fatalf("parseBuildArgs() error = %v", err)
	}
	if showHelp {
		t.Fatal("Maven arguments must not request CLI help")
	}
	want := []string{"-DskipTests", "-Pdev"}
	if !reflect.DeepEqual(mavenArgs, want) {
		t.Fatalf("Maven args = %#v, want %#v", mavenArgs, want)
	}
}

func TestParseBuildArgsRejectsRemovedModes(t *testing.T) {
	for _, removed := range []string{"--incremental", "--clean"} {
		if _, _, err := parseBuildArgs([]string{removed}); err == nil {
			t.Fatalf("parseBuildArgs(%q) = nil error, want migration hint", removed)
		}
	}
}

func TestParseBuildArgsSupportsSeparator(t *testing.T) {
	mavenArgs, _, err := parseBuildArgs([]string{"--", "-DskipTests", "-Pdev"})
	if err != nil {
		t.Fatalf("parseBuildArgs() error = %v", err)
	}
	if !reflect.DeepEqual(mavenArgs, []string{"-DskipTests", "-Pdev"}) {
		t.Fatalf("Maven args = %#v", mavenArgs)
	}
}
