package main

import (
	"fmt"
	"log"
	"net"
)

func main() {

	var addr = ":1234"

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("listen on ", addr, " ...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	remoteAddr := conn.RemoteAddr()
	log.Println("remoteAddr: ", remoteAddr.String(), " connected!")

	req := new(MethodsReq)
	if err := req.Read(conn); err != nil {
		log.Fatal(err)
	}

	log.Println("method request=> ", req.Methods)

	resp := new(MethodResp)
	resp.Method = chooseMethod(req.Methods)
	if resp.Method == NoAcceptableMethods {
		if err := resp.Write(conn); err != nil {
			log.Fatal(err)
		}
		if err := conn.Close(); err != nil {
			log.Fatal(err)
		}
		log.Println("remoteAddr closed! ", remoteAddr.String())
		return
	}

	// 认证
	if err := authMap[resp.Method](conn); err != nil {
		log.Fatal(err)
	}

	// 解析需要代理的地址
	connReq := new(ConnectReq)
	if err := connReq.Read(conn); err != nil {
		log.Fatal(err)
	}

	log.Println("connect remote addr ", connReq.Addr.String())

	fmt.Println("cmd ===> ", connReq.Cmd)

	if err := handleMap[connReq.Cmd](conn, connReq.Addr); err != nil {
		log.Fatal(err)
	}

	log.Println("connect exit!!! ", connReq.Addr.String())
}

func chooseMethod(methods []byte) byte {
	var supportNoAuth bool
	var supportUserAuth bool

	for _, m := range methods {
		if m == UserMethod {
			supportUserAuth = true
		} else if m == NoAuthMethod {
			supportNoAuth = true
		}
	}

	if supportUserAuth {
		return UserMethod
	}

	if supportNoAuth {
		return NoAuthMethod
	}
	return NoAuthMethod
}

func ReplyNoAuth(conn net.Conn) error {
	resp := new(MethodResp)
	resp.Method = NoAuthMethod
	return resp.Write(conn)
}

func ReplyUserAuth(conn net.Conn) error {
	resp := new(MethodResp)
	resp.Method = UserMethod
	if err := resp.Write(conn); err != nil {
		log.Fatal(err)
	}

	req := new(UserAuthReq)
	if err := req.Read(conn); err != nil {
		log.Fatal(err)
	}

	log.Println("userName ==> ", req.Version, req.Username, req.Password)

	userResp := new(UserAuthResp)
	userResp.Version = UserVer

	if req.Username == "admin" && req.Password == "123456" {
		userResp.Status = AuthSucceed
	} else {
		userResp.Status = AuthFailed
	}
	return userResp.Write(conn)
}

var authMap = map[byte]func(conn net.Conn) error{
	NoAuthMethod: ReplyNoAuth,
	UserMethod:   ReplyUserAuth,
}
