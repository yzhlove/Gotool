package main

import (
	"fmt"
	"github.com/yzhlove/Gotool/redis-cluster/app/handler"
	"github.com/yzhlove/Gotool/redis-cluster/app/module"
	"github.com/yzhlove/Gotool/redis-cluster/app/module/env"
	"github.com/yzhlove/Gotool/redis-cluster/app/module/log"
	"github.com/yzhlove/Gotool/redis-cluster/app/module/signal"
	"log/slog"
	"strings"
)

func main() {

	if err := module.Invoke(); err != nil {
		fmt.Println("[module.Invoke] error: ", err)
		return
	}

	if err := handler.Run(); err != nil {
		fmt.Println("[handler.Run] error: ", err)
		return
	}

	log.Info("redis-cluster run success! ", slog.String("monitor port", strings.Join(env.GetRedisPorts(), ",")))
	signal.Listen()
}
