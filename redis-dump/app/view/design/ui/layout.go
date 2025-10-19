package ui

import (
	"fmt"
	"image/color"
	"log/slog"
	"net/url"
	"os"
	"sync/atomic"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/hdt3213/rdb/parser"
	"rain.com/Gotool/redis-dump/app/config"
	"rain.com/Gotool/redis-dump/app/internal/rdb"
	"rain.com/Gotool/redis-dump/app/log"
	"rain.com/Gotool/redis-dump/app/view/design/data"
	"rain.com/Gotool/redis-dump/app/view/design/widgets"
)

type rdbLayout struct {
	window fyne.Window
	conf   *config.Config
	data   binding.List[*rdb.Meta]
	notify chan struct{}
	status atomic.Bool
}

func New(window fyne.Window, conf *config.Config) *rdbLayout {
	r := &rdbLayout{
		window: window,
		conf:   conf,
		notify: make(chan struct{}),
	}
	r.status.Store(false)
	r.data = binding.NewList[*rdb.Meta](func(meta *rdb.Meta, meta2 *rdb.Meta) bool { return meta == meta2 })
	r.monitor()
	r.registry()
	return r
}

func (r *rdbLayout) registry() {
	data.AddListen(func() {
		if r.status.CompareAndSwap(false, true) {
			select {
			case r.notify <- struct{}{}:
			default:
			}
		}
	})
}

func (r *rdbLayout) monitor() {
	go func() {
		for {
			select {
			case <-r.notify:
				r.run()
			}
		}
	}()
}

func (r *rdbLayout) run() {
	defer r.status.Store(false)
	for _, uri := range data.Get() {
		u, err := url.Parse(uri)
		if err != nil {
			log.Error("parse url failed! ", slog.String("uri", uri), log.ErrWrap(err))
			continue
		}

		if u.Scheme != "file" {
			log.Error("parse scheme failed! ", slog.String("uri", uri), log.ErrWrap(err))
			continue
		}

		if err = r.setValue(u.Path); err != nil {
			log.Error("parse scheme failed! ", slog.String("uri", uri), log.ErrWrap(err))
		}
	}
}

func (r *rdbLayout) setValue(path string) error {
	reader, err := os.Open(path)
	if err != nil {
		log.Error("open file failed! ", slog.String("uri", path), log.ErrWrap(err))
		return err
	}
	defer reader.Close()
	return rdb.Dump(reader, func(meta *rdb.Meta) {
		r.data.Append(meta)
	})
}

func (r *rdbLayout) Layout() fyne.CanvasObject {

	ls := widget.NewListWithData(r.data, func() fyne.CanvasObject {
		rectangle := widgets.NewRectangleText()
		showBtn := widget.NewButtonWithIcon(" -- 详  情 -- ", theme.GridIcon(), nil)
		showBtn.Importance = widget.HighImportance
		keyLabel := widget.NewLabel("")
		return container.NewBorder(nil, nil, rectangle, showBtn, keyLabel)
	}, func(item binding.DataItem, object fyne.CanvasObject) {
		it := item.(binding.Item[*rdb.Meta])
		if ret, err := it.Get(); err == nil {
			ctrl := object.(*fyne.Container)
			if len(ctrl.Objects) == 3 {

				redisKeyLabel := ctrl.Objects[0].(*widget.Label)
				redisTypeIcon := ctrl.Objects[1].(*widgets.RectangleWithText)
				showBtn := ctrl.Objects[2].(*widget.Button)
				showBtn.OnTapped = func() {
					NewList(r.window, ret).Show()
				}

				switch ret.RedisData.GetType() {
				case parser.StringType:
					redisTypeIcon.Update(
						widgets.WithText("--STR-"),
						widgets.WithFillColor(widgets.Green),
						widgets.WithStrokeColor(widgets.Gray),
						widgets.WithTextColor(color.White),
					)
				case parser.HashType:
					redisTypeIcon.Update(
						widgets.WithText("HASH"),
						widgets.WithFillColor(widgets.Red),
						widgets.WithStrokeColor(widgets.Gray),
						widgets.WithTextColor(color.White),
					)
				case parser.ListType:
					redisTypeIcon.Update(
						widgets.WithText("LIST"),
						widgets.WithFillColor(widgets.Orange),
						widgets.WithStrokeColor(widgets.Gray),
						widgets.WithTextColor(color.White),
					)
				case parser.SetType:
					redisTypeIcon.Update(
						widgets.WithText("--SET-"),
						widgets.WithFillColor(widgets.Yellow),
						widgets.WithStrokeColor(widgets.Gray),
						widgets.WithTextColor(color.White),
					)
				case parser.ZSetType:
					redisTypeIcon.Update(
						widgets.WithText("ZSET"),
						widgets.WithFillColor(widgets.Violet),
						widgets.WithStrokeColor(widgets.Gray),
						widgets.WithTextColor(color.White),
					)
				case parser.DBSizeType, parser.AuxType:
					redisTypeIcon.Update(
						widgets.WithText("  DEF "),
						widgets.WithFillColor(widgets.Indigo),
						widgets.WithStrokeColor(widgets.Gray),
						widgets.WithTextColor(color.White),
					)
				}

				var strText string
				if tm := ret.RedisData.GetExpiration(); tm != nil {
					strText = fmt.Sprintf("Key=%s|ExpireTime=%s",
						ret.RedisData.GetKey(),
						tm.Format(time.RFC3339))
				} else {
					strText = fmt.Sprintf("Key=%s", ret.RedisData.GetKey())
				}
				redisKeyLabel.SetText(strText)
			}
		}
	})

	line := canvas.NewLine(widgets.Violet)
	line.StrokeWidth = 2
	top := container.NewBorder(nil, line, nil, nil, nil)

	return widget.NewCard(
		"RedisDump工具",
		"Tips: 请将RDB文件拖拽进当前窗口!!!",
		container.NewBorder(top, nil, nil, nil, ls))
}
