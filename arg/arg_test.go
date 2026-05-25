package arg

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	arg := NewArg("testParse", []string{"-c", "config.json"})
	var config = ""
	arg.String(&config, "c", "config file")
	arg.ParseAndRunHook()
	a := assert.New(t)
	a.Equal("config.json", config)
}

func TestParseNotSupport(t *testing.T) {
	arg := NewArg("testParse", []string{"not.used", "-c", "config.json"})
	arg.FlagSet.Init(arg.FlagSet.Name(), flag.ContinueOnError)
	var builder strings.Builder
	arg.FlagSet.SetOutput(&builder)
	var config = ""
	arg.String(&config, "c", "config file")
	err := arg.ParseAndRunHookErr()
	a := assert.New(t)
	a.EqualError(err, "Error: NOT support arg 'not.used'")
	a.Equal("Error: NOT support arg 'not.used'\nUsage of testParse:\n  -c string\n    \tconfig file\n",
		builder.String())
}

func TestParseUnused(t *testing.T) {
	arg := NewArg("testParse", []string{"-not.used", "-c", "config.json"})
	arg.FlagSet.Init(arg.FlagSet.Name(), flag.ContinueOnError)
	var builder strings.Builder
	arg.FlagSet.SetOutput(&builder)
	var config = ""
	arg.String(&config, "c", "config file")
	err := arg.ParseAndRunHookErr()
	a := assert.New(t)
	a.EqualError(err, "flag provided but not defined: -not.used")
	a.Equal("flag provided but not defined: -not.used\nUsage of testParse:\n  -c string\n    \tconfig file\n",
		builder.String())
}
