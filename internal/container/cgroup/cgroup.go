package cgroup

import (
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

const baseCgroup = "/sys/fs/cgroup/kunteynir"

type Cgroup struct {
	Path string
}

func New(name string) (*Cgroup, error) {
	if err := os.MkdirAll(baseCgroup, 0o755); err != nil {
		return nil, fmt.Errorf("create base cgroup: %w", err)
	}

	if err := os.WriteFile(
		filepath.Join(baseCgroup, "cgroup.subtree_control"),
		[]byte("+memory +pids +cpu"),
		0o644,
	); err != nil {
		return nil, fmt.Errorf("write subtree_control: %w", err)
	}

	containerCgroup := filepath.Join(baseCgroup, name)
	if err := os.Mkdir(containerCgroup, 0o755); err != nil {
		return nil, fmt.Errorf("create %q cgroup: %w", name, err)
	}

	cg := &Cgroup{Path: containerCgroup}
	if err := cg.writeLimit("memory.swap.max", "0"); err != nil {
		cg.Cleanup()
		return nil, fmt.Errorf("set swap to 0: %w", err)
	}

	return cg, nil
}

func (cg *Cgroup) SetMaxMemory(value string) error {
	return cg.writeLimit("memory.max", value)
}

func (cg *Cgroup) SetMaxPids(value string) error {
	return cg.writeLimit("pids.max", value)
}

const cpuPeriod = 100_000

func (cg *Cgroup) SetMaxCPU(quota float64) error {
	return cg.writeLimit("cpu.max", strconv.Itoa(int(math.Round(quota*cpuPeriod))))
}

func (cg *Cgroup) writeLimit(name, value string) error {
	file := filepath.Join(cg.Path, name)
	if err := os.WriteFile(file, []byte(value), 0o644); err != nil {
		return fmt.Errorf("write %q limit: %w", name, err)
	}

	return nil
}

const cleanupRetries = 5

const cooldown = time.Millisecond * 10

var ErrX = errors.New("x")

func (cg *Cgroup) Cleanup() error {
	for range cleanupRetries {
		err := os.Remove(cg.Path)
		if err == nil {
			return nil
		}
		if !errors.Is(err, syscall.EBUSY) {
			return fmt.Errorf("remove cgroup: %w", err)
		}

		time.Sleep(cooldown)
	}

	return ErrX
}

func (cg *Cgroup) Kill() error {
	killpath := filepath.Join(cg.Path, "cgroup.kill")
	if err := os.WriteFile(killpath, []byte{'1'}, 0); err != nil {
		return fmt.Errorf("write to %q: %w", killpath, err)
	}

	return nil
}
