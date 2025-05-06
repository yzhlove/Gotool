package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

// 使用GET发送查询请求

func main() {

	values := url.Values{
		"query": {"hello", "world"},
	}

	resp, err := http.Get("http://localhost:1234" + "?" + values.Encode())
	if err != nil {
		slog.Error("http.Get failed", slog.Any("error", err))
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("io.ReadAll failed", slog.Any("error", err))
		return
	}
	fmt.Println(string(body))
}
