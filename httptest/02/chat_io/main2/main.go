package main

import (
	"bytes"
	"fmt"
	"io"
)

type readerWrapper struct {
	io.ReaderAt
	offset int64
}

func (r *readerWrapper) Read(p []byte) (n int, err error) {
	n, err = r.ReaderAt.ReadAt(p, r.offset)
	r.offset += int64(n)
	return
}

func main() {

	var text = "hello world!!!"
	reader := bytes.NewReader([]byte(text))
	wrapper := &readerWrapper{ReaderAt: reader, offset: 0}

	buffer := bytes.NewBuffer([]byte{})
	n, err := io.Copy(buffer, wrapper)
	fmt.Println(n, err, buffer.String())

	buffer.Reset()
	wrapper = &readerWrapper{ReaderAt: reader, offset: 0}

	cache := make([]byte, 2)

	n, err = io.CopyBuffer(buffer, wrapper, cache)
	fmt.Println(n, err, buffer.String())
}
