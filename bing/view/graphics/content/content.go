package content

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/yzhlove/Gotool/bing/module/services"
	"image"
)

func New(w fyne.Window) fyne.CanvasObject {

	wallpaper := canvas.NewImageFromImage(nil)
	wallpaper.FillMode = canvas.ImageFillContain

	services.AddListen(func(img image.Image) {
		wallpaper.Image = img
		wallpaper.Refresh()
	})
	defer services.Force()

	button := widget.NewButtonWithIcon("SAVE", theme.InfoIcon(), func() {
		dialog.ShowInformation("Wallpaper Save", "not found save wallpaper!!! ", w)
		return
	})
	button.Importance = widget.HighImportance

	return container.NewBorder(nil, button, nil, nil, wallpaper)
}
