package args

import (
  "github.com/xpwu/go-commandline/tinyFlag"
  "os"
  "path"
)

type args struct {
  ExeAbsDir  string
  ExeName    string
  ConfigFile string
}

var (
  Args     = args{}
  isInited = false
)

func init() {
  exe := os.Args[0]
  if path.IsAbs(exe) {
    Args.ExeAbsDir = path.Dir(exe)
  } else {
    pwd, err := os.Getwd()
    if err != nil {
      panic(err.Error())
    }
    Args.ExeAbsDir = path.Join(pwd, path.Dir(exe))
  }
  Args.ExeName = path.Base(exe)
}

func InitConfigFlag(defaultValue string, tips string)  {
  if isInited {
    return
  }
  isInited = true

  tinyFlag.StringVar(&Args.ConfigFile, "c", defaultValue, tips)
}

