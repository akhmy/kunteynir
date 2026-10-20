package container

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

const (
	tempDir   = "/tmp"
	upperDir  = tempDir + "/upper"
	workDir   = tempDir + "/work"
	mergedDir = tempDir + "/merged"

	// pivot_root moves the old root to putOld, which is
	// reachable as oldRoot from the new root.
	putOld  = mergedDir + "/.oldroot"
	oldRoot = "/.oldroot"
)

func Init(lowerDir, hostname string, args []string) error {
	err := syscall.Mount(
		"", "/", "",
		syscall.MS_PRIVATE|syscall.MS_REC, "",
	)
	if err != nil {
		return fmt.Errorf("make mounts private: %w", err)
	}

	if err := setUpFS(lowerDir); err != nil {
		return fmt.Errorf("init fs: %w", err)
	}

	err = syscall.Mount("proc", "/proc", "proc", 0, "")
	if err != nil {
		return fmt.Errorf("mount /proc: %w", err)
	}

	if err = syscall.Sethostname([]byte(hostname)); err != nil {
		return fmt.Errorf("sethostname: %w", err)
	}

	exe, err := exec.LookPath(args[0])
	if err != nil {
		return err
	}

	err = syscall.Exec(exe, args, os.Environ())
	if err != nil {
		return fmt.Errorf("exec %s: %w", exe, err)
	}

	return nil
}

func setUpFS(lowerDir string) error {
	lowerDir, err := filepath.Abs(lowerDir)
	if err != nil {
		return fmt.Errorf("abs lowerdir: %w", err)
	}

	if err := syscall.Mount("tmpfs", tempDir, "tmpfs", 0, ""); err != nil {
		return fmt.Errorf("mount tmpfs: %w", err)
	}

	if err := os.Mkdir(upperDir, 0755); err != nil {
		return fmt.Errorf("create upperdir: %w", err)
	}
	if err := os.Mkdir(workDir, 0755); err != nil {
		return fmt.Errorf("create workdir: %w", err)
	}
	if err := os.Mkdir(mergedDir, 0755); err != nil {
		return fmt.Errorf("create mergeddir: %w", err)
	}

	if err := syscall.Mount("overlay", mergedDir, "overlay", 0,
		fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", lowerDir, upperDir, workDir),
	); err != nil {
		return fmt.Errorf("mount overlay: %w", err)
	}

	if err := os.Mkdir(putOld, 0755); err != nil {
		return fmt.Errorf("create oldrootdir: %w", err)
	}

	if err := syscall.PivotRoot(mergedDir, putOld); err != nil {
		return fmt.Errorf("pivot root: %w", err)
	}

	if err := os.Chdir("/"); err != nil {
		return fmt.Errorf("chdir root: %w", err)
	}

	if err := syscall.Unmount(oldRoot, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount oldrootdir: %w", err)
	}

	if err := os.Remove(oldRoot); err != nil {
		return fmt.Errorf("remove oldrootdir: %w", err)
	}

	return nil
}
