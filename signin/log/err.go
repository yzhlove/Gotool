package log

import (
	"fmt"
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

func formatStacktrace(w io.Writer) error {

	u := make([]uintptr, 10)

	ss := runtime.CallersFrames(u)
	var more = true
	var tf runtime.Frame
	for more {
		tf, more = ss.Next()
		fmt.Println("---> ", tf.Function, tf.File, tf.Line)
	}

	stack := getStack()
	defer putStack(stack)

	num := runtime.Callers(1, stack)
	next := num > 0

	cursor := runtime.CallersFrames(stack[:num])
	var f runtime.Frame
	for next {
		f, next = cursor.Next()
		if err := formatFrame(w, f.Function, f.File, f.Line); err != nil {
			return err
		}
	}
	return nil
}

func formatFrame(w io.Writer, fun, file string, line int) error {
	if _, err := w.Write([]byte(fun)); err != nil {
		return err
	}
	if _, err := w.Write([]byte("\n")); err != nil {
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
