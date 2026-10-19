package cli

import (
	"errors"

	"github.com/akhmy/kunteynir/internal/container"
)

var ErrRunUsage = errors.New("usage: kunteynir run <lowerdir> <command>")

type RunCommand struct {
	LowerDir string
	Command  string
}

func (c *RunCommand) Execute() error {
	container.Run(c.LowerDir, c.Command)

	return nil
}
