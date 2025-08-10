package cipher

import (
	"fmt"
	"math/rand/v2"
	"testing"
)

func Test_GenBook(t *testing.T) {
	BuildBook()
}

func Test_ToString(t *testing.T) {

	var count = 100

	for range count {
		num := rand.Uint64()
		fmt.Println("number = ", num, "\t toString = ", ToString(num))
	}

}

func Test_EncodeDecode(t *testing.T) {

	var a = uint64(1234)
	fmt.Println("number => ", a)
	var str = ToString(a)
	fmt.Println("string => ", str)
	var b, _ = ToUint64(str)
	fmt.Println("number => ", b)

	var count = 100

	for range count {
		value1 := rand.Uint64()
		str1 := ToString(value1)
		value2, _ := ToUint64(str1)
		if value1 != value2 {
			t.Error("error")
		}
	}

}
