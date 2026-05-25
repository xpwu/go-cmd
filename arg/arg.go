package arg

import (
	"errors"
	"flag"
	"fmt"
	"os"
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

func (a *Arg) AddParseHook(f func()) {
	a.callbacks = append(a.callbacks, f)
}

// Deprecated: AddCallBack using: AddParseHook
func (a *Arg) AddCallBack(f func()) {
	a.AddParseHook(f)
}

func (a *Arg) ParseAndRunHook() {
	// ignore error, not panic
	_ = a.ParseAndRunHookErr()
}

// ParseAndRunHookErr All args must be parsed
func (a *Arg) ParseAndRunHookErr() error {
	err := a.FlagSet.Parse(a.args)
	if err != nil {
		return err
	}

	// All args must be parsed
	if a.FlagSet.NArg() != 0 {
		msg := fmt.Sprintf("Error: NOT support arg '%s'", a.FlagSet.Arg(0))
		_, _ = fmt.Fprintln(a.FlagSet.Output(), msg)
		a.FlagSet.Usage()

		err = errors.New(msg)
		// copy from 'flag'
		switch a.FlagSet.ErrorHandling() {
		case flag.ContinueOnError:
			return err
		case flag.ExitOnError:
			if err == flag.ErrHelp {
				os.Exit(0)
			}
			os.Exit(2)
		case flag.PanicOnError:
			panic(err)
		}
	}

	for _, f := range a.callbacks {
		f()
	}

	return nil
}

// Deprecated: Parse using: ParseAndRunHook
func (a *Arg) Parse() {
	a.ParseAndRunHook()
}

// Deprecated: ParseErr using: ParseAndRunHookErr
func (a *Arg) ParseErr() error {
	return a.ParseAndRunHookErr()
}
