package clientcli

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/xpwu/go-cmd/arg"
	cmdPackage "github.com/xpwu/go-cmd/cmd"
	"github.com/xpwu/go-log/log"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type (
	AckToClient = string
	Cmd         = func(args *arg.Arg) AckToClient
)

func format(len, maxLen int) string {
	format := "  %s"
	for i := maxLen - len + 2; i > 0; i-- {
		format += " "
	}
	format += "%s\n"

	return format
}

func usage(args *arg.Arg) AckToClient {
	ret := fmt.Sprintf("Usage: <command> [arguments]\nThe valid 'commands' are:\n")

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
		ret += fmt.Sprintf(format(len(k), maxLen), k, helps[k])
	}

	ret += fmt.Sprintf("Every 'argument' starts with '-'.\n" +
		"Use \"<command> -h\" for more information about the command.\n")

	return ret
}

var (
	cmds = map[string]Cmd{
		"-h":   usage,
		"help": usage,
		"hello": func(args *arg.Arg) AckToClient {
			return "connected, server's pid = " + strconv.Itoa(os.Getpid()) + "\n" + usage(args)
		},
	}
	reservedCmd = map[string]bool{
		"-h":    true,
		"help":  true,
		"hello": true,
	}

	helps = map[string]string{
		"help": "show this help info",
	}
)

func runServer(clientChan <-chan Request) {
	for {
		select {
		case req := <-clientChan:
			f := strings.Fields(req.Content)
			if len(f) < 1 {
				req.Response <- usage(nil)
				break
			}

			cmd, ok := cmds[f[0]]
			if !ok {
				req.Response <- usage(nil)
				break
			}

			args := arg.NewArg(f[0], f[1:])
			args.FlagSet.Init(f[0], flag.PanicOnError)
			stdout := &bytes.Buffer{}
			args.FlagSet.SetOutput(stdout)

			ack := ""
			func() {
				defer func() {
					if r := recover(); r != nil {
						if err, ok := r.(error); ok {
							if err != flag.ErrHelp {
								ack = err.Error()
							}
						} else {
							ack = fmt.Sprint("Error: ", r)
						}
					}
				}()

				ack = cmd(args)
			}()
			ack += " \n" + stdout.String()
			ack = strings.Trim(ack+"\n"+stdout.String(), " \t\n") + "\n"
			req.Response <- ack
		}
	}
}

func Start() {
	err := StartErr()
	if err == nil {
		return
	}
	_, _ = fmt.Fprint(os.Stderr, err.Error())

	os.Exit(2)
}

func StartErr() error {
	ctx, logger := log.WithCtx(context.TODO())

	clientChan, err := ChanFromClient(ctx)
	if err != nil {
		logger.Error(err)
		return err
	}

	go runServer(clientChan)

	return nil
}

func RegisterCmdErr(cmdName string, help string, cmd func(args *arg.Arg) AckToClient) error {
	if _, ok := reservedCmd[cmdName]; ok {
		return cmdPackage.ReservedCmdNameErr
	}

	if cmdName[0] == '-' {
		return cmdPackage.InvalidCmdNameErr
	}

	tryName := cmdName
	for i := 2; ; i++ {
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

func RegisterCmd(cmdName string, help string, cmd func(args *arg.Arg) AckToClient) {
	err := RegisterCmdErr(cmdName, help, cmd)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		os.Exit(2)
	}
}

func RegisterCmdNoArgs(cmdName string, help string, cmd func() Response) {
	RegisterCmd(cmdName, help, func(args *arg.Arg) Response {
		args.ParseAndRunHook()
		return cmd()
	})
}

// ---- Deprecated ----

// Deprecated: Response Listener
type (
	Response = string
	Listener = func(args *arg.Arg) Response
)

// Deprecated: running mu
var (
	running = false
	mu      = sync.Mutex{}
)

// Deprecated: run
// don't lock
func run(ctx context.Context) {
	ctx, logger := log.WithCtx(ctx)

	if running {
		return
	}

	clientChan, err := ChanFromClient(ctx)
	if err != nil {
		logger.Error(err)
		running = false
		return
	}

	running = true

	go func() {
		for {
			select {
			case req := <-clientChan:
				f := strings.Fields(req.Content)
				if len(f) < 1 {
					req.Response <- usage(nil)
					break
				}
				mu.Lock()
				cmd, ok := cmds[f[0]]
				mu.Unlock()

				if !ok {
					req.Response <- usage(nil)
					break
				}
				args := arg.NewArg(f[0], f[1:])
				args.FlagSet.Init(f[0], flag.ContinueOnError)
				buf := &bytes.Buffer{}
				args.FlagSet.SetOutput(buf)

				cmdRet := cmd(args)
				fOut := buf.String()
				if len(fOut) != 0 {
					cmdRet += " \n" + fOut
				}
				req.Response <- cmdRet
			}
		}
	}()
}

// Deprecated: Listen using: RegisterCmd
// Listen 与 RegisterCmd 不能混用，可能会有未知问题
func Listen(ctx context.Context, cmdName string, help string, ln Listener) {
	mu.Lock()
	defer mu.Unlock()

	tryName := cmdName
	for i := 2; ; i++ {
		_, ok := cmds[tryName]
		if !ok {
			break
		}

		tryName = fmt.Sprintf("%s%d", cmdName, i)
	}

	cmds[tryName] = ln
	helps[tryName] = help

	run(ctx)
}

// Deprecated: ListenNoArg using: RegisterCmdNoArgs
// ListenNoArg 与 RegisterCmd 不能混用，可能会有未知问题
func ListenNoArg(ctx context.Context, cmdName string, help string, ln func() Response) {
	Listen(ctx, cmdName, help, func(args *arg.Arg) Response {
		args.ParseAndRunHook()
		return ln()
	})
}
