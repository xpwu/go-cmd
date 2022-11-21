package arg

import (
  "fmt"
  "github.com/xpwu/go-cmd/exe"
  "github.com/xpwu/go-config/configs"
  "os"
  "path/filepath"
)

type option struct {
  name string
}

type Options func(o *option)

func Name(v string) Options {
  return func(o *option) {
    o.name = v
  }
}

func ReadConfig(arg *Arg, opts... Options) {
  opt := &option{
    name: "c",
  }
  for _,o := range opts {
    o(opt)
  }

  config := "config.json"
  arg.String(&config, opt.name, "config file path")

  arg.AddCallBack(func() {
    if !filepath.IsAbs(config) {
      config = filepath.Join(exe.Exe.AbsDir, config)
    }

    configs.SetConfigurator(&configs.JsonConfig{ReadFile: config})
    err := configs.ReadWithErr()
    if err != nil {
      fmt.Println(err)
      os.Exit(-1)
    }
  })
}
