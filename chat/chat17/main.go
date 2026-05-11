package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func main() {

	f := fx.New(
		fx.Provide(NewEchoHandler),
		fx.Provide(NewServerMux),
		fx.Provide(NewHTTPServer),
		fx.Provide(NewLogger),
		fx.WithLogger(func(l *slog.Logger) fxevent.Logger {
			return &fxevent.SlogLogger{Logger: l}
		}),
		fx.Invoke(func(s *http.Server) {}),
	)
	f.Run()
}

func NewLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	}))
}

type EchoHandler struct {
	log *slog.Logger
}

func NewEchoHandler(l *slog.Logger) *EchoHandler {
	return &EchoHandler{
		log: l,
	}
}

func (h *EchoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.Info("request", "method", r.Method, "url", r.URL.String())
	if _, err := io.Copy(w, r.Body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error: %v", err)
	}
}

func NewServerMux(echo *EchoHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/echo", echo)
	return mux
}

func NewHTTPServer(lc fx.Lifecycle, mux *http.ServeMux, l *slog.Logger) *http.Server {
	s := &http.Server{Addr: ":8433", Handler: mux}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", s.Addr)
			if err != nil {
				return err
			}
			l.Info("starting http server", "listener", ln.Addr().String())
			go s.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			l.Info("stopping http server")
			return s.Shutdown(ctx)
		},
	})
	return s
}
