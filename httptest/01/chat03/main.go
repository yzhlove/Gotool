package main

import (
	"log/slog"
	"net/http"
)

// 使用http.Head获取Body

func main() {

	resp, err := http.Head("http://localhost:1234")
	if err != nil {
		slog.Error("http.Head failed", slog.Any("error", err))
		return
	}
	slog.Info("resp info",
		slog.String("status", resp.Status),
		slog.Int("status code", resp.StatusCode),
		slog.Any("head", resp.Header))

}
