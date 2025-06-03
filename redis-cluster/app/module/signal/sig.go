package signal

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var once sync.Once
var sCh chan os.Signal

func New() error {
	once.Do(func() {
		sCh = make(chan os.Signal, 1)
		signal.Notify(sCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM,
			syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	})
	return nil
}

func Listen() {
	for {
		select {
		case <-sCh:
			os.Exit(0)
		}
	}
}
