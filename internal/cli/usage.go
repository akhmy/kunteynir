package cli

import (
	"fmt"
	"io"
)

type UsageCommand struct {
	Out io.Writer
}

func (c *UsageCommand) Execute() error {
	fmt.Fprintln(c.Out, "usage: kunteynir <command>")

	return nil
}
