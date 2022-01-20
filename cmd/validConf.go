package cmd

import (
  "fmt"
  "github.com/xpwu/go-cmd/arg"
  "github.com/xpwu/go-cmd/exe"
  "github.com/xpwu/go-config/config"
  "path/filepath"
)


func init() {
  argR := "config.json"
  RegisterCmd("vcjson", "valid config with json",func(args *arg.Arg) {
    args.String(&argR, "c", "the file name of config file")
    if !filepath.IsAbs(argR) {
      argR = filepath.Join(exe.Exe.AbsDir, argR)
    }
    config.SetConfigurator(&config.JsonConfig{ReadFile: argR})
    fmt.Print(config.Valid())
  })
}
