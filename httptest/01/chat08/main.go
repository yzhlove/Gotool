package main

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
)

// 对文件发送设置任意MIME类型

func main() {

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	part := make(textproto.MIMEHeader)
	part.Set("Content-Type", "text/plain")
	part.Set("Content-Disposition", `form-data; name="upload"; filename="main.go"`)
	fileWrite, err := writer.CreatePart(part)

	if err != nil {
		slog.Error("create form file error", slog.Any("error", err))
		return
	}

	readFile, err := os.Open("./chat01/main.go")
	if err != nil {
		slog.Error("open file error", slog.Any("error", err))
		return
	}
	defer readFile.Close()
	io.Copy(fileWrite, readFile)
	writer.Close()

	resp, err := http.Post("http://localhost:1234", writer.FormDataContentType(), &buffer)
	if err != nil {
		slog.Error("post error", slog.Any("error", err))
		return
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	fmt.Println(string(data))
}
