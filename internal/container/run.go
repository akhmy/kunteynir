package container

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type InitError struct {
	Msg string
}

func (e *InitError) Error() string {
	return "container init: " + e.Msg
}

type ExitError struct {
	Code int
}

func (e *ExitError) Error() string {
	return fmt.Sprintf("container exited with status %d", e.Code)
}

func Run(lowerDir, hostname string, args, env []string) error {
	r, w, err := os.Pipe()
	if err != nil {
		return fmt.Errorf("create init error pipe: %w", err)
	}
	defer r.Close()

	args = append([]string{"setup", "-name", hostname, lowerDir}, args...)

	cmd := exec.Command("/proc/self/exe", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWIPC,
	}
	cmd.ExtraFiles = append(cmd.ExtraFiles, w)

	if err := cmd.Start(); err != nil {
		w.Close()
		return fmt.Errorf("start init: %w", err)
	}
	w.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("read init error pipe: %w", err)
	}

	if len(data) != 0 {
		cmd.Wait()
		return &InitError{Msg: strings.TrimSpace(string(data))}
	}

	err = cmd.Wait()

	if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
		status := exitErr.Sys().(syscall.WaitStatus)
		if status.Signaled() {
			return &ExitError{Code: 128 + int(status.Signal())}
		}
		return &ExitError{Code: status.ExitStatus()}
	}

	if err != nil {
		return fmt.Errorf("wait init: %w", err)
	}

	return nil
}
