package validconf

import (
  "fmt"
  "github.com/xpwu/go-cmd/arg"
  "github.com/xpwu/go-cmd/cmd"
  "github.com/xpwu/go-cmd/exe"
  "github.com/xpwu/go-config/configs"
  "path/filepath"
)


func init() {
  argR := "config.json"
  cmd.RegisterCmd("vcjson", "valid config with json",func(args *arg.Arg) {
    args.String(&argR, "c", "the file name of config file")
    if !filepath.IsAbs(argR) {
      argR = filepath.Join(exe.Exe.AbsDir, argR)
    }
    configs.SetConfigurator(&configs.JsonConfig{ReadFile: argR})
    fmt.Print(configs.Valid())
  })
}
