package main

import (
	"context"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func main() {

	//fx.New(
	//	fx.Provide(fx.Annotate(NewEchoHandler, fx.As(new(Router)))),
	//	fx.Provide(NewServerMux),
	//	fx.Provide(NewHTTPServer),
	//	fx.Provide(NewLogger),
	//	fx.WithLogger(func(l *slog.Logger) fxevent.Logger {
	//		return &fxevent.SlogLogger{Logger: l}
	//	}),
	//	fx.Invoke(func(s *http.Server) {}),
	//).Run()

	fx.New(
		fx.Provide(
			fx.Annotate(NewEchoHandler, fx.As(new(Router)), fx.ResultTags(`name:"echo"`)),
			fx.Annotate(NewHelloHandler, fx.As(new(Router)), fx.ResultTags(`name:"hello"`)),
			fx.Annotate(NewServerMuxWithHelloAndEcho, fx.ParamTags(`name:"hello"`, `name:"echo"`)),
			NewHTTPServer,
			NewLogger,
		),
		fx.WithLogger(func(l *slog.Logger) fxevent.Logger {
			return &fxevent.SlogLogger{Logger: l}
		}),
		fx.Invoke(func(s *http.Server) {}),
	).Run()

}

func NewLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	}))
}

type Router interface {
	http.Handler
	Pattern() string
}

type EchoHandler struct {
	*slog.Logger
}

func (e EchoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.Info("request", "method", r.Method, "url", r.URL.String())
	if _, err := io.Copy(w, r.Body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		e.Error("error", "err", err)
	}
}

func (e EchoHandler) Pattern() string {
	return "/echo"
}

func NewEchoHandler(l *slog.Logger) *EchoHandler {
	return &EchoHandler{l}
}

type HelloHandler struct {
	*slog.Logger
}

func (h HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Info("request", "method", r.Method, "url", r.URL.String())
	w.Write([]byte("hello world"))
}

func (h HelloHandler) Pattern() string {
	return "/hello"
}

func NewHelloHandler(l *slog.Logger) *HelloHandler {
	return &HelloHandler{l}
}

func NewServerMux(router Router) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle(router.Pattern(), router)
	return mux
}

func NewServerMuxWithHelloAndEcho(router1, router2 Router) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle(router1.Pattern(), router1)
	mux.Handle(router2.Pattern(), router2)
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
