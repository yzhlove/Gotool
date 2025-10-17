package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"rain.com/Gotool/redis-dump/app/config"
	"rain.com/Gotool/redis-dump/app/view/design/data"
)

func MainWindow(conf *config.Config) {
	toolboxApp := app.New()
	windows := toolboxApp.NewWindow("Redis RDB Dump")
	layout := New(windows, conf)
	windows.SetOnDropped(inputCallback)
	windows.SetContent(layout.Layout())
	windows.Resize(fyne.NewSize(1280, 720))
	windows.ShowAndRun()
}

func inputCallback(position fyne.Position, uris []fyne.URI) {
	var strs = make([]string, len(uris))
	for k, v := range uris {
		strs[k] = v.String()
		fmt.Println("--------> ", v.String())
	}
	data.Update(strs)
}
