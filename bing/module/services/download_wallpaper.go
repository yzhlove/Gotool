package services

import (
	"bytes"
	"github.com/yzhlove/Gotool/bing/internal/api"
	"github.com/yzhlove/Gotool/bing/module/log"
	"image"
	"sync"
	"time"
)

var downloader *download
var once sync.Once

type download struct {
	*time.Timer
	*time.Ticker
	content image.Image
	notify  []func(img image.Image)
}

func (d *download) send() {
	if d.content != nil {
		for _, ntf := range d.notify {
			ntf(d.content)
		}
	}
}

func (d *download) sync() {
	d.content = nil
	msg, err := api.GetWallpaper()
	if err != nil {
		log.Error("[service] sync error", log.ErrAttr(err))
	} else {
		if obj, _, err := image.Decode(bytes.NewReader(msg.Content)); err != nil {
			log.Error("[service] image decoder error", log.ErrAttr(err))
		} else {
			d.content = obj
			d.send()
		}
	}
}

func (d *download) run() {
	go func() {
		d.sync()

		for {
			select {
			case <-d.Timer.C:
				d.sync()
			case <-d.Ticker.C:
				d.sync()
			}
		}
	}()
}

func AddListen(callback ...func(img image.Image)) {
	if downloader != nil {
		downloader.notify = append(downloader.notify, callback...)
	}
}

func Force() {
	if downloader != nil {
		downloader.send()
	}
}

func New() error {
	once.Do(func() {
		downloader = &download{
			Timer:   time.NewTimer(time.Hour * 12),
			Ticker:  time.NewTicker(time.Hour),
			content: nil,
			notify:  make([]func(img image.Image), 0, 4),
		}
		downloader.run()
	})
	return nil
}
