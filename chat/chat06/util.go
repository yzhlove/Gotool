package main

import (
	"io"
	"sync"
)

func Transport(src, dst io.ReadWriter) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		io.Copy(dst, src)
	}()

	go func() {
		defer wg.Done()
		io.Copy(src, dst)
	}()

	wg.Wait()
}
