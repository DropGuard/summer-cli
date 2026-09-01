package runner

import (
	"reflect"
	"testing"
)

func TestBuildArgsCleanByDefault(t *testing.T) {
	got := buildArgs(BuildOptions{}, "package", "-DskipTests")
	want := []string{"clean", "package", "-DskipTests"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("buildArgs() = %#v, want %#v", got, want)
	}
}

func TestBuildArgsIncrementalOmitsClean(t *testing.T) {
	got := buildArgs(BuildOptions{Incremental: true}, "package", "-Pdev")
	want := []string{"package", "-Pdev", "-Dsummer.build.incremental=true"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("buildArgs() = %#v, want %#v", got, want)
	}
}

func TestBuildArgsCheckUsesProcessClasses(t *testing.T) {
	got := buildArgs(BuildOptions{}, "process-classes", "-DskipTests")
	want := []string{"clean", "process-classes", "-DskipTests"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("buildArgs() = %#v, want %#v", got, want)
	}
}
