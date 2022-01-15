package command

import "github.com/xpwu/go-config/config"

type Command interface {
  Run()
}

type FuncCommand func()

func (f FuncCommand) Run() {
  f()
}

type Creator func() Command
type When func() bool

type option struct {
  autoReadConf bool
}

type Option func(*option)

func NotAutoReadConf() Option {
  return func(o *option) {
    o.autoReadConf = false
  }
}

type cmdTable struct {
  condition When
  creator   Creator
  option    *option
}

var (
  cmdTables            = make([]cmdTable, 0)
  defaultCmd   Creator = nil
  defaultCmdOp *option = &option{autoReadConf: true}
)

func RegisterCommand(when When, creator Creator, options ...Option) {
  op := &option{autoReadConf: true}
  for _, o := range options {
    o(op)
  }

  cmdTables = append(cmdTables, cmdTable{
    condition: when,
    creator:   creator,
    option:    op,
  })
}

func RegisterDefaultCommand(creator Creator, options ...Option) {
  for _, o := range options {
    o(defaultCmdOp)
  }

  defaultCmd = creator
}

type DefaultCommand struct {
  Handler func()
}

func (d *DefaultCommand) Register() {
  RegisterDefaultCommand(func() Command {
    return FuncCommand(d.Handler)
  })
}

func Run() {
  for _, cmd := range cmdTables {
    if cmd.condition() {
      if cmd.option.autoReadConf {
        config.Read()
      }
      cmd.creator().Run()
      goto over
    }
  }

  if defaultCmd != nil {
    if defaultCmdOp.autoReadConf {
      config.Read()
    }

    defaultCmd().Run()
  }

over:
}
