package main

import (
	"fmt"
	"log"
	"net"
)

func main() {

	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("listen to ", l.Addr().String())

	for {
		c, err := l.Accept()
		if err != nil {
			log.Println("AcceptError: ", err)
			continue
		}
		go func(c net.Conn) {
			if err = Handel(c); err != nil {
				log.Println("HandelError: ", err)
			}
		}(c)
	}

}
