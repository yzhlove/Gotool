package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"net"

	"golang.org/x/sync/errgroup"
)

const MinUdpDatagramLen = 10 // RSV(2)+FRAG(1)+ATYPR(1)+ADDR(4)+PORT(2)

func Transport(src, dst io.ReadWriter) error {
	var wg errgroup.Group
	wg.Go(func() error {
		_, err := io.Copy(dst, src)
		return err
	})
	wg.Go(func() error {
		_, err := io.Copy(src, dst)
		return err
	})

	if err := wg.Wait(); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}
	return nil
}

func Tunnel(ctx context.Context, udp net.PacketConn) error {
	var wg errgroup.Group
	wg.Go(func() error {

		b := GetBytes()
		defer PutBytes(b)

		for {
			n, cliAddr, err := udp.ReadFrom(b)
			if err != nil {
				return err
			}

			datagram := new(UDPDatagram)
			if err = datagram.Read(bytes.NewReader(b[:n])); err != nil {
				return err
			}

			remoteAddr, err := net.ResolveUDPAddr("udp", datagram.Header.Addr.String())
			if err != nil {
				return err
			}

			if _, err = udp.WriteTo(datagram.Data, remoteAddr); err != nil {
				return err
			}

			n2, _, err := udp.ReadFrom(b)
			if err != nil {
				return err
			}

			datagram.Header.Rsv = 0
			datagram.Data = b[:n2]
			buf := bytes.NewBuffer(nil)
			if err = datagram.Write(buf); err != nil {
				return err
			}

			_, err = udp.WriteTo(buf.Bytes(), cliAddr)
			return err
		}
	})
	return wg.Wait()
}

func TcpWaitEOF(conn io.ReadWriter) {
	b := GetBytes()
	defer PutBytes(b)
	for {
		if _, err := conn.Read(b); err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			log.Println("wait tcp failed! ", err)
			return
		}
	}
}
