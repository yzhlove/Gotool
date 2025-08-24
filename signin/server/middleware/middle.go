package middleware

import (
	"net/http"

	"github.com/yzhlove/Gotool/signin/server/context"
)

type Middleware interface {
	Next()
	Abort()
}

type BuildContext func(http.ResponseWriter, *http.Request) (*context.Context, error)

type (
	MiddleFunc func(ctx *context.Context, m Middleware)
	HandleFunc func(ctx *context.Context)
)

func handleWrap(fn HandleFunc) MiddleFunc {
	return func(ctx *context.Context, m Middleware) {
		fn(ctx)
		m.Next()
	}
}

func errMiddleware(w http.ResponseWriter, r *http.Request, err error) {

}

type middleware struct {
	buildFunc    BuildContext
	handleChains []MiddleFunc
	ctx          *context.Context
	next         int
}

func New(buildFunc BuildContext) *middleware {
	return &middleware{
		buildFunc: buildFunc,
	}
}

func (m *middleware) Use(fn MiddleFunc) {
	m.handleChains = append(m.handleChains, fn)
}

func (m *middleware) Next() {
	if m.next >= 0 {
		m.next++
		if m.next < len(m.handleChains) {
			m.handleChains[m.next](m.ctx, m)
		}
	}
}

func (m *middleware) Abort() {
	m.next = -1
}

func (m *middleware) run() {
	if m.next >= 0 && m.next < len(m.handleChains) {
		m.handleChains[m.next](m.ctx, m)
	}
}

func (m *middleware) build(ctx *context.Context) *middleware {
	middle := &middleware{ctx: ctx}
	middle.handleChains = make([]MiddleFunc, len(m.handleChains))
	copy(middle.handleChains, m.handleChains)
	return middle
}

func (m *middleware) Handle(fn HandleFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, err := m.buildFunc(w, r)
		if err != nil {
			errMiddleware(w, r, err)
			return
		}
		middle := m.build(ctx)
		middle.Use(handleWrap(fn))
		middle.run()
	}
}
