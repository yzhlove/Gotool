package module

import (
	"github.com/yzhlove/Gotool/redis-cluster/app/module/env"
	"github.com/yzhlove/Gotool/redis-cluster/app/module/log"
	"github.com/yzhlove/Gotool/redis-cluster/app/module/signal"
)

type startFunc func() error

func start(functions ...startFunc) error {
	for _, fn := range functions {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

func Invoke() error {
	return start(
		log.New,
		env.New,
		signal.New,
	)
}
