//go:build !windows

package runner

import (
	"context"
	"os/exec"
	"syscall"
)

// runDevProcess puts Maven into its own process group so a stop reaches the
// whole tree (mvn spawns the application JVM as a child in dev mode). A cancel
// sends SIGTERM to the group; if Maven is still running devStopTimeout later,
// WaitDelay kills it. Either way, once Maven has exited the rest of its group
// is SIGKILLed — the guarantee the Windows job object gives with
// KILL_ON_JOB_CLOSE: nothing Maven started outlives the CLI.
func runDevProcess(_ context.Context, stop context.CancelFunc, cmd *exec.Cmd) error {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		stop()
		return signalProcessGroup(cmd)
	}
	cmd.WaitDelay = devStopTimeout
	err := cmd.Run()
	killProcessGroup(cmd)
	return err
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

// killProcessGroup SIGKILLs whatever is left of Maven's process group after
// Maven itself has been reaped. A group ID is not recycled while any member is
// alive, so a live group is still Maven's; an empty one answers ESRCH, which is
// ignored. (The only gap is a brand-new group leader taking Maven's PID in the
// microseconds between reap and kill — pid allocation would have to wrap.)
func killProcessGroup(cmd *exec.Cmd) {
	if cmd.Process == nil {
		return
	}
	_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
}
