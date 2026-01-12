package main

import (
	"fmt"
	"net"
)

const (
	TcpCMD  = 0x01
	BindCMD = 0x02
	UdpCMD  = 0x03
)

var handleMap = map[byte]func(conn net.Conn, addr *Address) error{
	TcpCMD:  TcpConnect,
	UdpCMD:  UdpConnect,
	BindCMD: BindConnect,
}

func TcpConnect(conn net.Conn, addr *Address) (err error) {

	fmt.Println("tcp connect ing...")

	dial, err := net.Dial("tcp", addr.String())
	if err != nil {
		return err
	}
	defer dial.Close()

	fmt.Println("dial tcp on ", addr.String())

	resp := NewAddressResp(RepSucceeded, nil)
	if err = resp.Write(conn); err != nil {
		return err
	}

	fmt.Println("reply clinet !!!")

	Transport(conn, dial)

	fmt.Println("over transport!")
	return
}

func UdpConnect(conn net.Conn, addr *Address) (err error) {

	return nil
}

func BindConnect(conn net.Conn, addr *Address) (err error) {

	return nil
}
