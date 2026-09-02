//go:build !windows

package runner

import (
	"os/exec"
	"syscall"
)

// setProcessGroup puts Maven into its own process group so a cancel signal can
// reach the whole tree (mvn spawns the application JVM as a child in dev mode).
func setProcessGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

// signalProcessGroup signals Maven's process group with SIGTERM, the graceful
// termination signal the JVM's shutdown hook expects. Called from cmd.Cancel
// when the CLI's own context is cancelled (SIGINT/SIGTERM on the CLI).
func signalProcessGroup(cmd *exec.Cmd) error {
	if cmd.Process == nil {
		return nil
	}
	if pgid, err := syscall.Getpgid(cmd.Process.Pid); err == nil && pgid == cmd.Process.Pid {
		return syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)
	}
	return cmd.Process.Signal(syscall.SIGTERM)
}
