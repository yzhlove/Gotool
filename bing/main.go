package main

import (
	"fmt"
	"github.com/yzhlove/Gotool/bing/module"
	"github.com/yzhlove/Gotool/bing/module/log"
	"github.com/yzhlove/Gotool/bing/module/opts"
	"github.com/yzhlove/Gotool/bing/module/opts/app"
	"github.com/yzhlove/Gotool/bing/module/opts/cmd"
	"github.com/yzhlove/Gotool/bing/module/services"
	"github.com/yzhlove/Gotool/bing/view"
	"go.uber.org/dig"
	"time"
)

func main() {

	if err := module.Invoke(
		log.New,
		services.New,
	); err != nil {
		fmt.Println(err.Error())
		time.Sleep(time.Second * 5)
		return
	}

	type params struct {
		dig.In
		Cmd *cmd.Values
	}

	container := dig.New()
	container.Provide(cmd.New)

	err := container.Invoke(func(ps params) {
		view.New(opts.New(
			opts.WithApp(&app.AppInfo{
				Id:   "bing-wallpaper",
				Desc: "show bing wallpaper for today! ",
			}),
			opts.WithCmd(ps.Cmd))).Run()
	})

	if err != nil {
		log.Error("run application error", log.ErrAttr(err))
		time.Sleep(time.Minute)
	}
}
