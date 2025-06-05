package module

import (
	"github.com/yzhlove/Gotool/flocktest/module/in"
	"github.com/yzhlove/Gotool/flocktest/module/log"
)

type invokeFn func() error

func setInvoke(functions ...invokeFn) error {
	for _, fn := range functions {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

func Invoke() error {
	return setInvoke(
		in.New,
		log.New,
	)
}
