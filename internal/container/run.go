package container

import (
	"os"
	"os/exec"
	"syscall"
)

func Run(lowerDir, command string) error {
	cmd := exec.Command("/proc/self/exe", "init", lowerDir, command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS,
	}

	cmd.Run()

	return nil
}
