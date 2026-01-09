package main

import (
	"fmt"
	"testing"
)

func Test_A(t *testing.T) {

	var bbs = make([]byte, 1024)
	bbs[0] = 'H'
	bbs[1] = 'e'
	bbs[2] = 'l'
	bbs[3] = 'l'
	bbs[4] = 'o'

	t.Log(string(bbs))

	fmt.Println("==> 1", len(bbs), cap(bbs))

	clear(bbs)

	t.Log(string(bbs))

	fmt.Println("==> 2", len(bbs), cap(bbs))

}
