package http

import "slices"

type V struct {
	Key   string
	Value string
}

type H []V

type M struct {
	Head H
	Body []byte
}

func (h *M) SetBody(data []byte) {
	h.Body = data
}

func (h *M) SetHead(key, value string) {
	idx := slices.IndexFunc(h.Head, func(v V) bool {
		return v.Key == key
	})
	if idx == -1 {
		h.Head = append(h.Head, V{Key: key, Value: value})
		return
	}
	h.Head[idx] = V{Key: key, Value: value}
}
