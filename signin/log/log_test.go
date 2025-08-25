package log

import (
	"errors"
	"log/slog"
	"testing"
)

func Test_Log(t *testing.T) {

	New(defOpt{})
	Info("this is info log!", slog.String("key", "value"))
	Debug("this is info log!", slog.String("key", "value"))
	Warn("this is info log!", slog.String("key", "value"))

	Error(errors.New("this is error log! "))
}

type defOpt struct{}

func (d defOpt) GetLevel() slog.Leveler {
	return slog.LevelDebug
}

func (d defOpt) GetApp() string {
	return "test"
}

func (d defOpt) GetSrc() bool {
	return false
}

func (d defOpt) GetType() string {
	return "TEXT"
}
