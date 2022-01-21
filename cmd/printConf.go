package cmd

import (
  "github.com/xpwu/go-cmd/arg"
  "github.com/xpwu/go-cmd/exe"
  "github.com/xpwu/go-config/configs"
  "path/filepath"
)

func init() {
  argR := "config.json.default"
  RegisterCmd("pcjson", "print config with json", func(args *arg.Arg) {
    args.String(&argR, "c", "the file name of config file")
    args.Parse()
    if !filepath.IsAbs(argR) {
      argR = filepath.Join(exe.Exe.AbsDir, argR)
    }
    configs.SetConfigurator(&configs.JsonConfig{PrintFile: argR})
    configs.Print()
  })
}
