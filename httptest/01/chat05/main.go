package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
)

// 使用POST方法发送任意请求

func main() {

	file, err := os.Open("./chat01/main.go")
	if err != nil {
		slog.Error("open file error", slog.Any("error", err))
		return
	}
	defer file.Close()

	resp, err := http.Post("http://localhost:1234", "text/plain", file)
	if err != nil {
		slog.Error("post error", slog.Any("error", err))
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("read body error", slog.Any("error", err))
		return
	}
	fmt.Println(string(data))
}
