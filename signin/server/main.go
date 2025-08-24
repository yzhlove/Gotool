package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yzhlove/Gotool/signin/server/context"
	"github.com/yzhlove/Gotool/signin/server/handler"
	"github.com/yzhlove/Gotool/signin/server/middleware"
)

func main() {

	ctx := context.New()

	m := middleware.New(func(writer http.ResponseWriter, request *http.Request) (*context.Context, error) {
		ctx.WithHTTP(request, writer)
		return ctx, nil
	})
	m.Use(middleware.RecoverMiddleware)
	m.Use(middleware.LogMiddleware)

	mux := http.NewServeMux()
	mux.HandleFunc("/ike", m.Handle(handler.IkeHandle))
	mux.HandleFunc("/registry", m.Handle(handler.RegHandle))
	fmt.Println("listeners to 8080...")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
