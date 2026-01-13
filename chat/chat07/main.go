package main

import (
	"context"
	"log"
	"net"
	"time"
)

func main() {
	go udpServer1(context.Background())
	time.Sleep(time.Millisecond * 100)
	udpClient1()
	time.Sleep(time.Second)
}

func udpServer1(ctx context.Context) {

	ls, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 1234,
	})

	if err != nil {
		log.Fatal(err)
	}
	defer ls.Close()

	buffer := make([]byte, 1024)

	for {
		select {
		case <-ctx.Done():
			log.Println("context done! ", ctx.Err())
			return
		default:
		}

		n, addr, err := ls.ReadFromUDP(buffer)
		if err != nil {
			log.Println("read udp data failed! ", err)
		}
		log.Println("remote addr is ==> ", addr.String())
		if _, err = ls.WriteToUDP(buffer[:n], addr); err != nil {
			log.Println("send udp data failed! ", err)
		}
	}
}

func udpClient1() {
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	dstAddr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 1234}

	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		log.Fatal(err)
	}

	if _, err = conn.Write([]byte("Hello World!")); err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("recv msg => ", string(buf[:n]))
}
