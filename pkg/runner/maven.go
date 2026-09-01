package runner

import (
	"fmt"
	"os"
	"os/exec"
)

// BuildOptions controls the safety/performance trade-off of a Summer build.
// A clean target is the default because Maven does not remove every stale class
// when a source or generated type disappears. Incremental is an explicit opt-in.
type BuildOptions struct {
	Incremental bool
}

func checkMaven() error {
	_, err := exec.LookPath("mvn")
	if err != nil {
		return fmt.Errorf("maven ('mvn') is not installed or not in your PATH. Please install Maven to use this command")
	}
	return nil
}

func RunDev(extraArgs ...string) error {
	if err := checkMaven(); err != nil {
		return err
	}

	args := append([]string{"summer:dev"}, extraArgs...)
	cmd := exec.Command("mvn", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fmt.Println("🚀 Starting Summer Framework in Dev Mode...")
	return cmd.Run()
}

// RunBuild keeps the original safe-build entry point for callers of the runner
// package. Use RunBuildWithOptions when the incremental mode is intentional.
func RunBuild(extraArgs ...string) error {
	return RunBuildWithOptions(BuildOptions{}, extraArgs...)
}

func RunBuildWithOptions(options BuildOptions, extraArgs ...string) error {
	if err := checkMaven(); err != nil {
		return err
	}

	args := buildArgs(options, "package", extraArgs...)
	cmd := exec.Command("mvn", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if options.Incremental {
		fmt.Println("Building Summer application incrementally...")
	} else {
		fmt.Println("Building Summer application from a clean target...")
	}
	return cmd.Run()
}

// RunCheck validates compilation, Jandex discovery, dependency resolution, and
// AOT source compilation without packaging an application artifact. It uses a
// clean target by default for the same reason as RunBuild.
func RunCheck(options BuildOptions, extraArgs ...string) error {
	if err := checkMaven(); err != nil {
		return err
	}

	args := buildArgs(options, "process-classes", extraArgs...)
	cmd := exec.Command("mvn", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if options.Incremental {
		fmt.Println("Checking Summer application incrementally...")
	} else {
		fmt.Println("Checking Summer application from a clean target...")
	}
	return cmd.Run()
}

func buildArgs(options BuildOptions, goal string, extraArgs ...string) []string {
	args := make([]string, 0, len(extraArgs)+3)
	if !options.Incremental {
		args = append(args, "clean")
	}
	args = append(args, goal)
	args = append(args, extraArgs...)
	if options.Incremental {
		// The Maven plugin uses this property to run stale-class cleanup before
		// javac. Keep it coupled to the CLI mode instead of relying on users to
		// remember a second -D flag.
		args = append(args, "-Dsummer.build.incremental=true")
	}
	return args
}
