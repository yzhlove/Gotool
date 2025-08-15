package context

import (
	"encoding/json"
	"net/http"

	"github.com/yzhlove/Gotool/signin/helper"
)

func (ctx *Context) outputJson() {
	ctx.resp.Header().Set("Content-Type", "application/json")
}

func (ctx *Context) JSON(msg any) {
	ctx.outputJson()
	ctx.resp.WriteHeader(http.StatusOK)
	_, _ = ctx.resp.Write(helper.Try(json.Marshal(msg)).Must())
}
