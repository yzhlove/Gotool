package ssh

import (
	"fmt"
	"testing"
)

func Test_Console(t *testing.T) {

	err := LogString([]string{"ls", "/Users/yurisa"}, func(s string) {
		fmt.Println(s)
	})

	if err != nil {
		t.Fatal(err)
	}

}
