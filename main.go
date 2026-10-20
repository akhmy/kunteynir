package main

import (
	"errors"
	"fmt"
	"os"
	"syscall"

	"github.com/akhmy/kunteynir/internal/cli"
	"github.com/akhmy/kunteynir/internal/container"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "%s\n\n", cli.ErrNoCommand)
		cmd, _ := cli.NewUsageCommand(os.Stderr)
		cmd.Execute()
		os.Exit(2)
	}

	if os.Args[1] == "init" {
		pipe := os.NewFile(3, "init-error-pipe")
		syscall.CloseOnExec(3)

		cmd, err := cli.NewInitCommand(os.Args[2:])
		if err != nil {
			reportInitError(pipe, err)
			os.Exit(2)
		}

		if err := cmd.Execute(); err != nil {
			reportInitError(pipe, err)
			os.Exit(1)
		}

		os.Exit(0)
	}

	cmd, err := cli.NewCommand(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	err = cmd.Execute()

	if exitErr, ok := errors.AsType[*container.ExitError](err); ok {
		os.Exit(exitErr.Code)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "kunteynir: %s\n", err)
		os.Exit(1)
	}
}

// reportInitError sends err to kunteynir run through the pipe,
// falling back to stderr when init was started without one.
func reportInitError(pipe *os.File, err error) {
	if _, werr := fmt.Fprint(pipe, err); werr != nil {
		fmt.Fprintf(os.Stderr, "kunteynir: %s\n", err)
	}
}
