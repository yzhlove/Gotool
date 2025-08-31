package log

import (
	"bytes"
	"log/slog"
	"os"
	"runtime"
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
		var buf = bytes.NewBuffer([]byte{})
		stack := make([]uintptr, 1)
		runtime.Callers(2, stack)
		formatStacktrace(buf, stack)
		_logger.Error(err.Error(), slog.String("stacktrace", buf.String()))
	}
}
