package clientcli

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/xpwu/go-cmd/arg"
	_ "github.com/xpwu/go-cmd/cmd/interactive"
	"github.com/xpwu/go-cmd/interac"
	"github.com/xpwu/go-log/log"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// Response end: blank line.  space tab etc are not blank line. only ^\n$ is blank line

type Response = string
type Listener = func(args *arg.Arg) Response

func format(len, maxLen int) string {
	format := "    %s"
	for i := maxLen - len + 2; i > 0; i-- {
		format += " "
	}
	format += ": %s\n"

	return format
}

func usage(args *arg.Arg) Response {
	ret := fmt.Sprintf("Usage:\n    <command> [arguments]\n \nThe commands are: \n \n")

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

	ret += fmt.Sprintf(" \n    \" <command> -h\" for more information about the command.\n")

	return ret
}

var (
	cmds    = map[string]Listener{}
	helps   = map[string]string{}
	running = false

	mu = sync.Mutex{}
)

func initV() {
	cmds = map[string]Listener{
		"-h":   usage,
		"help": usage,
		"hello": func(args *arg.Arg) Response {
			return "connected, server's pid = " + strconv.Itoa(os.Getpid()) + "\n" + usage(args)
		},
	}

	helps = map[string]string{
		"help": "show this help info",
	}

	running = false
}

func init() {
	initV()
}

// dont lock
func run(ctx context.Context) {
	ctx, logger := log.WithCtx(ctx)

	if running {
		return
	}

	clientChan, err := interac.ChanFromClient(ctx)
	if err != nil {
		logger.Error(err)
		running = false
		return
	}

	running = true

	go func() {
		for {
			select {
			case <-ctx.Done():
				mu.Lock()
				initV()
				mu.Unlock()

				return

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

func Listen(ctx context.Context, cmdName string, help string, ln Listener) {
	mu.Lock()
	defer mu.Unlock()

	tryName := cmdName
	for i := 0; ; i++ {
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

func ListenNoArg(ctx context.Context, cmdName string, help string, ln func() Response) {
	Listen(ctx, cmdName, help, func(args *arg.Arg) Response {
		args.Parse()
		return ln()
	})
}
