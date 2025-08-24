package context

import (
	"crypto/cipher"
	"io"
	"net/http"

	"github.com/yzhlove/Gotool/signin/helper"
	ccipher "github.com/yzhlove/Gotool/signin/package/cipher"
)

type Context struct {
	req  *http.Request
	resp http.ResponseWriter
	aead cipher.AEAD
}

func New() *Context {
	return &Context{}
}

func (ctx *Context) WithHTTP(req *http.Request, resp http.ResponseWriter) {
	ctx.req = req
	ctx.resp = resp
}

func (ctx *Context) GetRequestName() string {
	return ctx.req.RequestURI
}

func (ctx *Context) GetBody() []byte {
	return helper.Try(io.ReadAll(ctx.req.Body)).Must()
}

func (ctx *Context) Bind(aead cipher.AEAD) {
	ctx.aead = aead
}

func (ctx *Context) Parse(data []byte) ([]byte, error) {
	return ccipher.GCMOpen(ctx.aead, data, nil)
}
