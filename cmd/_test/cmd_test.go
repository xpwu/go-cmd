package _test

import (
	"github.com/stretchr/testify/assert"
	"github.com/xpwu/go-cmd/cmd"
	_ "github.com/xpwu/go-cmd/cmd/interactive"
	_ "github.com/xpwu/go-cmd/cmd/printconf"
	_ "github.com/xpwu/go-cmd/cmd/validconf"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	var builder strings.Builder
	cmd.SetUsageOutput(&builder)
	cmd.Run()
	a := assert.New(t)
	a.Equal("Error: Did NOT Register 'run'(DefaultCmdName) command\n",
		builder.String())
}
