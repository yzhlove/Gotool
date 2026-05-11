package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"testing"
)

func Test_A(t *testing.T) {

	conn, err := GetRedisServer("localhost")
	if err != nil {
		t.Error(err)
	}

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
		go handler(c, conn)
	}

}

func handler(c net.Conn, rs net.Conn) {

	fmt.Println("dispatch -------------11111")

	errCh := make(chan error, 2)

	pt := NewParseTest(c)
	go func() {
		errCh <- pt.parse(func(reply *ArrayReply) error {
			send := reply.ToBytes()
			_, err := rs.Write(send)
			return err
		})
	}()

	go func() {
		w := io.MultiWriter(c, os.Stdout)
		_, err := io.Copy(w, rs)
		errCh <- err
	}()

	err := <-errCh
	fmt.Println("err --> ", err)
}

type parseTest struct {
	reader *bufio.Reader
	conn   net.Conn
}

func NewParseTest(conn net.Conn) *parseTest {
	return &parseTest{
		reader: bufio.NewReader(conn),
		conn:   conn,
	}
}

func (p *parseTest) parse(callback func(reply *ArrayReply) error) error {

	for {
		line, err := p.reader.ReadBytes('\n')
		if err != nil {
			return err
		}

		if len(line) <= 2 {
			return fmt.Errorf("invalid cmd %s", line)
		}

		if line[0] != '*' {
			return fmt.Errorf("invalid prefix %s", line)
		}

		reply, err := NewArrayReply(line[1:], p.reader)
		if err != nil {
			return err
		}

		fmt.Println("------------------ 1")
		r := reply.(*ArrayReply)
		for _, d := range r.data {
			fmt.Println(string(d))
		}
		fmt.Println("------------------ 2")

		if err = callback(r); err != nil {
			return err
		}
	}
}
