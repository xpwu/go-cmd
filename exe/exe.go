package exe

import (
	"os"
	"path/filepath"
)

type exe struct {
	AbsDir string
	Name   string
}

var (
	Exe = &exe{}
)

func init() {
	e := os.Args[0]
	f, err := filepath.Abs(e)
	if err != nil {
		Exe.AbsDir = e
	} else {
		Exe.AbsDir = f
	}

	Exe.Name = filepath.Base(e)
}
