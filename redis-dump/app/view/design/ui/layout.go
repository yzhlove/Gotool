package ui

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"sync/atomic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/hdt3213/rdb/parser"
	"rain.com/Gotool/redis-dump/app/config"
	"rain.com/Gotool/redis-dump/app/internal/rdb"
	"rain.com/Gotool/redis-dump/app/log"
	"rain.com/Gotool/redis-dump/app/view/design/data"
	"rain.com/Gotool/redis-dump/app/view/design/res"
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
		fmt.Println("-------- 1")
		if r.status.CompareAndSwap(false, true) {
			fmt.Println("-------- 2")
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
				fmt.Println("-------- 3")
				r.run()
			}
		}
	}()
}

func (r *rdbLayout) run() {
	defer r.status.Store(false)
	for _, uri := range data.Get() {
		fmt.Println("-----------> ", uri)
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
		hashIcon := widget.NewIcon(res.ResourceHashPng)
		showBtn := widget.NewButtonWithIcon(" 详情 ", theme.GridIcon(), nil)
		showBtn.Importance = widget.HighImportance
		keyLabel := widget.NewLabel("")
		return container.NewBorder(nil, nil, hashIcon, showBtn, keyLabel)
	}, func(item binding.DataItem, object fyne.CanvasObject) {
		it := item.(binding.Item[*rdb.Meta])
		if ret, err := it.Get(); err == nil {
			ctrl := object.(*fyne.Container)
			if len(ctrl.Objects) == 3 {

				redisKeyLabel := ctrl.Objects[0].(*widget.Label)
				redisTypeIcon := ctrl.Objects[1].(*widget.Icon)
				showBtn := ctrl.Objects[2].(*widget.Button)
				showBtn.OnTapped = func() {

				}
				switch ret.RedisData.GetType() {
				case parser.StringType:
					redisTypeIcon.SetResource(res.ResourceStringPng)
				case parser.HashType:
					redisTypeIcon.SetResource(res.ResourceHashPng)
				case parser.ListType:
					redisTypeIcon.SetResource(res.ResourceListPng)
				case parser.SetType:
					redisTypeIcon.SetResource(res.ResourceSetPng)
				case parser.ZSetType:
					redisTypeIcon.SetResource(res.ResourceZsetPng)
				}
				redisKeyLabel.SetText(ret.RedisData.GetKey())
			}
		}
	})
	return ls
}
