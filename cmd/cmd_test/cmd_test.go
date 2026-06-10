package cmd_test

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/xpwu/go-cmd/arg"
	"github.com/xpwu/go-cmd/cmd"
	_ "github.com/xpwu/go-cmd/cmd/interactive"
	_ "github.com/xpwu/go-cmd/cmd/printconf"
	_ "github.com/xpwu/go-cmd/cmd/validconf"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"
)

func setOsArgs(args ...string) (oldArgs []string) {
	oldArgs = os.Args
	os.Args = append([]string{oldArgs[0]}, args...)

	return
}

func TestRunErr(t *testing.T) {
	cmd.SetExeNameInUsageForTesting("cmd_test")

	oldArgs := setOsArgs()
	defer func() { os.Args = oldArgs }()

	err := cmd.RunErr()
	a := assert.New(t)
	a.Equal(`
Error: No command specified or NOT Register 'run'(default) command
Usage: cmd_test <command> [arguments]
The valid 'commands' are: (the default command is 'run')
  client  client cli mode
  pcjson  print config with json
  vcjson  valid config with json
Every 'argument' starts with '-'.
Use "cmd_test <command> -h" for more information about the command.
`,
		"\n"+err.Error())
}

func setOsExitToPanic() {
	cmd.SetOsExitForTesting(func(code int) {
		panic(code)
	})
}

func resetOsExit() {
	cmd.SetOsExitForTesting(os.Exit)
}

func assertExitCode(a *assert.Assertions, expectCode int) {
	r := recover()
	if exitCode, ok := r.(int); ok {
		a.Equal(expectCode, exitCode)
	} else {
		panic(r)
	}
}

type pipe struct {
	r *os.File
	w *os.File
}

func redirectStdErrTo(p *pipe) (oldStdErr *os.File) {
	oldStdErr = os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w
	p.r = r
	p.w = w
	return
}

func readPipe(p *pipe) string {
	_ = p.w.Close()

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, p.r)
	return buf.String()
}

func TestRunExit2(t *testing.T) {
	cmd.SetExeNameInUsageForTesting("cmd_test")

	p := &pipe{}
	oldStdErr := redirectStdErrTo(p)
	defer func() { os.Stderr = oldStdErr }()

	oldArgs := setOsArgs()
	defer func() { os.Args = oldArgs }()

	a := assert.New(t)

	func() {
		setOsExitToPanic()
		defer resetOsExit()
		defer assertExitCode(a, 2)
		cmd.Run() // exit(2)
	}()

	a.Equal(`
Error: No command specified or NOT Register 'run'(default) command
Usage: cmd_test <command> [arguments]
The valid 'commands' are: (the default command is 'run')
  client  client cli mode
  pcjson  print config with json
  vcjson  valid config with json
Every 'argument' starts with '-'.
Use "cmd_test <command> -h" for more information about the command.
`,
		"\n"+readPipe(p))
}

type helpInfo struct {
	name string
	info string
}

func testCmdNameLess7HelpInfo(t *testing.T, helpName string, helps ...*helpInfo) {

	a := assert.New(t)

	builder := strings.Builder{}
	last := ""

	for _, help := range helps {
		if len(help.name) >= 7 {
			panic("len(cmdName) MUST be less 7 in func 'testCmdNameLess7HelpInfo' ")
		}
		if help.name > "client" || help.name < last {
			panic("{$last} < cmdName < 'client' Not satisfied")
		}

		builder.WriteString("\n  " + help.name + strings.Repeat(" ", 8-len(help.name)) + help.info)

		last = help.name
	}

	cmd.SetExeNameInUsageForTesting("cmd_test")

	p := &pipe{}
	oldStdErr := redirectStdErrTo(p)
	defer func() { os.Stderr = oldStdErr }()

	oldArgs := setOsArgs(helpName)
	defer func() { os.Args = oldArgs }()

	func() {
		setOsExitToPanic()
		defer resetOsExit()
		defer assertExitCode(a, 0)
		cmd.Run() // exit(0)
	}()

	a.Equal(`
Usage: cmd_test <command> [arguments]
The valid 'commands' are: (the default command is 'run')`+builder.String()+`
  client  client cli mode
  pcjson  print config with json
  vcjson  valid config with json
Every 'argument' starts with '-'.
Use "cmd_test <command> -h" for more information about the command.
`,
		"\n"+readPipe(p))
}

func TestRunHelpExit0(t *testing.T) {
	testCmdNameLess7HelpInfo(t, "-h")
	testCmdNameLess7HelpInfo(t, "--help")
	testCmdNameLess7HelpInfo(t, "help")
	testCmdNameLess7HelpInfo(t, "-help")
}

func TestRegisterAndHelpInfo(t *testing.T) {
	c := &helpInfo{"acmd", "for testing help"}
	a := assert.New(t)

	a.Nil(cmd.RegisterCmdErr(c.name, c.info, func(args *arg.Arg) {

	}))

	testCmdNameLess7HelpInfo(t, "-h", c)
}

var registerErrCases = []struct {
	expectErr error
	name      string
}{{cmd.ReservedCmdNameErr, "help"},
	{cmd.ReservedCmdNameErr, "-h"},
	{cmd.ReservedCmdNameErr, "--help"},
	{cmd.ReservedCmdNameErr, "--help"},
	{cmd.InvalidCmdNameErr, "-c"}}

func TestRegisterErr(t *testing.T) {
	a := assert.New(t)

	testFuncs := []func(name string) error{
		func(name string) error { return cmd.RegisterCmdErr(name, "", func(args *arg.Arg) {}) },
		func(name string) error { return cmd.RegisterKeepAliveCmdErr(name, "", func(args *arg.Arg) {}) },
	}

	for _, cs := range registerErrCases {
		for _, f := range testFuncs {
			a.Equal(cs.expectErr, f(cs.name))
		}
	}
}

func TestRegister(t *testing.T) {
	a := assert.New(t)

	testFuncs := []func(name string){func(name string) { cmd.RegisterCmd(name, "", func(args *arg.Arg) {}) },
		func(name string) { cmd.RegisterKeepAliveCmd(name, "", func(args *arg.Arg) {}) },
		func(name string) { cmd.RegisterCmdNoArgs(name, "", func() {}) },
		func(name string) { cmd.RegisterKeepAliveCmdNoArgs(name, "", func() {}) },
	}

	for _, cs := range registerErrCases {
		for _, f := range testFuncs {
			func() {
				p := &pipe{}
				oldStdErr := redirectStdErrTo(p)
				defer func() { os.Stderr = oldStdErr }()

				setOsExitToPanic()
				defer resetOsExit()
				defer assertExitCode(a, 2)

				f(cs.name)
				a.Equal(cs.expectErr.Error(), readPipe(p))
			}()
		}
	}
}

func testHelpInfoInclude(t *testing.T, help *helpInfo) {
	a := assert.New(t)

	cmd.SetExeNameInUsageForTesting("cmd_test")

	p := &pipe{}
	oldStdErr := redirectStdErrTo(p)
	defer func() { os.Stderr = oldStdErr }()

	oldArgs := setOsArgs("-h")
	defer func() { os.Args = oldArgs }()

	func() {
		setOsExitToPanic()
		defer resetOsExit()
		defer assertExitCode(a, 0)
		cmd.Run() // exit(0)
	}()

	re := regexp.MustCompile("^\\nUsage: cmd_test <command> \\[arguments]\\n" +
		"The valid 'commands' are: \\(the default command is 'run'\\)[\\s\\S]*?" +
		"  " + help.name + ".*?  " + help.info + "[\\s\\S]*?" +
		"Every 'argument' starts with '-'.\\n" +
		"Use \"cmd_test <command> -h\" for more information about the command.\\n$")

	a.True(re.MatchString("\n" + readPipe(p)))
}

func testArgInfoInclude(t *testing.T, cmdName string, help *helpInfo) {
	a := assert.New(t)

	cmd.SetExeNameInUsageForTesting("cmd_test")

	p := &pipe{}
	oldStdErr := redirectStdErrTo(p)
	defer func() { os.Stderr = oldStdErr }()

	oldArgs := setOsArgs(cmdName, "-h")
	defer func() { os.Args = oldArgs }()

	func() {
		defer func() {
			r := recover()
			if exitInfo, ok := r.(string); ok {
				a.Equal("unexpected call to os.Exit(0) during test", exitInfo)
			} else {
				panic(r)
			}
		}()

		cmd.Run() // exit(0)
	}()

	re := regexp.MustCompile("^Usage of .*? " + cmdName + " :\\n" + "[\\s\\S]*?" +
		"  -" + help.name + "[\\s\\S]*?\\t" + help.info + "[\\s\\S]*?\\n$")

	a.True(re.MatchString(readPipe(p)))
}

func TestRegisterCmd(t *testing.T) {
	a := assert.New(t)

	normal := &helpInfo{"alpha", "for testing"}
	normalArgInfo1 := &helpInfo{"ok", "every is ok"}
	normalArgInfo2 := &helpInfo{"name", "user name"}
	ok := false
	name := "go-cmd"
	err := cmd.RegisterCmdErr(normal.name, normal.info, func(args *arg.Arg) {
		args.Bool(&ok, normalArgInfo1.name, normalArgInfo1.info)
		args.String(&name, normalArgInfo2.name, normalArgInfo2.info)
		args.ParseAndRunHook()
	})
	a.Nil(err)

	keepalive := &helpInfo{"beta", "for keepalive testing"}
	keepaliveArgInfo := &helpInfo{"config", "server config"}
	config := ""
	err = cmd.RegisterKeepAliveCmdErr(keepalive.name, keepalive.info, func(args *arg.Arg) {
		args.String(&config, keepaliveArgInfo.name, keepaliveArgInfo.info)
		args.ParseAndRunHook()
	})
	a.Nil(err)

	defaultCmd := &helpInfo{cmd.DefaultCmdName, "running"}
	defaultCmdArgInfo := &helpInfo{"who", "who are you?"}
	who := ""
	err = cmd.RegisterCmdErr(defaultCmd.name, defaultCmd.info, func(args *arg.Arg) {
		args.String(&who, defaultCmdArgInfo.name, defaultCmdArgInfo.info)
		args.ParseAndRunHook()
	})
	a.Nil(err)

	testHelpInfoInclude(t, normal)
	testArgInfoInclude(t, normal.name, normalArgInfo1)
	testArgInfoInclude(t, normal.name, normalArgInfo2)
	testHelpInfoInclude(t, keepalive)
	testArgInfoInclude(t, keepalive.name, keepaliveArgInfo)
	testHelpInfoInclude(t, defaultCmd)
	testArgInfoInclude(t, defaultCmd.name, defaultCmdArgInfo)

	func() {
		oldArgs := setOsArgs(normal.name, "-ok")
		defer func() { os.Args = oldArgs }()
		cmd.Run()

		a.True(ok)
	}()
	func() {
		expect := "github"
		oldArgs := setOsArgs("run", "-who", expect)
		defer func() { os.Args = oldArgs }()
		cmd.Run()

		a.Equal(expect, who)
	}()

	normal2 := &helpInfo{normal.name + "2", "for testing2"}
	normal2ArgInfo1 := &helpInfo{"ok2", "every is ok2"}
	normal2ArgInfo2 := &helpInfo{"name2", "user name2"}
	ok = false
	name = "go-cmd2"
	// normal.name,   not normal2.name
	err = cmd.RegisterCmdErr(normal.name, normal2.info, func(args *arg.Arg) {
		args.Bool(&ok, normal2ArgInfo1.name, normal2ArgInfo1.info)
		args.String(&name, normal2ArgInfo2.name, normal2ArgInfo2.info)
		args.ParseAndRunHook()
	})
	a.Nil(err)

	defaultCmd = &helpInfo{cmd.DefaultCmdName, "running2"}
	defaultCmdArgInfo = &helpInfo{"who2", "who are you, anyway?"}
	who = ""
	err = cmd.RegisterCmdErr(defaultCmd.name, defaultCmd.info, func(args *arg.Arg) {
		args.String(&who, defaultCmdArgInfo.name, defaultCmdArgInfo.info)
		args.ParseAndRunHook()
	})
	a.Nil(err)

	testHelpInfoInclude(t, normal)
	testArgInfoInclude(t, normal.name, normalArgInfo1)
	testArgInfoInclude(t, normal.name, normalArgInfo2)
	testHelpInfoInclude(t, keepalive)
	testArgInfoInclude(t, keepalive.name, keepaliveArgInfo)
	testHelpInfoInclude(t, defaultCmd)
	testArgInfoInclude(t, defaultCmd.name, defaultCmdArgInfo)
	testHelpInfoInclude(t, normal2)
	testArgInfoInclude(t, normal2.name, normal2ArgInfo1)
	testArgInfoInclude(t, normal2.name, normal2ArgInfo2)

	func() {
		oldArgs := setOsArgs(normal.name, "-ok")
		defer func() { os.Args = oldArgs }()
		cmd.Run()

		a.True(ok)
	}()
	func() {
		expect := "github2"
		oldArgs := setOsArgs("run", "-who2", expect)
		defer func() { os.Args = oldArgs }()
		cmd.Run()

		a.Equal(expect, who)
	}()
	func() {
		expect := "go-cmd2"
		ok = false
		oldArgs := setOsArgs(normal2.name, "-"+normal2ArgInfo2.name, expect)
		defer func() { os.Args = oldArgs }()
		cmd.Run()

		a.Equal(expect, name)
		a.False(ok)
	}()
}

func TestRegisterDefaultCmd(t *testing.T) {
	a := assert.New(t)

	c := &helpInfo{cmd.DefaultCmdName, "run for testing"}

	err := cmd.RegisterCmdErr(c.name, c.info, func(args *arg.Arg) {

	})

	a.Nil(err)

	testHelpInfoInclude(t, c)
}
