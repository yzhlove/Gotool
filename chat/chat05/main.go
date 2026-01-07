package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

const sockVer = 0x05
const noAuth = 0x00
const cmdConnect = 0x01
const aTypeIPv4 = 0x01
const aTypeIPv6 = 0x04
const aTypeDomain = 0x03
const respSucceeded = 0x00

var (
	errVerNotSupport  = errors.New("socks version not supported! ")
	errCmdNotSupport  = errors.New("socks cmd not supported! ")
	errAtypNotSupport = errors.New("socks atyp not supported! ")
)

func main() {

	s, err := net.Listen("tcp", ":1080")
	if err != nil {
		panic(err)
	}

	fmt.Println("start listen on :1080")

	for {
		c, err := s.Accept()
		if err != nil {
			fmt.Println("accept failed! ", err)
			continue
		}
		go Process(c)
	}

}

func Process(conn net.Conn) {
	defer conn.Close()

	if err := NoAuth(conn); err != nil {
		fmt.Println("NoAuth failed! ", err)
		return
	}

	addr, port, err := Connect(conn)
	if err != nil {
		fmt.Println("Connect failed! ", err)
		return
	}

	if err = Forward(addr, port, conn); err != nil {
		fmt.Println("Forward failed! ", err)
		return
	}
}

func NoAuth(conn net.Conn) (err error) {

	// +----+----------+----------+
	// |VER | NMETHODS | METHODS  |
	// +----+----------+----------+
	// | 1  |    1     | 1 to 255 |
	// +----+----------+----------+
	// VER: 协议版本，socks5为0x05
	// NMETHODS: 支持认证的方法数量
	// METHODS: 对应NMETHODS，NMETHODS的值为多少，METHODS就有多少个字节。RFC预定义了一些值的含义，内容如下:

	buf := make([]byte, 256)

	if _, err = io.ReadFull(conn, buf[:2]); err != nil {
		return err
	}

	ver, nMethods := int(buf[0]), int(buf[1])
	if ver != sockVer {
		return errVerNotSupport
	}

	if _, err = io.ReadFull(conn, buf[:nMethods]); err != nil {
		return err
	}

	fmt.Println("ver ==> ", ver, " nMethods ==> ", nMethods, "Methods ==> ", buf[:nMethods])

	// +----+--------+
	// |VER | METHOD |
	// +----+--------+
	// | 1  |   1    |
	// +----+--------+

	// 回复无需认证
	_, err = conn.Write([]byte{sockVer, noAuth})
	return err
}

func Connect(conn net.Conn) (addr string, port uint16, err error) {

	// 我们来回忆一下请求阶段的逻辑。浏览器会发送一个包，包里面包含如下6个字段
	// +----+-----+-------+------+----------+----------+
	// |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
	// +----+-----+-------+------+----------+----------+
	// | 1  |  1  | X'00' |  1   | Variable |    2     |
	// +----+-----+-------+------+----------+----------+
	// VER 版本号，socks5的值为0x05。version 版本号， 还是 5
	// CMD 0x01表示CONNECT请求。CMD 代表请求的类型，我们只支持 connection 请求，也就是让代理服务建立新的TCP连接。
	// RSV 保留字段，值为0x00，不理会。
	// ATYP 目标地址类型，DST.ADDR的数据对应这个字段的类型。可能是 IPV4 IPV6 或者域名。
	//   0x01表示IPv4地址，DST.ADDR为4个字节
	//   0x03表示域名，DST.ADDR是一个可变长度的域名
	// DST.ADDR 一个可变长度的值，这个地址的长度是根据 atype 的类型而不同的,port 端口号，两个字节, 我们需要逐个去读取这些字段。
	// DST.PORT 目标端口，固定2个字节

	buf := make([]byte, 4)
	if _, err = io.ReadFull(conn, buf); err != nil {
		return "", 0, err
	}

	ver, cmd, rsv, atype := buf[0], buf[1], buf[2], buf[3]
	if ver != sockVer {
		return "", 0, errVerNotSupport
	}
	_ = rsv // 保留字段

	if cmd != cmdConnect {
		return "", 0, errCmdNotSupport
	}

	// 获取地址
	switch atype {
	case aTypeIPv4:
		// ipv4为 4个字节
		if _, err = io.ReadFull(conn, buf); err != nil {
			return "", 0, err
		}
		addr = net.IPv4(buf[0], buf[1], buf[2], buf[3]).String()
	case aTypeIPv6:
		return "", 0, errAtypNotSupport
	case aTypeDomain:
		// 第一个字节表示长度，后面的字节表示 domain
		if _, err = io.ReadFull(conn, buf[:1]); err != nil {
			return "", 0, err
		}

		domain := make([]byte, int(buf[0]))
		if _, err = io.ReadFull(conn, domain); err != nil {
			return "", 0, err
		}
		addr = string(domain)
	default:
		return "", 0, errAtypNotSupport
	}

	// 获取端口
	// 端口为网络字节序，即 big-endian
	if _, err = io.ReadFull(conn, buf[:2]); err != nil {
		return "", 0, err
	}

	port = binary.BigEndian.Uint16(buf[:2])
	fmt.Println("===> addr ", addr, " port ==> ", port)

	// +----+-----+-------+------+----------+----------+
	// |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
	// +----+-----+-------+------+----------+----------+
	// | 1  |  1  | X'00' |  1   | Variable |    2     |
	// +----+-----+-------+------+----------+----------+
	// VER socks版本，这里为0x05，第一个是版本号还是 socket 5。
	// REP Relay field,内容取值如下 X’00’ succeeded，第二个，就是返回的类型，这里是成功就返回0。
	// RSV 保留字段，第三个是保留字段填 0。
	// ATYPE 地址类型，第四个 atype 地址类型填 1。
	// BND.ADDR 服务绑定的地址，第五个，第六个暂时用不到，都填成 0。
	// BND.PORT 服务绑定的端口DST.PORT
	_, err = conn.Write([]byte{sockVer, respSucceeded, rsv, aTypeIPv4, 0, 0, 0, 0, 0, 0})
	return
}

func Forward(addr string, port uint16, conn net.Conn) (err error) {
	dst, err := net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		return err
	}
	defer dst.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		io.Copy(dst, conn)
	}()

	go func() {
		defer wg.Done()
		io.Copy(conn, dst)
	}()

	wg.Wait()
	return
}
