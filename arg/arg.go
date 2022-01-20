package arg

import (
	"flag"
	"time"
)

type Arg struct {
	FlagSet *flag.FlagSet
	args    []string
}

func NewArg(name string, args []string) *Arg {
  return &Arg{
    FlagSet: flag.NewFlagSet(name, flag.ExitOnError),
    args:    args,
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

func (a *Arg) Parse() {
  _ = a.FlagSet.Parse(a.args)
}
