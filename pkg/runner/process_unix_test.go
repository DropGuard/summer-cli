//go:build !windows

package runner

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

// Maven exits while a child it started (in dev mode: the application JVM) is
// still running in its process group. The CLI must not leave that child behind
// — the Unix counterpart of the Windows job object's KILL_ON_JOB_CLOSE.
func TestRunDevProcessKillsGroupLeftoversAfterMavenExits(t *testing.T) {
	pidFile := filepath.Join(t.TempDir(), "child.pid")
	ctx, stop := context.WithCancel(context.Background())
	defer stop()
	cmd := exec.CommandContext(ctx, "sh", "-c", `sleep 60 >/dev/null 2>&1 & echo $! > "$1"`, "sh", pidFile)
	if err := runDevProcess(ctx, stop, cmd); err != nil {
		t.Fatalf("runDevProcess: %v", err)
	}
	data, err := os.ReadFile(pidFile)
	if err != nil {
		t.Fatalf("read child pid: %v", err)
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		t.Fatalf("parse child pid %q: %v", data, err)
	}
	deadline := time.Now().Add(3 * time.Second)
	for !processGone(pid) {
		if time.Now().After(deadline) {
			_ = syscall.Kill(pid, syscall.SIGKILL)
			t.Fatalf("child %d outlived Maven's exit", pid)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func processGone(pid int) bool {
	if err := syscall.Kill(pid, 0); errors.Is(err, syscall.ESRCH) {
		return true
	}
	// A killed but not yet reaped child still answers signal 0; /proc (Linux) tells it apart.
	stat, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	return err == nil && strings.Contains(string(stat), ") Z ")
}
