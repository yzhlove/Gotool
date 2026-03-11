package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var CRLF = "\r\n"

type Reply interface {
	ToBytes() []byte
}

type StatusReply struct {
	data string
}

func NewStatusReply(status string) (*StatusReply, error) {
	return &StatusReply{status}, nil
}

func (s *StatusReply) ToBytes() []byte {
	return []byte(s.data)
}

type NumberReply struct {
	code int64
}

func NewNumberReply(code string) (*NumberReply, error) {
	value, err := strconv.ParseInt(code, 10, 64)
	if err != nil {
		return nil, err
	}
	return &NumberReply{code: value}, nil
}

func (n *NumberReply) ToBytes() []byte {
	return []byte(strconv.FormatInt(n.code, 10))
}

type ErrorReply struct {
	data string
}

func NewErrorReply(err string) (*ErrorReply, error) {
	return &ErrorReply{err}, nil
}

func (e *ErrorReply) ToBytes() []byte {
	return []byte("-" + e.data + "\r\n")
}

func (e *ErrorReply) Error() string {
	return e.data
}

type BulkStrReply struct {
	data string
}

func NewStringReply(length string, reader *bufio.Reader) (*BulkStrReply, error) {
	l, err := strconv.Atoi(strings.Trim(length, CRLF))
	if err != nil {
		return nil, err
	}

	values := make([]byte, l+2)
	if _, err = io.ReadFull(reader, values); err != nil {
		return nil, err
	}
	return &BulkStrReply{data: string(values[:l])}, nil
}

func (b *BulkStrReply) ToBytes() []byte {
	return nil
}

type ArrayReply struct {
	data []string
}

func NewArrayReply(length string, reader *bufio.Reader) (*ArrayReply, error) {
	n, err := strconv.Atoi(strings.Trim(length, CRLF))
	if err != nil {
		return nil, err
	}

	fmt.Println("n --> ", n)
	a := &ArrayReply{}
	for range n {
		fmt.Println("1111 ")
		value, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("2222 ", err)
			return nil, err
		}

		fmt.Printf("vvvvake ==> %q \n ", string(value))

		l, err := strconv.Atoi(strings.Trim(string(value[1:]), CRLF))
		if err != nil {
			return nil, err
		}

		values := make([]byte, l+2)
		if _, err = io.ReadFull(reader, values); err != nil {
			return nil, err
		}
		a.data = append(a.data, string(values[:l]))
	}
	return a, nil
}

func (a *ArrayReply) ToBytes() []byte {
	return nil
}
