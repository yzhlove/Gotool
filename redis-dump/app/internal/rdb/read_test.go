package rdb

import (
	"fmt"
	"os"
	"testing"
)

func Test_RDB(t *testing.T) {

	a := "/Users/yostar/Desktop/player-11084.rdb"
	f, err := os.Open(a)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	values, err := Dump(f)
	if err != nil {
		t.Fatal(err)
	}

	for _, vv := range values {
		fmt.Println(vv.String())
	}

}
