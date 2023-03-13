package argp

import (
	"flag"
	"fmt"
	"strings"
)

// Command describe a command
type Command interface {
	// Run execute a command and return error
	Run(args []string) error
	// Help returns help message of a command
	Help() string
	// Name returns name of command
	Name() string
	// Describe describe the command itself, as short help message
	Describe() string
}

// Cmd is am implementation of Comamnd, it can add subcommand to it, command flags will always parsed whether
// sub command is executed or not.
// When sub command has been triggerd, the Command() won't be execute, but sub command will be executed instead
type Cmd[OptType any] struct {
	CmdName     string
	Usage       string
	Options     *OptType
	Command     func(opts *OptType) error
	SubCommands map[string]Command

	flagSet *flag.FlagSet
}

func (cmd *Cmd[OptType]) AddSubCmd(subcmd Command) {
	if subcmd == nil {
		panic("adding nil command")
	}
	if cmd.SubCommands == nil {
		cmd.SubCommands = make(map[string]Command)
	}
	if _, ok := cmd.SubCommands[subcmd.Name()]; ok {
		panic(fmt.Errorf("sub-command %q for command %q has already been registerd", subcmd.Name(), cmd.CmdName))
	}
	cmd.SubCommands[subcmd.Name()] = subcmd
}

func (cmd *Cmd[OptType]) Run(args []string) error {
	// parse cmds to options
	restArg, err := cmd.parseOptions(cmd.Options, args)
	if err != nil {
		return err
	}
	if len(restArg) >= 1 {
		// check sub command
		if subcmd, ok := cmd.SubCommands[restArg[0]]; ok {
			return subcmd.Run(restArg[1:])
		}
	}
	return cmd.Command(cmd.Options)
}

func (cmd *Cmd[OptType]) Help() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("Command %s:\n%s\n", cmd.Name(), cmd.Usage))
	cmd.flagSet.SetOutput(&s)
	cmd.flagSet.Usage()
	if len(cmd.SubCommands) > 0 {
		s.WriteString("sub-commands:\n")
		for _, cmd := range cmd.SubCommands {
			s.WriteString(fmt.Sprintf("\t%-25s%s\n", cmd.Name(), cmd.Describe()))
		}
	}
	return s.String()
}

func (cmd *Cmd[OptType]) Name() string {
	return cmd.CmdName
}

func (cmd *Cmd[OptType]) Describe() string {
	return cmd.Usage
}
