package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(NewHTTPServer),
		fx.Invoke(func(s *http.Server) {}),
	).Run()
}

func NewHTTPServer(lc fx.Lifecycle) *http.Server {
	s := &http.Server{Addr: ":8433"}
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
			fmt.Println("stopping http server")
			return s.Shutdown(ctx)
		},
	})
	return s
}
