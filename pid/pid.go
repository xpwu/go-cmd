package pid

import (
	"github.com/xpwu/go-cmd/exe"
	"io/ioutil"
	"path/filepath"
	"strconv"
)

func Read() (string, error) {
	pidFile := filepath.Join(exe.Exe.AbsDir, "pid")
	pid, err := ioutil.ReadFile(pidFile)
	if err != nil {
		return "", err
	}
	return string(pid), nil
}

func Write(pid int) error {
	pidFile := filepath.Join(exe.Exe.AbsDir, "pid")
	return ioutil.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0664)
}
