package main

import (
	"fmt"
	"testing"
	"time"
)

func Test_Read(t *testing.T) {

	_config = new(Config)
	_config.Path = "/Users/yostar/Desktop/JsonLog"

	ReadDir(func(c *Content) {
		fmt.Println(c.String())
	})

}

func Test_ParseTime(t *testing.T) {

	layout := "2025-10-20T06:37:19Z"

	xx, err := time.Parse(time.RFC3339, layout)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("xx --> ", xx.String())

}
