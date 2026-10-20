package cli

import "os"

type Command interface {
	Execute() error
}

type UnknownCommandError struct {
	Command string
}

func (e *UnknownCommandError) Error() string {
	return "unknown command: " + e.Command
}

func NewCommand(cmd string, args ...string) (Command, error) {
	switch cmd {
	case "help":
		return &UsageCommand{Out: os.Stdout}, nil
	case "run":
		if len(args) != 2 {
			return nil, ErrRunUsage
		}
		return &RunCommand{LowerDir: args[0], Command: args[1]}, nil
	default:
		return nil, &UnknownCommandError{Command: cmd}
	}
}
