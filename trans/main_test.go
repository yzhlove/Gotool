package main

import (
	"fmt"
	"testing"
)

func Test_Read(t *testing.T) {

	_config = new(Config)
	_config.Path = "/Users/yostar/Desktop/JsonLog"

	ReadDir(func(c *Content) {
		fmt.Println(c.String())
	})

}
