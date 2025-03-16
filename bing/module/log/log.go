package log

import (
	"log/slog"
	"os"
	"sync"
)

var log *slog.Logger
var once sync.Once

func New() error {
	once.Do(func() {
		handle := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: false,
			Level:     slog.LevelDebug,
		})
		log = slog.New(handle.WithAttrs([]slog.Attr{
			{Key: "app", Value: slog.StringValue("bing_wallpaper")},
		}))
	})
	return nil
}

func Debug(msg string, args ...interface{}) {
	log.Debug(msg, args...)
}

func Info(msg string, args ...interface{}) {
	log.Info(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	log.Warn(msg, args...)
}

func Error(msg string, args ...interface{}) {
	log.Error(msg, args...)
}

func ErrAttr(err error) any {
	return slog.Any("error", err)
}
