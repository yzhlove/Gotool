package middleware

import (
	"fmt"

	"github.com/yzhlove/Gotool/signin/helper"
	"github.com/yzhlove/Gotool/signin/server/context"
)

func RecoverMiddleware(ctx *context.Context, m Middleware) {
	defer func() {
		if x := recover(); x != nil {
			fmt.Println(helper.Trace(fmt.Sprintf("PANIC RECOVER Middleware error:%v", x)))
		}
	}()
	m.Next()
}
