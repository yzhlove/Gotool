package main

import (
	"io"
	"net"
	"testing"
)

func Test_ClientByAuth(t *testing.T) {

	conn, err := net.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		t.Fatal(err)
	}

	resp := new(MethodsResp)
	resp.Methods = []byte{NoAuthMethod, UserMethod}
	if err := resp.Write(conn); err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 2)
	if _, err := io.ReadFull(conn, buf); err != nil {
		t.Fatal(err)
	}

	if int(buf[0]) != SocksVer {
		t.Fatal("version not support")
	}

	switch buf[1] {
	case NoAuthMethod:
		t.Log("no auth")
		conn.Close()
	case UserMethod:
		t.Log("userAuth")

		req := new(UserAuthReq)
		req.Version = UserVer
		req.Username = "admin"
		req.Password = "123456"

		if err := req.Write(conn); err != nil {
			t.Fatal(err)
		}

		if _, err := io.ReadFull(conn, buf); err != nil {
			t.Fatal(err)
		}

		if int(buf[0]) != UserVer {
			t.Fatal("version not support")
		}

		if int(buf[1]) == AuthSucceed {
			t.Log("Auth Succeed! ")
		} else {
			t.Log("Auth Failed! ")
		}
	}

}
