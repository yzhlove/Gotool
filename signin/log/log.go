package log

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type Option interface {
	GetSrc() bool
	GetApp() string
	GetLevel() slog.Leveler
	GetType() string
}

var _logger *slog.Logger

func New(opt Option) {
	var handle slog.Handler
	options := &slog.HandlerOptions{
		AddSource:   opt.GetSrc(),
		Level:       opt.GetLevel(),
		ReplaceAttr: nil,
	}
	out := os.Stdout
	switch strings.ToUpper(opt.GetType()) {
	case "JSON":
		handle = slog.NewJSONHandler(out, options)
	case "TEXT":
		handle = slog.NewTextHandler(out, options)
	}

	_logger = slog.New(handle).With(slog.String("app", opt.GetApp()))
}

func With(args ...any) *slog.Logger {
	return _logger.With(args...)
}

func Debug(msg string, args ...any) {
	_logger.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	_logger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	_logger.Warn(msg, args...)
}

func Error(err error) {
	if err == nil {
		_logger.Error("BAD ERROR!")
	} else {
		var sb = &strings.Builder{}
		if _err := formatStacktrace(sb); _err != nil {
			_logger.Error(_err.Error())
		} else {
			_logger.Error(err.Error(), slog.String("stacktrace", sb.String()))
		}

		fmt.Println(sb.String())
	}
}
