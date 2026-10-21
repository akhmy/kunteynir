package fs

import (
	"fmt"
	"os"
	"syscall"
)

var pseudoMounts = []struct {
	target, fstype string
	flags          uintptr
	data           string
}{
	{"/proc", "proc", syscall.MS_NOSUID | syscall.MS_NODEV | syscall.MS_NOEXEC, ""},
	{"/sys", "sysfs", syscall.MS_RDONLY | syscall.MS_NOSUID | syscall.MS_NODEV | syscall.MS_NOEXEC, ""},
	{"/dev", "tmpfs", syscall.MS_NOSUID | syscall.MS_STRICTATIME, "mode=755,size=64m"},
	{"/dev/pts", "devpts", syscall.MS_NOSUID | syscall.MS_NOEXEC, "newinstance,ptmxmode=0666,mode=0620,gid=5"},
	{"/dev/shm", "tmpfs", syscall.MS_NOSUID | syscall.MS_NODEV | syscall.MS_NOEXEC, "mode=1777,size=64m"},
	{"/dev/mqueue", "mqueue", syscall.MS_NOSUID | syscall.MS_NODEV | syscall.MS_NOEXEC, ""},
}

func setupPseudoFS() error {
	for _, m := range pseudoMounts {
		if err := os.MkdirAll(m.target, 0755); err != nil {
			return err
		}
		if err := syscall.Mount(
			m.fstype,
			m.target,
			m.fstype,
			m.flags,
			m.data,
		); err != nil {
			return fmt.Errorf("mount %s: %w", m.target, err)
		}
	}

	return nil
}
