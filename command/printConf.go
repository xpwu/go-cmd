package command

import (
  "github.com/xpwu/go-commandline/tinyFlag"
  "github.com/xpwu/go-config/config"
  "os"
)

type printConfig struct {

}

//func NewPrintConfig() command.Command {
//  return printConfig{}
//}

func (p printConfig)Run()  {
  config.Print()
  os.Exit(0)
}

func init() {
  argP := false
  tinyFlag.BoolVar(&argP, "p", false, "打印配置项及默认值")
  RegisterCommand(func() bool {
    return argP
  }, func() Command {
    return printConfig{}
  }, NotAutoReadConf())
}
