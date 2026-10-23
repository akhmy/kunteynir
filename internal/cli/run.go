package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/akhmy/kunteynir/internal/container"
	"github.com/akhmy/kunteynir/internal/pkg"
)

var ErrRunUsage = errors.New("usage: kunteynir run [-e KEY[=VALUE]] [-name NAME] <lowerdir> <command> [args...]")

var ErrNotDir = errors.New("not a directory")

type RunCommand struct {
	lowerDir string
	hostname string
	args     []string
	env      []string

	memory string
	pids   int
	cpus   float64
}

const (
	pathEnv = "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
)

func NewRunCommand(args []string) (*RunCommand, error) {
	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	var env pkg.EnvValue

	if err := env.Set(pathEnv); err != nil {
		return nil, fmt.Errorf("set path env: %w", err)
	}

	hostname := fs.String("name", "default_hostname", "--name flag usage placeholder")
	memory := fs.String("memory", "max", "--memory flag usage placeholder")
	pids := fs.Int("pids", -1, "--pids flag usage placeholder")
	cpus := fs.Float64("cpus", -1.0, "--cpus flag usage placeholder")
	fs.Var(&env, "e", "-e flag usage placeholder")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	args = fs.Args()

	if len(args) < 2 {
		return nil, ErrRunUsage
	}

	return &RunCommand{
		lowerDir: args[0],
		hostname: *hostname,
		args:     args[1:],
		env:      env.Env(),

		memory: *memory,
		pids:   *pids,
		cpus:   *cpus,
	}, nil
}

func (c *RunCommand) Execute() error {
	info, err := os.Stat(c.lowerDir)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("lowerdir %s: %w", c.lowerDir, ErrNotDir)
	}

	return container.Run(
		c.lowerDir,
		c.hostname,
		c.args,
		c.env,
		c.memory,
		c.pids,
		c.cpus,
	)
}
