package main

import (
	"fmt"
	"net"
)

var serversMap = map[string]string{
	"localhost": "127.0.0.1:6379",
}

func GetRedisServer(name string) (net.Conn, error) {
	host, ok := serversMap[name]
	if ok {
		return net.Dial("tcp", host)
	}
	return nil, fmt.Errorf("Host: %s no registry! ", name)
}
