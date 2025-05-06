package main

import (
	"fmt"
	"golang.org/x/net/idna"
	"log/slog"
)

// 国际化域名

func main() {

	src := "芙蓉王"
	resp, err := idna.ToASCII(src)
	if err != nil {
		slog.Error("idna.ToASCII failed", slog.Any("error", err))
		return
	}
	fmt.Println("resp => ", resp)
}
