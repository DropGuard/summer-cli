package runner

import (
	"fmt"
	"os"
	"os/exec"
)

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

func RunBuild(extraArgs ...string) error {
	if err := checkMaven(); err != nil {
		return err
	}

	args := append([]string{"clean", "package"}, extraArgs...)
	cmd := exec.Command("mvn", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("📦 Building Summer application...")
	return cmd.Run()
}
