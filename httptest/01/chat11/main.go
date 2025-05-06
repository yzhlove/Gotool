package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
)

// 客户端访问文件系统

func main() {

	transport := &http.Transport{}
	transport.RegisterProtocol("file", http.NewFileTransport(http.Dir(".")))

	client := &http.Client{Transport: transport}
	resp, err := client.Get("file://./chat01/main.go")
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
