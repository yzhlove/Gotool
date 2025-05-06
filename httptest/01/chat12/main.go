package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
)

// 发送任意方法

func main() {

	client := &http.Client{}
	request, err := http.NewRequest(http.MethodDelete, "http://localhost:1234", nil)
	if err != nil {
		slog.Error("new request error", slog.Any("error", err))
		return
	}

	resp, err := client.Do(request)
	if err != nil {
		slog.Error("do request error", slog.Any("error", err))
		return
	}

	data, err := httputil.DumpResponse(resp, true)
	if err != nil {
		slog.Error("dump response error", slog.Any("error", err))
		return
	}
	fmt.Println(string(data))
}
