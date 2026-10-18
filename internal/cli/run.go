package cli

import (
	"errors"
	"fmt"

	"github.com/akhmy/kunteynir/internal/container"
)

var ErrRunUsage = errors.New("usage: kunteynir run <lowerdir> <command>")

type RunCommand struct {
	LowerDir string
	Command  string
}

func (c *RunCommand) Execute() error {
	err := container.Run(c.LowerDir, c.Command)
	if err != nil {
		return fmt.Errorf("run command: %w", err)
	}

	return nil
}
