package context

import "context"

type Context struct {
	ctx    context.Context
	cancel context.CancelFunc
}
