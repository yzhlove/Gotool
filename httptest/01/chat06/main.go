package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

// 使用POST发送实现io.Reader接口的数据
func main() {

	asang := strings.NewReader("给你的爱一直很安静\n来交换你偶尔给的关心")
	resp, err := http.Post("http://localhost:1234", "text/plain", asang)
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
