package interactive

import (
  "bufio"
  "fmt"
  "github.com/xpwu/go-cmd/arg"
  "github.com/xpwu/go-cmd/cmd"
  "github.com/xpwu/go-cmd/exe"
  "github.com/xpwu/go-cmd/interac"
  "github.com/xpwu/go-cmd/pid"
  "os"
  "time"
)

func init() {
  cmd.RegisterCmd("client", "interactive mode", func(args *arg.Arg) {
    args.Parse()
    client()
  })
}

func client() {
  pd, err := pid.Read()
  if err != nil {
    fmt.Print(exe.Exe.Name + "服务未启动")
    os.Exit(0)
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
      fmt.Println("连接pid=" + pd + "的" + exe.Exe.Name + "服务失败，可能服务没有启动")
      os.Exit(1)
    }
    fmt.Println(out)
    fmt.Print("\n>")
  case <-time.After(3 * time.Second):
    fmt.Println("连接pid=" + pd + "的" + exe.Exe.Name + "服务超时，可能服务没有启动")
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
