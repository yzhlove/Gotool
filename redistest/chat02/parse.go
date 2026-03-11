package main

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

	reader := bufio.NewReader(c)
	writer := bufio.NewWriter(c)
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return err
	}

	if len(line) <= 2 {
		_, _ = writer.Write(NewUnknownCMDReply().ToBytes())
		return fmt.Errorf("invalid cmd %s", line)
	}

	if line[0] != '*' {
		_, _ = writer.Write(NewUnknownCMDReply().ToBytes())
		return fmt.Errorf("invalid prefix %s", line)
	}

	reply, err := NewArrayReply(line[1:], reader)
	if err != nil {
		_, _ = writer.Write(NewParseErrReply().ToBytes())
		return err
	}

	sevName, err := Auth(reply.(*ArrayReply))
	if err != nil {
		_, _ = writer.Write(NewErrReply([]byte(err.Error())).ToBytes())
		return err
	}

	fmt.Println("--------- 111 ")

	server, err := GetRedisServer(sevName)
	if err != nil {
		_, _ = writer.Write(NewErrReply([]byte(err.Error())).ToBytes())
		return err
	}

	fmt.Println("--------- 222 ")

	var wg sync.WaitGroup
	wg.Go(func() {
		if _, err := io.Copy(server, reader); err != nil {
			log.Println("io.Copy: ", err)
		}
	})
	wg.Go(func() {
		if _, err := io.Copy(writer, server); err != nil {
			log.Println("io.Copy: ", err)
		}
	})

	fmt.Println("--------- 333 ")

	// 回复认证成功
	resp := NewOkReply().ToBytes()
	fmt.Println("===> resp ", string(resp))
	_, _ = writer.Write(resp)
	wg.Wait()
	return nil
}

func Auth(reply *ArrayReply) (string, error) {

	fmt.Println("---------------- 1")
	for _, v := range reply.data {
		fmt.Println(string(v))
	}
	fmt.Println("---------------- 2")

	if len(reply.data) == 0 {
		return "", fmt.Errorf("Auth: params is 0! ")
	}

	if strings.ToUpper(string(reply.data[0])) != "AUTH" {
		return "", fmt.Errorf("Auth: cmd is not AUTH! ")
	}

	if len(reply.data) != 2 {
		return "", fmt.Errorf("Auth: params is not 2! ")
	}
	return string(reply.data[1]), nil
}
