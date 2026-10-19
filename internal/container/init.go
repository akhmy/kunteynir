package container

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

func Init(lowerDir, command string) error {
	err := syscall.Mount(
		"", "/", "",
		syscall.MS_PRIVATE|syscall.MS_REC, "",
	)
	if err != nil {
		return fmt.Errorf("change root propagation type: %w", err)
	}

	err = syscall.Mount(lowerDir, lowerDir, "", syscall.MS_BIND|syscall.MS_REC, "")
	if err != nil {
		return fmt.Errorf("mount lower directory: %w", err)
	}

	oldRootDir := filepath.Join(lowerDir, ".oldroot")

	err = os.MkdirAll(oldRootDir, 0700)
	if err != nil {
		return fmt.Errorf("create .oldroot directory: %w", err)
	}

	err = syscall.PivotRoot(lowerDir, oldRootDir)
	if err != nil {
		return fmt.Errorf("pivot root: %w", err)
	}

	err = os.Chdir("/")
	if err != nil {
		return fmt.Errorf("cd root: %w", err)
	}

	err = syscall.Unmount("/.oldroot", syscall.MNT_DETACH)
	if err != nil {
		return fmt.Errorf("unmount .oldroot directory: %w", err)
	}

	err = os.Remove("/.oldroot")
	if err != nil {
		return fmt.Errorf("remove .oldroot directory: %w", err)
	}

	err = syscall.Mount("proc", "/proc", "proc", 0, "")
	if err != nil {
		return fmt.Errorf("mount proc: %w", err)
	}

	err = syscall.Sethostname([]byte("kunteynir"))
	if err != nil {
		return fmt.Errorf("sethostname: %w", err)
	}

	err = syscall.Exec(command, []string{command}, os.Environ())
	if err != nil {
		return fmt.Errorf("exec command: %w", err)
	}

	return nil
}
