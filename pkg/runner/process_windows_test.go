package runner

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestStdioFilesAcceptsFileBackedStreams(t *testing.T) {
	cmd := exec.Command("cmd.exe")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	files, err := stdioFiles(cmd)
	if err != nil {
		t.Fatalf("stdioFiles: %v", err)
	}
	if files != [3]*os.File{os.Stdin, os.Stdout, os.Stderr} {
		t.Fatalf("stdioFiles returned %v, want the command's own streams", files)
	}
}

func TestStdioFilesRejectsNonFileStreams(t *testing.T) {
	cases := map[string]func(*exec.Cmd){
		"nil stdin":     func(cmd *exec.Cmd) { cmd.Stdin = nil },
		"buffer stdout": func(cmd *exec.Cmd) { cmd.Stdout = &bytes.Buffer{} },
		"nil-file stderr": func(cmd *exec.Cmd) {
			var file *os.File
			cmd.Stderr = file
		},
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			cmd := exec.Command("cmd.exe")
			cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
			mutate(cmd)
			_, err := stdioFiles(cmd)
			if err == nil || !strings.Contains(err.Error(), "file-backed standard streams") {
				t.Fatalf("stdioFiles error = %v, want a file-backed streams error", err)
			}
		})
	}
}
