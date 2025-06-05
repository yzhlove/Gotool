package in

import (
	"flag"
	"sync"
)

type V struct {
	Path     string
	Duration string
	Mode     string
}

var mV V
var once sync.Once

func New() error {
	once.Do(func() {
		mV = parse()
	})
	return nil
}

func Get() V {
	return mV
}

func parse() V {
	var tv V
	flag.StringVar(&tv.Path, "p", "", "set mutex path")
	flag.StringVar(&tv.Duration, "d", "", "set try lock time duration")
	flag.StringVar(&tv.Mode, "m", "", "set lock mode")
	flag.Parse()
	return tv
}
