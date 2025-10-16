package log

import (
	"log/slog"
	"os"
	"sync"

	"rain.com/Gotool/redis-dump/app/config"
)

var mLog *slog.Logger
var once sync.Once

func ErrWrap(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.AnyValue(err),
	}
}

func New(c *config.Config) {
	once.Do(func() {
		opt := &slog.HandlerOptions{
			AddSource: false,
			Level:     slog.LevelDebug,
		}
		h := slog.NewTextHandler(os.Stdout, opt)
		mLog = slog.New(h)
	})
}

func Debug(msg string, args ...interface{}) {
	mLog.Debug(msg, args...)
}

func Info(msg string, args ...interface{}) {
	mLog.Info(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	mLog.Warn(msg, args...)
}

func Error(msg string, args ...interface{}) {
	mLog.Error(msg, args...)
}
