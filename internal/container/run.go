package container

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"uuid"

	"github.com/akhmy/kunteynir/internal/container/cgroup"
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
	cg, err := cgroup.New(uuid.New().String())
	if err != nil {
		return fmt.Errorf("create cgroup: %w", err)
	}

	defer cg.Kill()
	defer cg.Cleanup()

	if err := cg.SetMaxMemory("64m"); err != nil {
		return err
	}

	if err := cg.SetMaxCPU(0.5); err != nil {
		return err
	}

	if err := cg.SetMaxPids("20"); err != nil {
		return err
	}

	r, w, err := os.Pipe()
	if err != nil {
		return fmt.Errorf("create init error pipe: %w", err)
	}
	defer r.Close()

	args = append([]string{"setup", "-name", hostname, lowerDir}, args...)

	dir, err := os.Open(cg.Path)
	if err != nil {
		return fmt.Errorf("open cgroup: %w", err)
	}

	defer dir.Close()

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
		UseCgroupFD: true,
		CgroupFD:    int(dir.Fd()),
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
