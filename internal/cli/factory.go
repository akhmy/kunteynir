package cli

import (
	"io"
)

type Command interface {
	Execute() error
}

type CommandFactory struct {
	In  io.Writer
	Out io.Writer
	Err io.Writer
}

type UnknownCommandError struct {
	Command string
}

func (e *UnknownCommandError) Error() string {
	return "unknown command: " + e.Command
}

func (cf *CommandFactory) New(cmd string, args ...string) (Command, error) {
	switch cmd {
	case "help":
		return &UsageCommand{Out: cf.Out}, nil
	case "run":
		if len(args) != 2 {
			return nil, ErrRunUsage
		}
		return &RunCommand{LowerDir: args[0], Command: args[1]}, nil
	case "init":
		if len(args) != 2 {
			return nil, ErrInitUsage
		}
		return &InitCommand{LowerDir: args[0], Command: args[1]}, nil
	default:
		return nil, &UnknownCommandError{Command: cmd}
	}
}
