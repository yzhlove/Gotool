package helper

type wrap[T any] struct {
	value T
	err   error
}

func (w wrap[T]) Must() T {
	if w.err != nil {
		panic(w.err)
	}
	return w.value
}

func Try[T any](v T, err error) wrap[T] {
	return wrap[T]{value: v, err: err}
}
