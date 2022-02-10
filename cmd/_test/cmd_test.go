package _test

import (
  "github.com/xpwu/go-cmd/cmd"
  _ "github.com/xpwu/go-cmd/cmd/interactive"
  _ "github.com/xpwu/go-cmd/cmd/printconf"
  _ "github.com/xpwu/go-cmd/cmd/validconf"
  "testing"
)

func TestRun(t *testing.T) {
  cmd.Run()
}
