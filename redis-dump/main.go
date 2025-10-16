package main

import (
	"log"

	"go.uber.org/dig"
	"rain.com/Gotool/redis-dump/app/config"
	log2 "rain.com/Gotool/redis-dump/app/log"
	"rain.com/Gotool/redis-dump/app/view"
)

func main() {
	type moduleFunc func(c *config.Config) error
	type params struct {
		dig.In
		Config  *config.Config
		Modules []moduleFunc `group:"modules"`
	}

	container := dig.New()
	container.Provide(config.New)
	container.Provide(log2.New, dig.Group("modules"))

	if err := container.Invoke(func(in params) error {
		for _, module := range in.Modules {
			if err := module(in.Config); err != nil {
				return err
			}
		}
		return view.Adapter(in.Config)
	}); err != nil {
		log.Fatal(err)
	}
}
