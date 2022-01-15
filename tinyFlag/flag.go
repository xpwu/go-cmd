package tinyFlag

import (
  "flag"
  "fmt"
  "os"
  "reflect"
  "strconv"
  "time"
)

var flagNames = make(map[string]int)

func flagName(expect string) (actual string) {
  num, ok := flagNames[expect]
  if !ok {
    flagNames[expect] = 1
    return expect
  }

  flagNames[expect] = num + 1
  return expect + strconv.Itoa(num)
}

func BoolVar(p *bool, expectName string, defaultValue bool, usage string) (actualName string) {
  actualName = flagName(expectName)
  flag.BoolVar(p, actualName, defaultValue, usage)
  return
}

func Bool(expectName string, defaultValue bool, usage string) (p *bool, actualName string) {
  actualName = flagName(expectName)
  return flag.Bool(actualName, defaultValue, usage), actualName
}

func IntVar(p *int, expectName string, defaultValue int, usage string) (actualName string) {
  actualName = flagName(expectName)
  flag.IntVar(p, actualName, defaultValue, usage)
  return
}

func Int(expectName string, defaultValue int, usage string) (p *int, actualName string) {
  actualName = flagName(expectName)
  return flag.Int(actualName, defaultValue, usage), actualName
}

func Int64Var(p *int64, expectName string, defaultValue int64, usage string) (actualName string) {
  actualName = flagName(expectName)
  flag.Int64Var(p, actualName, defaultValue, usage)
  return
}

func Int64(expectName string, defaultValue int64, usage string) (p *int64, actualName string) {
  actualName = flagName(expectName)
  return flag.Int64(actualName, defaultValue, usage), actualName
}

func UintVar(p *uint, expectName string, defaultValue uint, usage string) (actualName string) {
  actualName = flagName(expectName)
  flag.UintVar(p, actualName, defaultValue, usage)
  return
}

func Uint(expectName string, defaultValue uint, usage string) (p *uint, actualName string) {
  actualName = flagName(expectName)
  return flag.Uint(actualName, defaultValue, usage), actualName
}

func Uint64Var(p *uint64, expectName string, defaultValue uint64, usage string) (actualName string) {
  actualName = flagName(expectName)
  flag.Uint64Var(p, actualName, defaultValue, usage)
  return
}

func Uint64(expectName string, defaultValue uint64, usage string) (p *uint64, actualName string) {
  actualName = flagName(expectName)
  return flag.Uint64(actualName, defaultValue, usage), actualName
}

func StringVar(p *string, expectName string, defaultValue string, usage string) (actualName string) {
  actualName = flagName(expectName)
  flag.StringVar(p, actualName, defaultValue, usage)
  return
}

func String(expectName string, defaultValue string, usage string) (p *string, actualName string) {
  actualName = flagName(expectName)
  return flag.String(actualName, defaultValue, usage), actualName
}

func Float64Var(p *float64, expectName string, defaultValue float64, usage string) (actualName string) {
  actualName = flagName(expectName)
  flag.Float64Var(p, actualName, defaultValue, usage)
  return
}

func Float64(expectName string, defaultValue float64, usage string) (p *float64, actualName string) {
  actualName = flagName(expectName)
  return flag.Float64(actualName, defaultValue, usage), actualName
}

func DurationVar(p *time.Duration, expectName string, defaultValue time.Duration,
  usage string) (actualName string) {

  actualName = flagName(expectName)
  flag.DurationVar(p, actualName, defaultValue, usage)
  return
}

func Duration(expectName string, defaultValue time.Duration, usage string) (p *time.Duration, actualName string) {
  actualName = flagName(expectName)
  return flag.Duration(actualName, defaultValue, usage), actualName
}

func Parse() {
  flag.Parse()
}

// arguments remaining
type ArgRemaining struct {
  flagSet *flag.FlagSet
  argName string
}

// flagSt 中可以使用 note tag对每一个域做使用说明
func NewArgRemaining(argName string, flagSt interface{}, additionalTips ...interface{}) *ArgRemaining {
  argName = fmt.Sprintf("%s -%s  -- ", os.Args[0], argName)
  flagSet := flag.NewFlagSet(argName, flag.ExitOnError)

  val := reflect.ValueOf(flagSt)
  if val.Kind() != reflect.Ptr || val.IsNil() {
    panic(&reflect.ValueError{Method: flagSet.Name(), Kind: val.Kind()})
  }
  ele := val.Elem()

  tp := ele.Type()
  if tp.Kind() != reflect.Struct {
    panic(&reflect.ValueError{Method: flagSet.Name(), Kind: tp.Kind()})
  }

  for i := 0; i < tp.NumField(); i++ {
    field := tp.Field(i)
    if field.Type.Kind() != reflect.Ptr {
      panic(fmt.Errorf("field(%s) of %s is not ptr", field.Name, flagSet.Name()))
    }

    e := ele.Field(i)
    note := field.Tag.Get("note")

    switch field.Type.Elem().Kind() {
    case reflect.Int:
      v := 0
      if !e.IsNil() {
        v = int(e.Elem().Int())
      }
      e.Set(reflect.ValueOf(flagSet.Int(field.Name, v, note)))
    case reflect.Int64:
      var v int64 = 0
      if !e.IsNil() {
        v = e.Elem().Int()
      }
      e.Set(reflect.ValueOf(flagSet.Int64(field.Name, v, note)))
    case reflect.Uint:
      var v uint = 0
      if !e.IsNil() {
        v = uint(e.Elem().Uint())
      }
      e.Set(reflect.ValueOf(flagSet.Uint(field.Name, v, note)))
    case reflect.Uint64:
      var v uint64 = 0
      if !e.IsNil() {
        v = e.Elem().Uint()
      }
      e.Set(reflect.ValueOf(flagSet.Uint64(field.Name, v, note)))
    case reflect.Bool:
      v := false
      if !e.IsNil() {
        v = e.Elem().Bool()
      }
      e.Set(reflect.ValueOf(flagSet.Bool(field.Name, v, note)))
    case reflect.String:
      v := ""
      if !e.IsNil() {
        v = e.Elem().String()
      }
      e.Set(reflect.ValueOf(flagSet.String(field.Name, v, note)))
    default:
      panic(&reflect.ValueError{Method: flagSet.Name(), Kind: field.Type.Kind()})
    }
  }

  flagSet.Usage = func() {
    _,_ = fmt.Fprintf(flagSet.Output(), "Usage of %s:\n", argName)
    flagSet.PrintDefaults()
    _,_ = fmt.Fprintln(flagSet.Output(), additionalTips...)
  }

  return &ArgRemaining{
    flagSet: flagSet,
    argName: argName,
  }
}

func (a *ArgRemaining) Tips() string {
  return fmt.Sprintf("%s -h 查看具体参数说明", a.argName)
}

func (a *ArgRemaining) Parse() {
  // Ignore errors; CommandLine is set for ExitOnError.
  _ = a.flagSet.Parse(flag.Args())
}
