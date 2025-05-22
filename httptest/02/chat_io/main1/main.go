package main

import (
	"bytes"
	"fmt"
)

func main() {

	var text = []byte("hello world!")

	buf := bytes.NewReader(text)

	data := make([]byte, 5)

	n, err := buf.ReadAt(data, 0)
	fmt.Println(n, err, string(data))

	data = make([]byte, 12)
	n, err = buf.ReadAt(data, 0)
	fmt.Println(n, err, string(data))

	data = make([]byte, 18)
	n, err = buf.ReadAt(data, 0)
	fmt.Println(n, err, string(data))

	data = make([]byte, 12)
	n, err = buf.ReadAt(data, 2)
	fmt.Println(n, err, string(data))

	data = make([]byte, 12)
	n, err = buf.ReadAt(data, 12)
	fmt.Println(n, err, string(data))

	data = make([]byte, 12)
	n, err = buf.ReadAt(data, 15)
	fmt.Println(n, err, string(data))

	buf2 := bytes.NewBuffer(text)

	data = make([]byte, 18)
	n, err = buf2.Read(data)
	fmt.Println(n, err, string(data))

}
