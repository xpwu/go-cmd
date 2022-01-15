package command

import (
  "fmt"
  "github.com/xpwu/go-commandline/tinyFlag"
  "github.com/xpwu/go-config/config"
  "os"
)

type validConf struct {

}

//func NewValidConfig() Command {
//  return validConf{}
//}

func (p validConf)Run()  {
  err := config.Valid()
  if err != nil {
    fmt.Printf(err.Error())
    os.Exit(1)
  }
  fmt.Print("\nvalid ok!\n")
  os.Exit(0)
}

func init() {
  arg := false
  tinyFlag.BoolVar(&arg, "v", false, "验证配置项的完备性")
  RegisterCommand(func() bool {
    return arg
  }, func() Command {
    return validConf{}
  }, NotAutoReadConf())
}
