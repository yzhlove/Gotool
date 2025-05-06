package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// 客户端代理

func main() {

	proxyUrl, err := url.Parse("http://localhost:1234")
	if err != nil {
		slog.Error("parse proxy url error", slog.Any("error", err))
		return
	}

	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}

	resp, err := client.Get("http://github.com")
	if err != nil {
		slog.Error("get error", slog.Any("error", err))
		return
	}
	data, err := httputil.DumpResponse(resp, true)
	if err != nil {
		slog.Error("dump response error", slog.Any("error", err))
		return
	}
	fmt.Println(string(data))
}
