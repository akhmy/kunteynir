package fs

import (
	"fmt"
	"syscall"
)

var devices = []struct {
	path         string
	major, minor int
}{
	{"/dev/null", 1, 3},
	{"/dev/zero", 1, 5},
	{"/dev/full", 1, 7},
	{"/dev/random", 1, 8},
	{"/dev/urandom", 1, 9},
	{"/dev/tty", 5, 0},
}

func setupDevices() error {
	old := syscall.Umask(0)
	defer syscall.Umask(old)

	for _, d := range devices {
		if err := syscall.Mknod(d.path, syscall.S_IFCHR|0o666, d.major<<8|d.minor); err != nil {
			return fmt.Errorf("mknod %s: %w", d.path, err)
		}
	}

	return nil
}
