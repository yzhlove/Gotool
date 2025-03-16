package opts

import (
	"github.com/yzhlove/Gotool/bing/module/opts/app"
	"github.com/yzhlove/Gotool/bing/module/opts/cmd"
)

type OptionFunc func(options *Options)

type Options struct {
	App *app.AppInfo
	Cmd *cmd.Values
}

func New(functions ...OptionFunc) *Options {
	ret := &Options{}
	for _, fn := range functions {
		fn(ret)
	}
	return ret
}

func WithApp(app *app.AppInfo) OptionFunc {
	return func(options *Options) {
		options.App = app
	}
}

func WithCmd(values *cmd.Values) OptionFunc {
	return func(options *Options) {
		options.Cmd = values
	}
}
