package arg

import (
  "flag"
  "time"
)

type Arg struct {
  FlagSet   *flag.FlagSet
  args      []string
  callbacks []func()
}

func NewArg(name string, args []string) *Arg {
  return &Arg{
    FlagSet:   flag.NewFlagSet(name, flag.ExitOnError),
    args:      args,
    callbacks: make([]func(), 0),
  }
}

func (a *Arg) Bool(defaultValue *bool, name string, usage string) {
  a.FlagSet.BoolVar(defaultValue, name, *defaultValue, usage)
}

func (a *Arg) Int(defaultValue *int, name string, usage string) {
  a.FlagSet.IntVar(defaultValue, name, *defaultValue, usage)
}

func (a *Arg) Int64(defaultValue *int64, name string, usage string) {
  a.FlagSet.Int64Var(defaultValue, name, *defaultValue, usage)
}

func (a *Arg) Uint(defaultValue *uint, name string, usage string) {
  a.FlagSet.UintVar(defaultValue, name, *defaultValue, usage)
}

func (a *Arg) Uint64(defaultValue *uint64, name string, usage string) {
  a.FlagSet.Uint64Var(defaultValue, name, *defaultValue, usage)
}

func (a *Arg) String(defaultValue *string, name string, usage string) {
  a.FlagSet.StringVar(defaultValue, name, *defaultValue, usage)
}

func (a *Arg) Float64(defaultValue *float64, name string, usage string) {
  a.FlagSet.Float64Var(defaultValue, name, *defaultValue, usage)
}

func (a *Arg) Duration(defaultValue *time.Duration, name string, usage string) {
  a.FlagSet.DurationVar(defaultValue, name, *defaultValue, usage)
}

// Parse 执行后，自动执行所有添加的callback
func (a *Arg) AddCallBack(f func()) {
  a.callbacks = append(a.callbacks, f)
}

func (a *Arg) Parse() {
	// ignore error, not panic
	_ = a.ParseErr()
}

func (a *Arg) ParseErr() error {
	err := a.FlagSet.Parse(a.args)
	if err != nil {
		return err
	}

	for _,f := range a.callbacks {
		f()
	}

	return nil
}
