package cmd

import (
	"fmt"
  "github.com/xpwu/go-cmd/arg"
  "io"
  "os"
  "sort"
)

type Cmd func(args *arg.Arg)

var usageOutput io.Writer = os.Stderr

func format(len, maxLen int) string {
  format := "        %s"
  for i := maxLen - len + 2; i > 0; i-- {
    format += " "
  }
  format += ": %s\n"

  return format
}

func usage(args *arg.Arg) {
  _,_ = fmt.Fprintf(usageOutput,
    "\nUsage:\n\n        %s <command> [arguments]\n\nThe commands are: (the default command is %s) \n\n",
    os.Args[0], DefaultCmdName)

  maxLen := 0
  keys := make([]string, 0, len(helps))
  for k,_ := range helps {
    if len(k) > maxLen {
      maxLen = len(k)
    }
    keys = append(keys, k)
  }

  sort.Slice(keys, func(i, j int) bool {
    return keys[i] < keys[j]
  })

  for _,k := range keys {
    _,_ = fmt.Fprintf(usageOutput, format(len(k), maxLen), k, helps[k])
  }

  _,_ = fmt.Fprintf(usageOutput,
    "\nUse \"%s <command> -h\" for more information about the command.\n\n", os.Args[0])

}

var cmds = map[string]Cmd{
	"-h":   usage,
	"help": usage,
	DefaultCmdName: func(args *arg.Arg) {},
}

var helps = map[string]string {
  DefaultCmdName: "<not implement>",
}

const DefaultCmdName = "run"

func RegisterCmd(cmdName string, help string, cmd Cmd) {
	// run 命令直接替换
	if cmdName == DefaultCmdName {
		cmds[cmdName] = cmd
		helps[cmdName] = help
		return
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
}

func RegisterCmdNoArgs(cmdName string, help string, cmd func()) {
	RegisterCmd(cmdName, help, func(args *arg.Arg) {
		args.Parse()
		cmd()
	})
}

func Run() {
	args := os.Args[1:]
  if len(args) == 0 {
    if cmd, ok := cmds[DefaultCmdName]; ok {
      cmd(arg.NewArg(os.Args[0], args))
    }
    return
  }

	tryCmd := args[0]

  if cmd, ok := cmds[tryCmd]; ok {
    cmd(arg.NewArg(os.Args[0] + " " + tryCmd + " ", args[1:]))
    return
  }

  // default
  cmd, ok := cmds[DefaultCmdName]
  if !ok {
    return
  }
  cmd(arg.NewArg(args[0], args))
}

