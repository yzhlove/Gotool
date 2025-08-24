package middleware

import (
	"fmt"
	"testing"

	"github.com/yzhlove/Gotool/signin/server/context"
)

func NextA(ctx *context.Context, m Middleware) {
	fmt.Println("a 1")
	m.Next()
	fmt.Println("a 2")
}

func NextB(ctx *context.Context, m Middleware) {
	fmt.Println("b 1")
	m.Next()
	fmt.Println("b 2")
}

func NextC(ctx *context.Context, m Middleware) {
	fmt.Println("c 1")
	m.Next()
	fmt.Println("c 2")
}

func Test_AAA(t *testing.T) {

	s := &middleware{
		handleChains: nil,
	}
	s.Use(NextA)
	s.Use(NextB)
	s.Use(NextC)
	s.Use(handleWrap(func(ctx *context.Context) {
		fmt.Println("ddd")
	}))
	s.run()
}
