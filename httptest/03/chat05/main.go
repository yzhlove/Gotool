package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// tls 通信

func main() {

	go tcpServer()
	time.Sleep(time.Second)
	connectClient()
}

var (
	rootPATH = "/Users/yostar/Develop/Go/GoPath/src/rain.com/Gotool/httptest/03/chat04"
)

func loadCACert(cert string) *x509.CertPool {

	pemCA, err := os.ReadFile(filepath.Join(rootPATH, cert))
	if err != nil {
		log.Fatal(err)
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(pemCA)
	return pool
}

func loadCert(keyPath, certPath string) *tls.Certificate {

	pemKey, err := os.ReadFile(filepath.Join(rootPATH, keyPath))
	if err != nil {
		log.Fatal(err)
	}

	pemCert, err := os.ReadFile(filepath.Join(rootPATH, certPath))
	if err != nil {
		log.Fatal(err)
	}

	certificate, err := tls.X509KeyPair(pemCert, pemKey)
	if err != nil {
		log.Fatal(err)
	}
	return &certificate
}

func tcpServer() {

	poolCA := loadCACert("rootCA.cer")
	serverCert := loadCert("server.key", "server.cer")

	cfg := &tls.Config{
		Certificates: []tls.Certificate{*serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert, // 强制客户端必须提供证书
		ClientCAs:    poolCA,
		MinVersion:   tls.VersionTLS12,
	}

	rsolveTCPAddr, err := net.ResolveTCPAddr("tcp", "localhost:1443")
	if err != nil {
		log.Fatal(err)
	}

	ls, err := tls.Listen("tcp", rsolveTCPAddr.String(), cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("listen to ", ls.Addr().String())

	for {
		cc, err := ls.Accept()
		if err != nil {
			if errors.Is(err, io.EOF) {
				continue
			}
			log.Println("accept error: ", err)
		} else {
			go func(cc net.Conn) {
				defer cc.Close()
				sc := bufio.NewScanner(cc)
				for sc.Scan() {
					text := sc.Text()
					fmt.Println("client msg: ", strings.TrimSuffix(strings.TrimSpace(text), "\n"))

					if strings.Contains(text, "Bye!") {
						fmt.Println("client is bye bye! ")
						return
					}
					fmt.Fprintf(cc, "%s\n", text)
				}
			}(cc)
		}
	}

}

func connectClient() {

	poolCA := loadCACert("rootCA.cer")
	clientCert := loadCert("client.key", "client.cer")

	tlsCfg := &tls.Config{
		Certificates: []tls.Certificate{*clientCert},
		RootCAs:      poolCA,
		MinVersion:   tls.VersionTLS12,
		ServerName:   "www.yzhdomain.com",
	}

	conn, err := tls.Dial("tcp", "localhost:1443", tlsCfg)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		s := bufio.NewScanner(conn)
		for s.Scan() {
			if err := s.Err(); err != nil {
				log.Println("scan err=> ", err)
			}
			fmt.Println("server msg:", s.Text())
		}
	}()

	fmt.Fprintf(conn, "hello world\n")
	time.Sleep(time.Second)

	fmt.Fprintf(conn, "测试\n")
	time.Sleep(time.Second)

	fmt.Fprintf(conn, "今天吃什么!\n")
	time.Sleep(time.Second)

	fmt.Fprintf(conn, "this is sunday! \n")
	time.Sleep(time.Second)

	fmt.Fprintf(conn, "what are you doing? \n")
	time.Sleep(time.Second)

	fmt.Fprintf(conn, "yes yes yes!!!\n")
	time.Sleep(time.Second)

	fmt.Fprintf(conn, "Bye!\n")
	time.Sleep(time.Second)

}
