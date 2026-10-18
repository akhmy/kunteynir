package container

import (
	"fmt"
	"os"
	"syscall"
)

func Init(lowerDir, command string) error {
	err := syscall.Sethostname([]byte("kunteynir"))
	if err != nil {
		return fmt.Errorf("sethostname: %w", err)
	}

	err = syscall.Exec(command, []string{command}, os.Environ())
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}
