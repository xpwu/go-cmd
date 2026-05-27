package cmd

import (
	"errors"
	"fmt"
	"github.com/xpwu/go-cmd/arg"
	"os"
	"sort"
	"strings"
)

type Cmd func(args *arg.Arg)

func format(len, maxLen int) string {
	format := "        %s"
	for i := maxLen - len + 2; i > 0; i-- {
		format += " "
	}
	format += ": %s\n"

	return format
}

func usage(builder *strings.Builder) {
	_, _ = fmt.Fprintf(builder,
		"\nUsage:\n\n        %s <command> [arguments]\n\nThe commands are: (the default command is '%s') \n\n",
		os.Args[0], KeepAliveCmd)

	maxLen := 0
	keys := make([]string, 0, len(helps))
	for k, _ := range helps {
		if len(k) > maxLen {
			maxLen = len(k)
		}
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for _, k := range keys {
		_, _ = fmt.Fprintf(builder, format(len(k), maxLen), k, helps[k])
	}

	_, _ = fmt.Fprintf(builder,
		"\nUse \"%s <command> -h\" for more information about the command.\n\n", os.Args[0])

}

var helpCmds = map[string]struct{}{
	"-h":     {},
	"-help":  {},
	"--help": {},
	"help":   {},
}

var cmds = map[string]Cmd{}

var helps = map[string]string{}

// Deprecated: DefaultCmdName using: KeepAliveCmd
const DefaultCmdName = KeepAliveCmd

// KeepAliveCmd keep process alive and  is also default cmd
const KeepAliveCmd = "run"

func RegisterCmd(cmdName string, help string, cmd Cmd) {
	err := RegisterCmdErr(cmdName, help, cmd)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}
}

func RegisterCmdErr(cmdName string, help string, cmd Cmd) error {
	// replace KeepAliveCmd directly
	if cmdName == KeepAliveCmd {
		cmds[KeepAliveCmd] = func(args *arg.Arg) {
			cmd(args)

			// keep process alive
			exit := make(chan struct{})
			<-exit
		}
		helps[KeepAliveCmd] = help
		return nil
	}

	if cmdName[0] == '-' {
		return errors.New("Error: cmdName can NOT start with '-'\n")
	}

	if _, ok := helpCmds[cmdName]; ok {
		return errors.New("Error: Not Register 'help' cmd\n")
	}

	tryName := cmdName
	for i := 0; ; i++ {
		_, ok := cmds[tryName]
		if !ok {
			break
		}

		tryName = fmt.Sprintf("%s%d", cmdName, i)
	}
	cmds[tryName] = cmd
	helps[tryName] = help

	return nil
}

func RegisterCmdNoArgs(cmdName string, help string, cmd func()) {
	RegisterCmd(cmdName, help, func(args *arg.Arg) {
		args.ParseAndRunHook()
		cmd()
	})
}

func fail(err string) error {
	builder := &strings.Builder{}
	_, _ = fmt.Fprintln(builder, err)
	usage(builder)
	return errors.New(builder.String())
}

type HelpErr error

func help() HelpErr {
	builder := &strings.Builder{}
	usage(builder)
	return errors.New(builder.String())
}

func Run() {
	err := RunErr()
	_, _ = fmt.Fprintln(os.Stderr, err.Error())
	if _, ok := err.(HelpErr); ok {
		os.Exit(0)
	}
	os.Exit(2)
}

func RunErr() error {

	args := os.Args[1:]
	if len(args) == 0 {
		if cmd, ok := cmds[KeepAliveCmd]; ok {
			cmd(arg.NewArg(os.Args[0], args))
			return nil
		}

		return fail("Error: No command specified or NOT Register 'run'(default) command")
	}

	cmdName := args[0]

	if _, ok := helpCmds[cmdName]; ok {
		return help()
	}

	// cmdName is an argument
	if cmdName[0] == '-' {
		if cmd, ok := cmds[KeepAliveCmd]; ok {
			cmd(arg.NewArg(os.Args[0], args))
			return nil
		}

		return fail("Error: NOT Register 'run'(default) command")
	}

	if cmd, ok := cmds[cmdName]; ok {
		cmd(arg.NewArg(os.Args[0]+" "+cmdName+" ", args[1:]))
		return nil
	}

	return fail("Error: NOT Register '" + cmdName + "' command")
}
