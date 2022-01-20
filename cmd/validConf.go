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
    args.String(&argR, "c", "config file name")
    config.SetConfigurator(&config.JsonConfig{ReadFile: filepath.Join(exe.Exe.AbsDir, argR)})
    fmt.Print(config.Valid())
  })
}
