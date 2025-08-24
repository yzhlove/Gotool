package context

import (
	stdctx "context"
	"crypto/cipher"
)

type Context struct {
	stdctx.Context
	stdctx.CancelFunc
	aead  cipher.AEAD
	token string
}

func New() *Context {
	ctx, cancel := stdctx.WithCancel(stdctx.Background())
	return &Context{
		Context:    ctx,
		CancelFunc: cancel,
	}
}

func (c *Context) BindAEAD(aead cipher.AEAD) {
	c.aead = aead
}

func (c *Context) BindToken(token string) {
	c.token = token
}
