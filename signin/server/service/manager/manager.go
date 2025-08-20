package manager

import (
	stdctx "context"
	"sync"

	"github.com/yzhlove/Gotool/signin/server/context"
	"github.com/yzhlove/Gotool/signin/server/service"
)

type serviceContext struct {
	data   sync.Map
	ctx    stdctx.Context
	cancel stdctx.CancelFunc
}

func (m *serviceContext) Init() error {
	return nil
}

func (m *serviceContext) Start() error {
	return nil
}

func (m *serviceContext) Stop() error {
	if m.cancel != nil {
		m.cancel()
	}
	return nil
}

func Bind(token string, ctx *context.Context) {
	if _manager != nil {
		_manager.data.Store(token, ctx)
	}
}

func Get(token string) *context.Context {
	if _manager != nil {
		if value, ok := _manager.data.Load(token); ok {
			return value.(*context.Context)
		}
	}
	return nil
}

func New() (service.Service, error) {
	return &serviceContext{}, nil
}

var _manager *serviceContext
