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

func RunDev() error {
	if err := checkMaven(); err != nil {
		return err
	}
	
	cmd := exec.Command("mvn", "summer:dev")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	
	fmt.Println("🚀 Starting Summer Framework in Dev Mode...")
	return cmd.Run()
}

func RunBuild() error {
	if err := checkMaven(); err != nil {
		return err
	}
	
	cmd := exec.Command("mvn", "clean", "package")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	fmt.Println("📦 Building Summer application...")
	return cmd.Run()
}
