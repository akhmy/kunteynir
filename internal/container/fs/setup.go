package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

func Setup(lowerDir string) error {
	if err := syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("make mounts private: %w", err)
	}

	if err := setupRootFS(lowerDir); err != nil {
		return fmt.Errorf("set up root fs: %w", err)
	}

	if err := setupPseudoFS(); err != nil {
		return fmt.Errorf("set up pseudo fs: %w", err)
	}

	if err := setupDevices(); err != nil {
		return fmt.Errorf("set up devices: %w", err)
	}

	if err := setupSymlinks(); err != nil {
		return fmt.Errorf("set up symlinks: %w", err)
	}

	return nil
}

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

func setupRootFS(lowerDir string) error {
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
