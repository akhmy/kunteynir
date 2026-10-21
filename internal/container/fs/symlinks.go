package fs

import (
	"fmt"
	"os"
)

var symlinks = []struct{ target, link string }{
	{"/proc/self/fd", "/dev/fd"},
	{"/proc/self/fd/0", "/dev/stdin"},
	{"/proc/self/fd/1", "/dev/stdout"},
	{"/proc/self/fd/2", "/dev/stderr"},
	{"/dev/pts/ptmx", "/dev/ptmx"},
}

func setupSymlinks() error {
	for _, s := range symlinks {
		if err := os.Symlink(s.target, s.link); err != nil {
			return fmt.Errorf("symlink %s: %w", s.link, err)
		}
	}

	return nil
}
