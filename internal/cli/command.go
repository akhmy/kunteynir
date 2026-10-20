package cli

import (
	"errors"
	"fmt"
	"os"
)

type Command interface {
	Execute() error
}

type UnknownCommandError struct {
	Command string
}

func (e *UnknownCommandError) Error() string {
	return fmt.Sprintf("unknown command %q", e.Command)
}

var ErrNoCommand = errors.New("no command")

func NewCommand(args []string) (Command, error) {
	if len(args) == 0 {
		return nil, ErrNoCommand
	}

	switch args[0] {
	case "help":
		return NewUsageCommand(os.Stdout)
	case "run":
		return NewRunCommand(args[1:])
	default:
		return nil, &UnknownCommandError{Command: args[0]}
	}
}
