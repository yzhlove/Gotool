package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"time"
)

// 发送HEAD与超时控制

func main() {

	client := &http.Client{}
	ctx, cacnel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cacnel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:1234", nil)
	if err != nil {
		slog.Error("new request error", slog.Any("error", err))
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth("admin-test", "123456")
	req.AddCookie(&http.Cookie{
		Name:  "auth",
		Value: "yurisa",
		Path:  "/",
	})

	resp, err := client.Do(req)
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
