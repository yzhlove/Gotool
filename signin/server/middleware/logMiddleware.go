package middleware

import (
	"fmt"
	"time"

	"github.com/yzhlove/Gotool/signin/server/context"
)

func LogMiddleware(ctx *context.Context, m Middleware) {
	start := time.Now()
	m.Next()
	fmt.Printf("RequestName:%s Duration:%dms \n", ctx.GetRequestName(), time.Since(start).Milliseconds())
}
