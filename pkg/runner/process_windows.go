package runner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

func runDevProcess(ctx context.Context, stop context.CancelFunc, cmd *exec.Cmd) error {
	if cmd.Err != nil {
		return cmd.Err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	job, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return fmt.Errorf("create Maven process job: %w", err)
	}
	defer windows.CloseHandle(job)
	var limits windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION
	limits.BasicLimitInformation.LimitFlags = windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE
	if _, err := windows.SetInformationJobObject(job, windows.JobObjectExtendedLimitInformation,
		uintptr(unsafe.Pointer(&limits)), uint32(unsafe.Sizeof(limits))); err != nil {
		return fmt.Errorf("configure Maven process job: %w", err)
	}
	process, err := startInJob(job, cmd)
	if err != nil {
		return fmt.Errorf("start Maven process job: %w", err)
	}
	finished := make(chan error, 1)
	go func() {
		state, err := process.Wait()
		if err == nil && !state.Success() {
			err = &exec.ExitError{ProcessState: state}
		}
		finished <- err
	}()
	select {
	case err := <-finished:
		return err
	case <-ctx.Done():
		stop()
	}
	return awaitConsoleShutdown(job, finished)
}

// Maven shares the CLI console, so keyboard interrupts reach both; the job enforces the deadline.
func awaitConsoleShutdown(job windows.Handle, finished <-chan error) error {
	timer := time.NewTimer(devStopTimeout)
	defer timer.Stop()
	select {
	case err := <-finished:
		return err
	case <-timer.C:
		if err := windows.TerminateJobObject(job, 130); err != nil {
			return fmt.Errorf("terminate Maven process job: %w", err)
		}
		return <-finished
	}
}

// stdioFiles returns the command's standard streams as files: their handles are
// passed explicitly to CreateProcess, so a non-file Reader/Writer is rejected
// with an error instead of a panicking type assertion.
func stdioFiles(cmd *exec.Cmd) ([3]*os.File, error) {
	var files [3]*os.File
	for index, stream := range []any{cmd.Stdin, cmd.Stdout, cmd.Stderr} {
		file, ok := stream.(*os.File)
		if !ok || file == nil {
			return files, fmt.Errorf("dev mode requires file-backed standard streams (stream %d is %T)", index, stream)
		}
		files[index] = file
	}
	return files, nil
}

// Start Maven suspended so it cannot spawn a JVM before it belongs to the job.
func startInJob(job windows.Handle, cmd *exec.Cmd) (*os.Process, error) {
	path, err := windows.UTF16PtrFromString(cmd.Path)
	if err != nil {
		return nil, err
	}
	args := append([]string{cmd.Path}, cmd.Args[1:]...)
	commandLine, err := windows.UTF16PtrFromString(windows.ComposeCommandLine(args))
	if err != nil {
		return nil, err
	}
	stdio, err := stdioFiles(cmd)
	if err != nil {
		return nil, err
	}
	var handles [3]windows.Handle
	for index, file := range stdio {
		if err := windows.DuplicateHandle(windows.CurrentProcess(), windows.Handle(file.Fd()),
			windows.CurrentProcess(), &handles[index], 0, true, windows.DUPLICATE_SAME_ACCESS); err != nil {
			return nil, err
		}
		defer windows.CloseHandle(handles[index])
	}
	attrs, err := windows.NewProcThreadAttributeList(1)
	if err != nil {
		return nil, err
	}
	defer attrs.Delete()
	if err := attrs.Update(windows.PROC_THREAD_ATTRIBUTE_HANDLE_LIST, unsafe.Pointer(&handles[0]), unsafe.Sizeof(handles)); err != nil {
		return nil, err
	}
	startup := windows.StartupInfoEx{
		StartupInfo: windows.StartupInfo{
			Flags:     windows.STARTF_USESTDHANDLES,
			StdInput:  handles[0],
			StdOutput: handles[1],
			StdErr:    handles[2],
		},
		ProcThreadAttributeList: attrs.List(),
	}
	startup.Cb = uint32(unsafe.Sizeof(startup))
	var child windows.ProcessInformation
	flags := uint32(windows.CREATE_SUSPENDED | windows.EXTENDED_STARTUPINFO_PRESENT)
	if err := windows.CreateProcess(path, commandLine, nil, nil, true, flags, nil, nil, &startup.StartupInfo, &child); err != nil {
		return nil, err
	}
	defer windows.CloseHandle(child.Thread)
	defer windows.CloseHandle(child.Process)
	if err := windows.AssignProcessToJobObject(job, child.Process); err != nil {
		_ = windows.TerminateProcess(child.Process, 1)
		return nil, err
	}
	process, err := os.FindProcess(int(child.ProcessId))
	if err != nil {
		return nil, err
	}
	if _, err := windows.ResumeThread(child.Thread); err != nil {
		_ = process.Release()
		return nil, err
	}
	return process, nil
}
