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

type cmdInfo struct {
	cmdName string
	info    string
}

func testCmdNameLess7HelpInfo(t *testing.T, helpName string, helps ...*cmdInfo) {

	a := assert.New(t)

	builder := strings.Builder{}
	last := ""

	for _, help := range helps {
		if len(help.cmdName) >= 7 {
			panic("len(cmdName) MUST be less 7 in func 'testCmdNameLess7HelpInfo' ")
		}
		if help.cmdName > "client" || help.cmdName < last {
			panic("{$last} < cmdName < 'client' Not satisfied")
		}

		builder.WriteString("\n  " + help.cmdName + strings.Repeat(" ", 8-len(help.cmdName)) + help.info)

		last = help.cmdName
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
	c := &cmdInfo{"acmd", "for testing help"}
	a := assert.New(t)

	a.Nil(cmd.RegisterCmdErr(c.cmdName, c.info, func(args *arg.Arg) {

	}))

	testCmdNameLess7HelpInfo(t, "-h", c)
}

func TestRegisterErr(t *testing.T) {
	a := assert.New(t)

	a.Equal(cmd.ReservedCmdNameErr,
		cmd.RegisterCmdErr("help", "", func(args *arg.Arg) {

		}))

	a.Equal(cmd.ReservedCmdNameErr,
		cmd.RegisterCmdErr("-h", "", func(args *arg.Arg) {

		}))

	a.Equal(cmd.ReservedCmdNameErr,
		cmd.RegisterCmdErr("--help", "", func(args *arg.Arg) {

		}))

	a.Equal(cmd.ReservedCmdNameErr,
		cmd.RegisterCmdErr("--help", "", func(args *arg.Arg) {

		}))

	a.Equal(cmd.InvalidCmdNameErr,
		cmd.RegisterCmdErr("-c", "", func(args *arg.Arg) {

		}))
}

func testHelpInfoInclude(t *testing.T, help *cmdInfo) {
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
		"  " + help.cmdName + ".*?  " + help.info + "[\\s\\S]*?" +
		"Every 'argument' starts with '-'.\\n" +
		"Use \"cmd_test <command> -h\" for more information about the command.\\n$")

	a.True(re.MatchString("\n" + readPipe(p)))
}

func TestRegisterACmd(t *testing.T) {
	a := assert.New(t)

	c := &cmdInfo{"alpha", "for testing"}
	ok := false

	err := cmd.RegisterCmdErr(c.cmdName, c.info, func(args *arg.Arg) {
		args.Bool(&ok, "ok", "every is ok")
		args.ParseAndRunHook()
	})

	a.Nil(err)

	testHelpInfoInclude(t, c)

	func() {
		oldArgs := setOsArgs("alpha", "-ok")
		defer func() { os.Args = oldArgs }()
		cmd.Run()

		a.True(ok)
	}()

	//func() {
	//	oldArgs := setOsArgs("alpha", "-h")
	//	defer func() { os.Args = oldArgs }()
	//	cmd.Run()
	//}()
}

func TestRegisterKeepAliveCmd(t *testing.T) {
	a := assert.New(t)

	c := &cmdInfo{cmd.KeepAliveCmd, "run for testing"}

	err := cmd.RegisterCmdErr(c.cmdName, c.info, func(args *arg.Arg) {

	})

	a.Nil(err)

	testHelpInfoInclude(t, c)
}
