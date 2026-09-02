package runner

import (
	"reflect"
	"testing"
)

func TestBuildGoalAlwaysClean(t *testing.T) {
	// Always-clean is the contract: the CLI gives up the clean/incremental
	// choice so stale classes can never enter the artifact, and the Maven
	// plugin's reconciliation covers incremental builds done via plain mvn.
	got := runGoalArgs("package", "-DskipTests")
	want := []string{"clean", "package", "-DskipTests"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("runGoalArgs() = %#v, want %#v", got, want)
	}
}

func TestCheckGoalUsesProcessClasses(t *testing.T) {
	got := runGoalArgs("process-classes", "-DskipTests")
	want := []string{"clean", "process-classes", "-DskipTests"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("runGoalArgs() = %#v, want %#v", got, want)
	}
}

func TestDevGoalUsesFullyQualifiedCoordinates(t *testing.T) {
	// The prefix form (summer:dev) fails on default Maven installs — prefix
	// resolution only consults pluginGroups from settings.xml. The CLI must
	// spell out the coordinates so `summer dev` works with zero configuration.
	if mavenDevGoal != "io.github.dropguard:summer-maven-plugin:dev" {
		t.Fatalf("mavenDevGoal = %q, want fully-qualified coordinates", mavenDevGoal)
	}
}
