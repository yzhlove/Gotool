package data

import "sync/atomic"

var _input atomic.Pointer[Input]

type Input struct {
	uris     []string
	registry []func()
	holder   *atomic.Pointer[Input]
}

func newInput() *Input {
	return &Input{
		holder: &_input,
	}
}

func (in *Input) update(values []string) {
	in.uris = values
	if p := in.holder.Load(); p != nil {
		in.registry = p.registry
	}
}

func (in *Input) reg(fn func()) {
	in.registry = append(in.registry, fn)
}

func (in *Input) notify() {
	for _, fn := range in.registry {
		fn()
	}
}

func (in *Input) apply() {
	in.holder.Store(in)
}

func Update(values []string) {
	ls := newInput()
	ls.update(values)
	ls.apply()
	ls.notify()
}

func AddListen(fn func()) {
	ls := newInput()
	ls.reg(fn)
	ls.apply()
}

func Get() []string {
	if res := _input.Load(); res != nil {
		return res.uris
	}
	return nil
}
