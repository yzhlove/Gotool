package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"strings"
)

// Cookie的发送和接收

func main() {

	jar, err := cookiejar.New(nil)
	if err != nil {
		slog.Error("new cookiejar error", slog.Any("error", err))
		return
	}

	//cc := http.Client{Jar: jar}

	http.DefaultClient = &http.Client{Jar: jar}

	for i := 0; i < 3; i++ {
		//resp, err := cc.Get("http://localhost:1234")
		resp, err := http.Get("http://localhost:1234")
		if err != nil {
			slog.Error("get error", slog.Any("error", err))
			return
		}

		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			slog.Error("dump response error", slog.Any("error", err))
			return
		}
		fmt.Println(string(dump))
		fmt.Println(strings.Repeat("-", 50))
	}

}
