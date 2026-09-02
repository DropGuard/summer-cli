package runner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func checkMaven() error {
	_, err := exec.LookPath("mvn")
	if err != nil {
		return fmt.Errorf("maven ('mvn') is not installed or not in your PATH. Please install Maven to use this command")
	}
	return nil
}

// mavenDevGoal is the fully-qualified dev goal. Maven's short prefix form
// (`summer:dev`) resolves only through pluginGroups configured in
// settings.xml — a default installation fails with "No plugin found for
// prefix 'summer'". The CLI spells out the coordinates so `summer dev` works
// on a pristine machine; the length never reaches the user's command line.
const mavenDevGoal = "io.github.dropguard:summer-maven-plugin:dev"

func RunDev(extraArgs ...string) error {
	if err := checkMaven(); err != nil {
		return err
	}

	// dev is long-lived: the CLI's own SIGINT/SIGTERM must reach the Maven
	// process group so a non-interactive kill (docker stop, CI timeout, kill
	// from another terminal) does not orphan the JVM still holding port 8080.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	args := append([]string{mavenDevGoal}, extraArgs...)
	cmd := exec.CommandContext(ctx, "mvn", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	setProcessGroup(cmd)
	// Cancel sends the same signal to the whole group (not just the mvn child),
	// then restores default signal behavior so a second Ctrl-C hard-kills.
	cmd.Cancel = func() error {
		stop()
		return signalProcessGroup(cmd)
	}
	// Never hang forever on a stuck child that ignores SIGTERM.
	cmd.WaitDelay = 5 * time.Second

	fmt.Println("🚀 Starting Summer Framework in Dev Mode...")
	return cmd.Run()
}

// RunBuild always starts from a clean target so stale classes cannot enter the
// artifact. The Maven plugin's source-set reconciliation runs unconditionally
// on every generate-aot invocation, but a clean compile is the only guarantee
// that holds for users who switch engines, branches, or plugin versions.
func RunBuild(extraArgs ...string) error {
	return runGoal("package", "Building Summer application...", extraArgs)
}

// RunCheck validates compilation, Jandex discovery, dependency resolution, and
// AOT source compilation without packaging an application artifact.
func RunCheck(extraArgs ...string) error {
	return runGoal("process-classes", "Checking Summer application...", extraArgs)
}

func runGoal(goal, message string, extraArgs []string) error {
	if err := checkMaven(); err != nil {
		return err
	}

	cmd := exec.Command("mvn", runGoalArgs(goal, extraArgs...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fmt.Println(message)
	return cmd.Run()
}

// runGoalArgs assembles the Maven invocation: always a clean build, then the
// lifecycle goal, then user-passed Maven arguments verbatim.
func runGoalArgs(goal string, extraArgs ...string) []string {
	args := append([]string{"clean", goal}, extraArgs...)
	return args
}
