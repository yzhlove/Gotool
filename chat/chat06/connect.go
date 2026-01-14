package main

import (
	"context"
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

	/*
		客户端先和 SOCKS5 服务器建立 TCP 连接，发送 UDP ASSOCIATE 请求（CMD=0x03）。
		服务器响应客户端，告知客户端用于 UDP 转发的本地端口。
		客户端后续将 UDP 数据包封装成 SOCKS5 UDP 包格式，发送到服务器指定的 UDP 端口。
		服务器解析 SOCKS5 UDP 包，转发到目标地址，再将响应包封装后返回给客户端。
	*/

	fmt.Println("udp listener connect ing...", addr.String())

	ls, err := net.ListenUDP("udp", nil)
	if err != nil {
		resp := NewAddressResp(RepFailure, nil)
		if err = resp.Write(conn); err != nil {
			return err
		}
		return err
	}
	defer ls.Close()

	fmt.Println("udp listener address...", ls.LocalAddr().String())

	udpAddr, err := NewAddrFromAddr(ls.LocalAddr(), conn.LocalAddr())
	if err != nil {
		return err
	}

	fmt.Println("udp addr...", udpAddr.String())

	resp := NewAddressResp(RepSucceeded, udpAddr)
	if err = resp.Write(conn); err != nil {
		return err
	}

	fmt.Println("udp resp ok ...")

	if err = Tunnel(context.Background(), ls); err != nil {
		return err
	}

	fmt.Println("udp tunnel ok ...")

	TcpWaitEOF(conn)
	return nil
}

func BindConnect(conn net.Conn, addr *Address) (err error) {

	// 不支持改命令
	resp := NewAddressResp(RepCmdUnsupported, nil)
	return resp.Write(conn)
}
