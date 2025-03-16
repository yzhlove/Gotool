package view

import (
	"github.com/yzhlove/Gotool/bing/module/log"
	"github.com/yzhlove/Gotool/bing/module/opts"
	"github.com/yzhlove/Gotool/bing/view/graphics"
)

type screen struct {
	*opts.Options
}

func New(options *opts.Options) *screen {
	return &screen{options}
}

func (s *screen) Run() {
	if s.Cmd.UI {
		graphics.MainWindow(s.Options)
	} else {
		log.Error("not support console !!! ")
	}
}
