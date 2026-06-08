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
	format := "  %s"
	for i := maxLen - len + 2; i > 0; i-- {
		format += " "
	}
	format += "%s\n"

	return format
}

var exeNameInUsage = os.Args[0]

func usage(builder *strings.Builder) {
	_, _ = fmt.Fprintf(builder,
		"Usage: %s <command> [arguments]\nThe valid 'commands' are: (the default command is '%s')\n",
		exeNameInUsage, DefaultCmdName)

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
		"Every 'argument' starts with '-'.\n"+
			"Use \"%s <command> -h\" for more information about the command.\n", exeNameInUsage)

}

var helpCmds = map[string]struct{}{
	"-h":     {},
	"-help":  {},
	"--help": {},
	"help":   {},
}

var cmds = map[string]Cmd{}

var helps = map[string]string{}

const DefaultCmdName = "run"

var osExit = os.Exit

func RegisterKeepAliveCmd(cmdName string, help string, cmd Cmd) {
	err := RegisterKeepAliveCmdErr(cmdName, help, cmd)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		osExit(2)
	}
}

func RegisterCmd(cmdName string, help string, cmd Cmd) {
	err := RegisterCmdErr(cmdName, help, cmd)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		osExit(2)
	}
}

var InvalidCmdNameErr = errors.New("Error: cmdName can NOT start with '-'\n")
var ReservedCmdNameErr = errors.New("Error: Not Register 'help' cmd\n")

var exitKeepAlive = make(chan struct{})

func RegisterKeepAliveCmdErr(cmdName string, help string, cmd Cmd) error {
	return RegisterCmdErr(cmdName, help, func(args *arg.Arg) {
		cmd(args)

		// keep process alive
		<-exitKeepAlive
	})
}

func RegisterCmdErr(cmdName string, help string, cmd Cmd) error {
	// replace DefaultCmd directly
	if cmdName == DefaultCmdName {
		cmds[DefaultCmdName] = cmd
		helps[DefaultCmdName] = help
		return nil
	}

	if _, ok := helpCmds[cmdName]; ok {
		return ReservedCmdNameErr
	}

	if cmdName[0] == '-' {
		return InvalidCmdNameErr
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

func RegisterKeepAliveCmdNoArgs(cmdName string, help string, cmd func()) {
	RegisterKeepAliveCmd(cmdName, help, func(args *arg.Arg) {
		args.ParseAndRunHook()
		cmd()
	})
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

type HelpErr string

func (e *HelpErr) Error() string {
	return string(*e)
}

func help() *HelpErr {
	builder := &strings.Builder{}
	usage(builder)
	ret := HelpErr(builder.String())
	return &ret
}

func Run() {
	err := RunErr()
	if err == nil {
		return
	}
	_, _ = fmt.Fprint(os.Stderr, err.Error())
	if _, ok := err.(*HelpErr); ok {
		osExit(0)
	}
	osExit(2)
}

func RunErr() error {

	args := os.Args[1:]
	if len(args) == 0 {
		if cmd, ok := cmds[DefaultCmdName]; ok {
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
		if cmd, ok := cmds[DefaultCmdName]; ok {
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
