package main

import (
	"math/rand/v2"
	"sync"
	"testing"
)

var numbersQueue = [64]uint32{
	1, 2, 3, 4, 5, 6, 7, 8,
	11, 22, 33, 44, 55, 66, 77, 88,
	21, 22, 23, 24, 252, 62, 72, 28,
	1, 2, 3, 4, 5, 6, 7, 8,
	1, 2, 3, 4, 5, 6, 7, 8,
	11, 22, 33, 44, 55, 66, 77, 88,
	21, 22, 23, 24, 252, 62, 72, 28,
	1, 2, 3, 4, 5, 6, 7, 8,
}

func Benchmark_A(b *testing.B) {

	b.ResetTimer()
	b.StartTimer()

	var sum uint32

	for i := 0; i < b.N; i++ {
		num := rand.IntN(60) + 4

		meta := pStruct.Get().([]uint32)

		for k := 0; k < num; k++ {
			meta = append(meta, numbersQueue[k])
		}
		for _, v := range meta {
			sum += v
		}
		pStruct.Put(meta)
	}
	b.StopTimer()
	b.ReportAllocs()

	// Benchmark_A-10    	   22372	    224496 ns/op	     912 B/op	       1 allocs/op

}

func Benchmark_B(b *testing.B) {

	b.ResetTimer()
	b.StartTimer()

	var sum uint32

	for i := 0; i < b.N; i++ {
		num := rand.IntN(60) + 4

		point := pPoint.Get().(*[]uint32)
		meta := *point

		for k := 0; k < num; k++ {
			meta = append(meta, numbersQueue[k])
		}
		for _, v := range meta {
			sum += v
		}
		pPoint.Put(&meta)
	}
	b.StopTimer()
	b.ReportAllocs()

	// Benchmark_B-10    	   21712	    113600 ns/op	     797 B/op	       1 allocs/op

}

var pStruct = sync.Pool{
	New: func() interface{} {
		data := make([]uint32, 0, 4)
		return data
	},
}

var pPoint = sync.Pool{
	New: func() interface{} {
		data := make([]uint32, 0, 4)
		return &data
	},
}
