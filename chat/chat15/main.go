package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"

	"go.uber.org/fx"
)

func main() {

	f := fx.New(
		fx.Provide(NewEchoHandler),
		fx.Provide(NewServerMux),
		fx.Provide(NewHTTPServer),
		fx.Invoke(func(s *http.Server) {}),
	)
	f.Run()

}

type EchoHandler struct{}

func NewEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

func (h *EchoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func NewHTTPServer(lc fx.Lifecycle, mux *http.ServeMux) *http.Server {
	s := &http.Server{Addr: ":8433", Handler: mux}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", s.Addr)
			if err != nil {
				return err
			}
			fmt.Println("starting http server on listener: ", ln.Addr().String())
			go s.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return s.Shutdown(ctx)
		},
	})
	return s
}
