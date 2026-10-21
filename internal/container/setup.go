package container

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/akhmy/kunteynir/internal/container/fs"
)

func Setup(lowerDir, hostname string, args []string) error {
	if err := fs.Setup(lowerDir); err != nil {
		return fmt.Errorf("set up fs: %w", err)
	}

	if err := syscall.Sethostname([]byte(hostname)); err != nil {
		return fmt.Errorf("set hostname: %w", err)
	}

	arg0, err := exec.LookPath(args[0])
	if err != nil {
		return err
	}

	err = syscall.Exec(arg0, args, os.Environ())
	if err != nil {
		return fmt.Errorf("exec %s: %w", arg0, err)
	}

	return nil
}
