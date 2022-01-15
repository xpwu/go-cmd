package pid

import (
  "github.com/xpwu/go-commandline/args"
  "io/ioutil"
  "path"
  "strconv"
)

func Read() (string, error) {
  pidFile := path.Join(args.Args.ExeAbsDir, "pid")
  pid,err := ioutil.ReadFile(pidFile)
  if err != nil {
    return "", err
  }
  return string(pid), nil
}

func Write(pid int) error {
  pidFile := path.Join(args.Args.ExeAbsDir, "pid")
  return ioutil.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0664)
}
