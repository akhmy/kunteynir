package cli

import (
	"fmt"
	"io"
)

type UsageCommand struct {
	out io.Writer
}

func NewUsageCommand(out io.Writer) (*UsageCommand, error) {
	return &UsageCommand{out: out}, nil
}

func (c *UsageCommand) Execute() error {
	fmt.Fprintln(c.out, "usage: kunteynir <command>")

	return nil
}
