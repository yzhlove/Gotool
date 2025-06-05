package main

import (
	"github.com/yzhlove/Gotool/flocktest/app"
	"github.com/yzhlove/Gotool/flocktest/module"
	logs "github.com/yzhlove/Gotool/flocktest/module/log"
	"log"
)

func main() {

	if err := module.Invoke(); err != nil {
		log.Fatalf("module.Invoke init error: %v ", err)
		return
	}
	if err := app.Run(); err != nil {
		logs.Error("app.Run error", logs.ErrWrap(err))
	}
	logs.Info("app success! ")
}
