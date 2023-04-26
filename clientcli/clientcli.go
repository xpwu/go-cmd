package clientcli

import (
	"context"
	"flag"
	"fmt"
	"github.com/xpwu/go-cmd/arg"
	"github.com/xpwu/go-cmd/interac"
	_ "github.com/xpwu/go-cmd/interac"
	"github.com/xpwu/go-log/log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Response = string
type Listener = func(args *arg.Arg) Response

func format(len, maxLen int) string {
	format := "		%s"
	for i := maxLen - len + 2; i > 0; i-- {
		format += "	"
	}
	format += ": %s\n"

	return format
}

func usage(args *arg.Arg) Response {
	ret := fmt.Sprintf("\nUsage:\n		<command> [arguments]\n\nThe commands are: \n\n")

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

	ret += fmt.Sprintf("\n		\" <command> -h\" for more information about the command.\n\n")

	return ret
}

var cmds = map[string]Listener{
	"-h":   usage,
	"help": usage,
	"hello": func(args *arg.Arg) Response {
		return "connected, server's pid = " + strconv.Itoa(os.Getpid()) + usage(args)
	},
}

var helps = map[string]string{
	"help": "show this help info",
}

var running = false

func run(ctx context.Context) {
	ctx, logger := log.WithCtx(ctx)

	if running {
		return
	}

	clientChan, err := interac.ChanFromClient(ctx)
	if err != nil {
		logger.Error(err)
		return
	}

	running = true

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

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
				args.FlagSet.Init(f[0], flag.ContinueOnError)

				req.Response <- cmd(args)
			}
		}
	}()
}

// 不是协程安全的

func Listen(ctx context.Context, cmdName string, help string, ln Listener) {
	for i := 0; ; i++ {
		_, ok := cmds[cmdName]
		if !ok {
			break
		}

		cmdName = fmt.Sprintf("%s%d", cmdName, i)
	}

	cmds[cmdName] = ln
	helps[cmdName] = help

	run(ctx)
}

func ListenNoArg(ctx context.Context, cmdName string, help string, ln func() Response) {
	Listen(ctx, cmdName, help, func(args *arg.Arg) Response {
		args.Parse()
		return ln()
	})
}
