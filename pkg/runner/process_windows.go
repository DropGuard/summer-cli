//go:build windows

package runner

import (
	"os/exec"
	"syscall"
)

// Windows has no process groups via os/exec's portable surface; cancellation
// degrades to signalling the direct child.
func setProcessGroup(cmd *exec.Cmd) {}

func signalProcessGroup(cmd *exec.Cmd) error {
	if cmd.Process == nil {
		return nil
	}
	return cmd.Process.Signal(syscall.SIGTERM)
}
