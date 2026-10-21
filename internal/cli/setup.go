package cli

import (
	"errors"
	"flag"
	"os"

	"github.com/akhmy/kunteynir/internal/container"
)

type SetupCommand struct {
	lowerDir string
	hostname string
	args     []string
}

var ErrSetupUsage = errors.New("usage: kunteynir setup <lowerdir> <command> [args...]")

var ErrNotContainerSetup = errors.New("setup must run as PID 1 of a container, use kunteynir run")

func NewSetupCommand(args []string) (*SetupCommand, error) {
	fs := flag.NewFlagSet("setup", flag.ContinueOnError)
	hostname := fs.String("name", "default_name", "usage placeholder")

	fs.Parse(args)

	args = fs.Args()

	if len(args) < 2 {
		return nil, ErrSetupUsage
	}

	return &SetupCommand{
		lowerDir: args[0],
		hostname: *hostname,
		args:     args[1:],
	}, nil
}

func (c *SetupCommand) Execute() error {
	if os.Getpid() != 1 {
		return ErrNotContainerSetup
	}

	return container.Setup(c.lowerDir, c.hostname, c.args)
}
