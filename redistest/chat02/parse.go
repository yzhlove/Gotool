package main

import "C"
import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

func Handel(c net.Conn) error {

	fmt.Println("handle running...")

	defer func() {
		if err := c.Close(); err != nil {
			log.Println("CloseError: ", err)
		}
	}()

	reader := bufio.NewReader(c)
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return err
	}

	if len(line) <= 2 {
		_, _ = c.Write(NewUnknownCMDReply().ToBytes())
		return fmt.Errorf("invalid cmd %s", line)
	}

	if line[0] != '*' {
		_, _ = c.Write(NewUnknownCMDReply().ToBytes())
		return fmt.Errorf("invalid prefix %s", line)
	}

	reply, err := NewArrayReply(line[1:], reader)
	if err != nil {
		_, _ = c.Write(NewParseErrReply().ToBytes())
		return err
	}

	value := reply.(*ArrayReply)

	fmt.Println("P---------------- 1")
	for _, v := range value.data {
		fmt.Println(string(v))
	}
	fmt.Println("P---------------- 2")

	if err = Hello(value); err != nil {
		fmt.Println("hello 2 resp error --> ", err)
		_, _ = c.Write(NewErrReply([]byte(err.Error())).ToBytes())
		return err
	}

	sevName, err := Auth(reply.(*ArrayReply))
	if err != nil {
		_, _ = c.Write(NewErrReply([]byte(err.Error())).ToBytes())
		return err
	}

	server, err := GetRedisServer(sevName)
	if err != nil {
		_, _ = c.Write(NewErrReply([]byte(err.Error())).ToBytes())
		return err
	}

	// 回复认证成功
	_, _ = c.Write(NewOkReply().ToBytes())
	Tcp(c, server)
	return nil
}

func Hello(reply *ArrayReply) error {
	if strings.ToUpper(string(reply.data[0])) == "HELLO" {
		if string(reply.data[1]) == "2" {
			return fmt.Errorf("ERR unknown command `HELLO`, with args beginning with: `2`")
		}
	}
	return nil
}

func Auth(reply *ArrayReply) (string, error) {
	if len(reply.data) == 0 {
		return "", fmt.Errorf("Auth: params is 0! ")
	}

	if strings.ToUpper(string(reply.data[0])) != "AUTH" {
		return "", fmt.Errorf("Auth: cmd is not AUTH! ")
	}

	server := reply.data[len(reply.data)-1]
	return string(server), nil
}

func Tcp(cli, sev net.Conn) {
	var wg sync.WaitGroup
	wg.Go(func() {
		if _, err := io.Copy(cli, sev); err != nil {
			log.Println("io.Copy: ", err)
		}
	})
	wg.Go(func() {
		if _, err := io.Copy(sev, cli); err != nil {
			log.Println("io.Copy: ", err)
		}
	})
	wg.Wait()
}
