package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

// 使用POST发送x-www-form-urlencoded形式的表单

func main() {

	values := url.Values{
		"query": {"hello", "world"},
	}

	resp, err := http.PostForm("http://localhost:1234", values)
	if err != nil {
		slog.Error("http.PostForm failed", slog.Any("error", err))
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
