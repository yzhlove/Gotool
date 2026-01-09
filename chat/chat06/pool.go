package main

import "sync"

var bytesPool = sync.Pool{
	New: func() any {
		return make([]byte, 1024)
	},
}

func GetBytes() []byte {
	return bytesPool.Get().([]byte)
}

func PutBytes(b []byte) {
	if cap(b) > 1024 {
		return
	}
	clear(b)
	bytesPool.Put(b)
}
