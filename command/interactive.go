package command

import (
  "bufio"
  "fmt"
  "github.com/xpwu/go-commandline/args"
  "github.com/xpwu/go-commandline/interac"
  "github.com/xpwu/go-commandline/pid"
  "github.com/xpwu/go-commandline/tinyFlag"
  "os"
  "time"
)

type interactiveCmd struct {
}

func init() {
  arg := false
  tinyFlag.BoolVar(&arg, "i", false, "进入交互模式")
  RegisterCommand(func() bool {
    return arg
  }, func() Command {
    return &interactiveCmd{}
  }, NotAutoReadConf())
}

func (i *interactiveCmd) Run() {

  pd, err := pid.Read()
  if err != nil {
    panic(args.Args.ExeName + "服务并没有运行")
  }

  write := interac.ChanFromServer(pd)

  response := make(chan string)
  select {
  case write <- interac.RequestChan{
    Request:  "h\n",
    Response: response,
  }:
    out, ok := <-response
    if !ok {
      fmt.Println("连接pid=" + pd + "的" + args.Args.ExeName + "服务失败，可能服务没有启动")
      os.Exit(1)
    }
    fmt.Println(out)
    fmt.Print("\n>")
  case <-time.After(3 * time.Second):
    fmt.Println("连接pid=" + pd + "的" + args.Args.ExeName + "服务超时，可能服务没有启动")
    return
  }

  input := bufio.NewScanner(os.Stdin)
  for input.Scan() {
    line := input.Text()
    if line == "q" {
      break
    }
    if line == "" {
      fmt.Print(">")
      continue
    }

    write <- interac.RequestChan{
      Request:  line + "\n",
      Response: response,
    }
    out, ok := <-response
    if !ok {
      fmt.Println("error! Maybe server stopped")
      os.Exit(1)
    }
    fmt.Println(out)
    fmt.Print(">")
  }
}
