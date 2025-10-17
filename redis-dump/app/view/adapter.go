package view

import (
	"rain.com/Gotool/redis-dump/app/config"
	"rain.com/Gotool/redis-dump/app/view/console"
	"rain.com/Gotool/redis-dump/app/view/design/ui"
)

func Adapter(c *config.Config) error {
	if c.UI {
		ui.MainWindow(c)
	} else {
		console.Console(c.Path)
	}
	return nil
}
