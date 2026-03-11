package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
)

var CRLF = []byte("\r\n")

type Reply interface {
	ToBytes() []byte
}

type StatusReply struct {
	data []byte
}

func NewStatus(data []byte) Reply {
	return &StatusReply{data: bytes.TrimSuffix(data, CRLF)}
}

func NewOkReply() *StatusReply {
	return &StatusReply{data: []byte("OK")}
}

func (s *StatusReply) ToBytes() []byte {
	var buf bytes.Buffer
	buf.WriteByte('+')
	buf.Write(s.data)
	buf.Write(CRLF)
	return buf.Bytes()
}

type ErrReply struct {
	err []byte
}

func NewErrReply(data []byte) Reply {
	return &ErrReply{err: bytes.TrimSuffix(data, CRLF)}
}

func (e *ErrReply) ToBytes() []byte {
	var buf bytes.Buffer
	buf.WriteByte('-')
	buf.Write(e.err)
	buf.Write(CRLF)
	return buf.Bytes()
}

func (e *ErrReply) Error() string {
	return string(e.err)
}

func NewAuthErrReply() Reply {
	return NewErrReply([]byte("Auth: You must select the server to connect to!"))
}

func NewParseErrReply() Reply {
	return NewErrReply([]byte("Parse: ERR parse error!"))
}

func NewUnknownCMDReply() Reply {
	return NewErrReply([]byte("ERR unknown command!"))
}

type IntReply struct {
	code int64
}

func NewIntReply(data []byte) (Reply, error) {
	v, err := strconv.Atoi(string(bytes.TrimSuffix(data, CRLF)))
	if err != nil {
		return nil, err
	}
	return &IntReply{code: int64(v)}, nil
}

func (i *IntReply) ToBytes() []byte {
	var buf bytes.Buffer
	buf.WriteByte(':')
	buf.WriteString(strconv.Itoa(int(i.code)))
	buf.Write(CRLF)
	return buf.Bytes()
}

type BulkReply struct {
	data []byte
}

func NewBulkReply(prefix []byte, buffer *bufio.Reader) (Reply, error) {
	l, err := strconv.Atoi(string(bytes.TrimSuffix(prefix, CRLF)))
	if err != nil {
		return nil, err
	}

	b := &BulkReply{}
	b.data = make([]byte, l+2)
	_, err = io.ReadFull(buffer, b.data)
	b.data = b.data[:l]
	return b, err
}

func (b *BulkReply) ToBytes() []byte {
	var buf bytes.Buffer
	buf.WriteByte('$')
	buf.WriteString(strconv.Itoa(len(b.data)))
	buf.Write(CRLF)
	buf.Write(b.data)
	buf.Write(CRLF)
	return buf.Bytes()
}

type ArrayReply struct {
	data [][]byte
}

func NewArrayReply(prefix []byte, buffer *bufio.Reader) (Reply, error) {
	count, err := strconv.Atoi(string(bytes.TrimSuffix(prefix, CRLF)))
	if err != nil {
		return nil, err
	}

	a := &ArrayReply{}
	for range count {
		line, err := buffer.ReadBytes('\n')
		if err != nil {
			return nil, err
		}
		if line[0] != '$' {
			return nil, fmt.Errorf("ArrayReply: invalid prefix %s", line)
		}
		reply, err := NewBulkReply(line[1:], buffer)
		if err != nil {
			return nil, err
		}
		a.data = append(a.data, reply.(*BulkReply).data)
	}
	return a, nil
}

func (a *ArrayReply) ToBytes() []byte {
	var buf bytes.Buffer
	buf.WriteByte('*')
	buf.WriteString(strconv.Itoa(len(a.data)))
	buf.Write(CRLF)
	for _, v := range a.data {
		buf.WriteByte('$')
		buf.WriteString(strconv.Itoa(len(v)))
		buf.Write(CRLF)
		buf.Write(v)
		buf.Write(CRLF)
	}
	return buf.Bytes()
}
