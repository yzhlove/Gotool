package main

import (
	"bytes"
	"fmt"
	"golang.org/x/exp/mmap"
	"io"
)

func main() {

	var file = "/Users/yurisa/Develop/GoWork/src/WorkSpace/Gotool/httptest/02/chat_io/main2/main.go"

	reader, err := mmap.Open(file)
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	wrapper := &readerWrapper{reader, 0}

	var buffer = bytes.NewBuffer([]byte{})
	n, err := io.Copy(buffer, wrapper)
	fmt.Println(n, err)
	fmt.Println()
	fmt.Println(buffer.String())
}

type readerWrapper struct {
	io.ReaderAt
	offset int64
}

func (r *readerWrapper) Read(p []byte) (n int, err error) {
	n, err = r.ReaderAt.ReadAt(p, r.offset)
	r.offset += int64(n)
	return
}
