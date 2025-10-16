package view

import (
	"rain.com/Gotool/redis-dump/app/config"
	"rain.com/Gotool/redis-dump/app/view/console"
)

func Adapter(c *config.Config) error {
	if c.UI {

	} else {
		console.Console(c.Path)
	}
	return nil
}
