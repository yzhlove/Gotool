package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"rain.com/Gotool/redis-dump/app/internal/decode"
	"rain.com/Gotool/redis-dump/app/internal/rdb"
	"rain.com/Gotool/redis-dump/app/view/design/widgets"
)

type List struct {
	value  *rdb.Meta
	data   binding.List[rdb.Data]
	view   dialog.Dialog
	window fyne.Window
}

func NewList(window fyne.Window, value *rdb.Meta) *List {
	ls := &List{
		value:  value,
		window: window,
	}
	ls.data = binding.NewList[rdb.Data](func(data rdb.Data, data2 rdb.Data) bool {
		return data == data2
	})
	for _, vdata := range value.Values {
		ls.data.Append(vdata)
	}
	ls.layout()
	return ls
}

func (l *List) Show() {
	l.view.Show()
}

func (l *List) Hide() {
	l.view.Hide()
}

func (l *List) layout() {
	ls := widget.NewListWithData(l.data, func() fyne.CanvasObject {
		dataType := widgets.NewRectangleText()
		keyLabel := widget.NewLabel("")
		valueLabel := widget.NewLabel("")

		return container.NewBorder(nil, nil,
			container.NewHBox(dataType, keyLabel), nil,
			container.NewHScroll(valueLabel))

	}, func(item binding.DataItem, object fyne.CanvasObject) {
		it := item.(binding.Item[rdb.Data])
		if ret, err := it.Get(); err == nil {
			ctrl := object.(*fyne.Container)
			if len(ctrl.Objects) == 2 {
				scorll := ctrl.Objects[0].(*container.Scroll)
				hbox := ctrl.Objects[1].(*fyne.Container)
				text := hbox.Objects[0].(*widgets.RectangleWithText)
				keyLabel := hbox.Objects[1].(*widget.Label)
				valueLabel := scorll.Content.(*widget.Label)

				valueLabel.SetText(ret.Value)
				if len(ret.Name) != 0 {
					keyLabel.SetText(ret.Name)
				}

				switch ret.Type {
				case decode.Unknown:
					text.Update(
						widgets.WithText("Unknown"),
						widgets.WithFillColor(color.White),
						widgets.WithStrokeColor(widgets.Gray),
						widgets.WithTextColor(widgets.Red),
					)
				case decode.JSON:
					text.Update(
						widgets.WithText("---JSON---"),
						widgets.WithFillColor(color.White),
						widgets.WithStrokeColor(widgets.Gray),
						widgets.WithTextColor(widgets.Red),
					)
				case decode.MsgPack:
					text.Update(
						widgets.WithText("MsgPack"),
						widgets.WithFillColor(color.White),
						widgets.WithStrokeColor(widgets.Gray),
						widgets.WithTextColor(widgets.Red),
					)
				}
			}
		}
	})

	l.view = dialog.NewCustom("显示详情!", "Close",
		widget.NewCard(
			l.value.RedisData.GetKey(), "",
			container.NewGridWrap(fyne.NewSize(1000, 480), ls)),
		l.window)
}
