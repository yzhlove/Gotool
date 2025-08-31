package log

import (
	"io"
	"runtime"
	"strconv"
	"sync"
)

var stackPools = sync.Pool{
	New: func() interface{} {
		return make([]uintptr, 64)
	},
}

func getStack() []uintptr {
	return stackPools.Get().([]uintptr)
}

func putStack(stack []uintptr) {
	clear(stack)
	stack = stack[:0]
	stackPools.Put(stack)
}

func formatStacktrace(buf io.Writer, trace []uintptr) {
	stash, _ := runtime.CallersFrames(trace).Next()

	stack := getStack()
	defer putStack(stack)

	num := runtime.Callers(1, stack)
	cursor := runtime.CallersFrames(stack[:num])
	var frame runtime.Frame
	var next = num > 0
	var found bool

	for next {
		frame, next = cursor.Next()
		if found {
			_ = formatFrame(buf, frame.Function, frame.File, frame.Line)
			continue
		}
		if stash.Function == frame.Function &&
			stash.File == frame.File &&
			stash.Line == frame.Line {
			_ = formatFrame(buf, frame.Function, frame.File, frame.Line)
			found = true
		}
	}
	if !found {
		_ = formatFrame(buf, stash.Function, stash.File, stash.Line)
	}
}

func formatFrame(w io.Writer, fun, file string, line int) error {
	if _, err := w.Write([]byte(fun)); err != nil {
		return err
	}
	if _, err := w.Write([]byte("\n\t")); err != nil {
		return err
	}
	if _, err := w.Write([]byte(file)); err != nil {
		return err
	}
	if _, err := w.Write([]byte(":")); err != nil {
		return err
	}
	if _, err := w.Write([]byte(strconv.Itoa(line))); err != nil {
		return err
	}
	if _, err := w.Write([]byte("\n")); err != nil {
		return err
	}
	return nil
}
