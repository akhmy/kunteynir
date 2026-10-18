package cli

import (
	"fmt"
	"io"
	"os"
)

type UsageCommand struct {
	Out io.Writer
}

func (c *UsageCommand) Execute() error {
	fmt.Fprintln(c.Out, "usage: kunteynir <command>")

	return nil
}

type UnknownWrapper struct {
	Value string
	Out   io.Writer
}

func (c *UnknownWrapper) Execute() error {
	fmt.Fprintln(c.Out, "unknown command: "+c.Value+"\n")

	return os.ErrNotExist
}
