package http

type V struct {
	Key   string
	Value string
}

type H []V

type M struct {
	Head H
	Body []byte
}
