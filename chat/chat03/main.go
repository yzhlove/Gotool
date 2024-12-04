package main

import (
	"fmt"
	"unsafe"
)

func main() {

	var x Driver = &AMD{name: "@amd"}
	fmt.Println(x.Name(), x.CPU())

	var a []uint32

	fmt.Println(unsafe.Sizeof(a))
	fmt.Println(unsafe.Alignof(a))

	fmt.Println(unsafe.Sizeof(User{}))
	fmt.Println(unsafe.Alignof(User{}))

}

type User struct {
	A int32 // 4 -> 8
	D bool  // 1
	E struct{}
	B []int32 // 24 -> 24
	C string  // 16 -> 16
}

type Driver interface {
	Name() string
	CPU() int
}

type AMD struct {
	name string
}

func (a *AMD) Name() string {
	return a.name
}

func (a AMD) CPU() int {
	return 128
}
