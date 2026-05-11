package main

import (
	"log/slog"
	"os"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func main() {

	fx.New(
		fx.WithLogger(func() fxevent.Logger {
			handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				AddSource:   true,
				Level:       slog.LevelDebug,
				ReplaceAttr: nil,
			})
			return &fxevent.SlogLogger{Logger: slog.New(handler)}
		}),
	).Run()

}
