package graphics

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/yzhlove/Gotool/bing/module/opts"
	"github.com/yzhlove/Gotool/bing/view/graphics/content"
)

func MainWindow(opts *opts.Options) {

	bingApp := app.NewWithID(opts.App.Id)
	windows := bingApp.NewWindow(opts.App.Desc)

	windows.SetContent(content.New(windows))
	windows.Resize(fyne.NewSize(640, 420))
	windows.ShowAndRun()
}
