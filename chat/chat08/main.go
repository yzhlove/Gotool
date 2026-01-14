package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"
)

// SOCKS5 UDP 包头部结构（不考虑分片）
// +----+------+------+----------+----------+----------+
// |RSV | FRAG | ATYP | DST.ADDR | DST.PORT |   DATA   |
// +----+------+------+----------+----------+----------+
// | 2  |  1   |  1   | Variable |    2     | Variable |
// +----+------+------+----------+----------+----------+

const (
	socks5Version   = 0x05
	cmdUDPAssociate = 0x03 // UDP 关联命令
	atypIPv4        = 0x01 // IPv4 地址类型
	atypDomain      = 0x03 // 域名地址类型
	atypIPv6        = 0x04 // IPv6 地址类型
)

// 处理TCP连接的SOCKS5 UDP ASSOCIATE请求
func handleTCPConn(tcpConn net.Conn) {
	defer tcpConn.Close()

	// 1. 读取客户端的SOCKS5认证请求（只处理无需认证的情况）
	buf := make([]byte, 256)
	n, err := tcpConn.Read(buf)
	if err != nil {
		fmt.Printf("读取认证请求失败: %v\n", err)
		return
	}

	// 验证版本和认证方法（只支持无需认证 0x00）
	if buf[0] != socks5Version {
		fmt.Println("不支持的SOCKS5版本")
		return
	}
	authMethods := buf[1:n]
	hasNoAuth := false
	for _, m := range authMethods {
		if m == 0x00 {
			hasNoAuth = true
			break
		}
	}
	if !hasNoAuth {
		fmt.Println("只支持无需认证的方式")
		tcpConn.Write([]byte{socks5Version, 0xff}) // 拒绝所有认证方法
		return
	}

	// 2. 响应无需认证
	tcpConn.Write([]byte{socks5Version, 0x00})

	// 3. 读取UDP ASSOCIATE请求
	n, err = tcpConn.Read(buf)
	if err != nil {
		fmt.Printf("读取UDP关联请求失败: %v\n", err)
		return
	}
	// 验证请求格式
	if buf[0] != socks5Version || buf[1] != cmdUDPAssociate || buf[2] != 0x00 {
		fmt.Println("无效的UDP关联请求")
		// 响应失败（REP=0x01）
		tcpConn.Write([]byte{socks5Version, 0x01, 0x00, atypIPv4, 0, 0, 0, 0, 0, 0})
		return
	}

	// 4. 启动UDP服务器，监听随机端口
	udpAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:0")
	if err != nil {
		fmt.Printf("解析UDP地址失败: %v\n", err)
		return
	}
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Printf("启动UDP服务器失败: %v\n", err)
		return
	}
	defer udpConn.Close()

	// 5. 响应UDP关联请求，告知客户端UDP监听端口
	localUDPAddr := udpConn.LocalAddr().(*net.UDPAddr)
	resp := make([]byte, 10)
	resp[0] = socks5Version                                           // 版本
	resp[1] = 0x00                                                    // REP=0x00 成功
	resp[2] = 0x00                                                    // RSV
	resp[3] = atypIPv4                                                // ATYP=IPv4
	copy(resp[4:8], localUDPAddr.IP.To4())                            // 服务器UDP监听IP
	binary.BigEndian.PutUint16(resp[8:10], uint16(localUDPAddr.Port)) // 服务器UDP监听端口
	tcpConn.Write(resp)

	fmt.Printf("UDP转发服务已启动，监听UDP端口: %d\n", localUDPAddr.Port)

	// 6. 处理UDP数据包转发
	udpBuf := make([]byte, 65535) // 足够大的缓冲区（不考虑分片）
	for {
		// 读取客户端发来的SOCKS5封装的UDP包
		n, clientUDPAddr, err := udpConn.ReadFromUDP(udpBuf)
		if err != nil {
			fmt.Printf("读取UDP数据失败: %v\n", err)
			return
		}
		if n < 8 { // 最小头部长度（RSV(2)+FRAG(1)+ATYP(1)+ADDR(最小4)+PORT(2)）
			fmt.Println("UDP数据包长度不足")
			continue
		}

		// 解析SOCKS5 UDP包头部
		rsv := binary.BigEndian.Uint16(udpBuf[0:2])
		frag := udpBuf[2]
		atyp := udpBuf[3]
		if rsv != 0 || frag != 0 { // 不考虑分片，FRAG必须为0
			fmt.Println("不支持分片或无效RSV")
			continue
		}

		// 解析目标地址和端口
		var dstAddr string
		portOffset := 0
		switch atyp {
		case atypIPv4:
			if n < 10 {
				fmt.Println("IPv4地址长度不足")
				continue
			}
			ip := net.IPv4(udpBuf[4], udpBuf[5], udpBuf[6], udpBuf[7])
			port := binary.BigEndian.Uint16(udpBuf[8:10])
			dstAddr = fmt.Sprintf("%s:%d", ip.String(), port)
			portOffset = 10
		case atypDomain:
			domainLen := int(udpBuf[4])
			if n < 5+domainLen+2 {
				fmt.Println("域名地址长度不足")
				continue
			}
			domain := string(udpBuf[5 : 5+domainLen])
			port := binary.BigEndian.Uint16(udpBuf[5+domainLen : 5+domainLen+2])
			dstAddr = fmt.Sprintf("%s:%d", domain, port)
			portOffset = 5 + domainLen + 2
		case atypIPv6:
			fmt.Println("暂不支持IPv6")
			continue
		default:
			fmt.Println("不支持的地址类型")
			continue
		}

		// 提取真实UDP数据
		data := udpBuf[portOffset:n]
		if len(data) == 0 {
			fmt.Println("无有效UDP数据")
			continue
		}

		// 解析目标UDP地址
		dstUDPAddr, err := net.ResolveUDPAddr("udp", dstAddr)
		if err != nil {
			fmt.Printf("解析目标UDP地址失败: %v\n", err)
			continue
		}

		// 转发UDP数据到目标地址
		_, err = udpConn.WriteToUDP(data, dstUDPAddr)
		if err != nil {
			fmt.Printf("转发UDP数据失败: %v\n", err)
			continue
		}

		// 读取目标地址的响应数据
		respBuf := make([]byte, 65535)
		udpConn.SetReadDeadline(time.Now().Add(5 * time.Second)) // 设置5秒超时
		respN, _, err := udpConn.ReadFromUDP(respBuf)
		udpConn.SetReadDeadline(time.Time{}) // 取消超时
		if err != nil {
			fmt.Printf("读取目标UDP响应失败: %v\n", err)
			continue
		}

		// 将响应数据封装成SOCKS5 UDP包，返回给客户端
		// 重新组装头部（目标地址改为客户端请求的地址）
		respPacket := new(bytes.Buffer)
		// RSV(2) + FRAG(1)
		respPacket.Write([]byte{0x00, 0x00, 0x00})
		// ATYP + 目标地址 + 目标端口
		respPacket.Write(udpBuf[3:portOffset])
		// 响应数据
		respPacket.Write(respBuf[:respN])

		// 发送回客户端
		_, err = udpConn.WriteToUDP(respPacket.Bytes(), clientUDPAddr)
		if err != nil {
			fmt.Printf("发送UDP响应给客户端失败: %v\n", err)
			continue
		}
	}
}

func main() {

	// 监听TCP端口（SOCKS5握手用）
	listenAddr := "127.0.0.1:1234"
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Printf("监听TCP端口失败: %v\n", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("SOCKS5 UDP转发服务器已启动，监听TCP端口: %s\n", listenAddr)

	// 处理客户端TCP连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("接受连接失败: %v\n", err)
			continue
		}
		go handleTCPConn(conn) // 并发处理每个连接
	}
}
