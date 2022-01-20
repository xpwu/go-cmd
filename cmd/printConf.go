package cmd

import (
  "github.com/xpwu/go-cmd/arg"
  "github.com/xpwu/go-cmd/exe"
  "github.com/xpwu/go-config/config"
  "path/filepath"
)

func init() {
  argR := "config.json.default"
  RegisterCmd("pcjson", "print config with json", func(args *arg.Arg) {
    args.String(&argR, "c", "config file name")
    config.SetConfigurator(&config.JsonConfig{PrintFile: filepath.Join(exe.Exe.AbsDir, argR)})
    config.Print()
  })
}
