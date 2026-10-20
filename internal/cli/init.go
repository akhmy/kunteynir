package cli

import (
	"errors"
	"flag"
	"os"

	"github.com/akhmy/kunteynir/internal/container"
)

type InitCommand struct {
	lowerDir string
	hostname string
	args     []string
}

var ErrInitUsage = errors.New("usage: kunteynir init <lowerdir> <command> [args...]")

var ErrNotContainerInit = errors.New("init must run as PID 1 of a container, use kunteynir run")

func NewInitCommand(args []string) (*InitCommand, error) {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	hostname := fs.String("name", "default_name", "usage placeholder")

	fs.Parse(args)

	args = fs.Args()

	if len(args) < 2 {
		return nil, ErrInitUsage
	}

	return &InitCommand{
		lowerDir: args[0],
		hostname: *hostname,
		args:     args[1:],
	}, nil
}

func (c *InitCommand) Execute() error {
	if os.Getpid() != 1 {
		return ErrNotContainerInit
	}

	return container.Init(c.lowerDir, c.hostname, c.args)
}
