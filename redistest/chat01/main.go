package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {

	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("listen to ", l.Addr().String())

	for {
		c, err := l.Accept()
		if err != nil {
			log.Println("AcceptError: ", err)
			continue
		}
		go handle(c)
	}

}

func handle(c net.Conn) {
	cc := NewCollect(c)
	if err := cc.parse(); err != nil {
		log.Fatal(err)
	}
}

type collect struct {
	r  io.Reader
	ch chan Reply
}

func NewCollect(reader io.Reader) *collect {
	c := &collect{
		r:  reader,
		ch: make(chan Reply, 1),
	}
	go c.run()
	return c
}

func (c *collect) run() {
	for reply := range c.ch {
		fmt.Println(reply)
	}
}

func (c *collect) exit() {
	close(c.ch)
}

func (c *collect) submit(r Reply) {
	c.ch <- r
}

func (c *collect) parse() error {
	b := bufio.NewReader(c.r)
	defer c.exit()
	for {
		line, err := b.ReadBytes('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		fmt.Printf("line ==> %q \n ", string(line))
		if len(line) <= 2 || line[len(line)-2] != '\r' {
			log.Println("WARNING: read line is ", line)
			continue
		}
		var reply Reply
		switch line[0] {
		case '+':
			reply, err = NewStatusReply(string(line[1:]))
		case '-':
			reply, err = NewErrorReply(string(line[1:]))
		case ':':
			reply, err = NewNumberReply(string(line[1:]))
		case '$':
			reply, err = NewStringReply(string(line[1:]), b)
		case '*':
			reply, err = NewArrayReply(string(line[1:]), b)
		default:
			log.Println("WARNING: not support sept ", line[0])
		}

		if err != nil {
			return err
		}
		c.submit(reply)
	}
}
