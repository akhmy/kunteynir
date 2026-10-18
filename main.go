package main

import (
	"fmt"
	"os"

	"github.com/akhmy/kunteynir/internal/cli"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, "no command\n\n")
		(&cli.UsageCommand{Out: os.Stderr}).Execute()
		os.Exit(2)
	}

	factory := &cli.CommandFactory{}

	cmd, err := factory.New(os.Args[1], os.Args[2:]...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}
	err = cmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute command: %s\n", err.Error())
		os.Exit(1)
	}
}
