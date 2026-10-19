package cli

import (
	"errors"
	"os"

	"github.com/akhmy/kunteynir/internal/container"
)

var ErrRunUsage = errors.New("usage: kunteynir run <lowerdir> <command>")

type RunCommand struct {
	LowerDir string
	Command  string
}

func (c *RunCommand) Execute() error {
	info, err := os.Stat(c.LowerDir)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return ErrNotDir
	}

	container.Run(c.LowerDir, c.Command)

	return nil
}
