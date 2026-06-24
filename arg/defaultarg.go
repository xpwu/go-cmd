package arg

import (
	"fmt"
	"github.com/xpwu/go-config/configs"
	"github.com/xpwu/go-x/exe"
	"os"
	"path/filepath"
)

type readConfigOption struct {
	name string
}

type ReadConfigOption func(o *readConfigOption)

// Deprecated: Options using: ReadConfigOption
type Options = ReadConfigOption

// Deprecated: Name using: ConfigFlag
func Name(v string) ReadConfigOption {
	return func(o *readConfigOption) {
		o.name = v
	}
}

func ConfigFlag(v string) ReadConfigOption {
	return func(o *readConfigOption) {
		o.name = v
	}
}

// Deprecated: ReadConfig using: HookReadConfigTo
func ReadConfig(arg *Arg, opts ...ReadConfigOption) {
	HookReadConfigTo(arg, opts...)
}

func HookReadConfigTo(arg *Arg, opts ...ReadConfigOption) {
	opt := &readConfigOption{
		name: "c",
	}
	for _, o := range opts {
		o(opt)
	}

	config := "config.json"
	arg.String(&config, opt.name, "config file path")

	arg.AddParseHook(func() {
		if !filepath.IsAbs(config) {
			config = filepath.Join(exe.AbsDir, config)
		}

		configs.SetConfigurator(&configs.JsonConfig{ReadFile: config})
		err := configs.ReadWithErr()
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	})
}
