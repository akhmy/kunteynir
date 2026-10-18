package cli

import (
	"errors"
	"os"

	"github.com/akhmy/kunteynir/internal/container"
)

type InitCommand struct {
	LowerDir string
	Command  string
}

var ErrInitUsage = errors.New("usage: kunteynir init <lowerdir> <command>")

var ErrNotDir = errors.New("lowerdir should be a directory")

func (c *InitCommand) Execute() error {
	info, err := os.Stat(c.LowerDir)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return ErrNotDir
	}

	return container.Init(c.LowerDir, c.Command)
}
