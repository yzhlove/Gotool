package main

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
)

// 使用multipart/form-data发送文件

func main() {

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	writer.WriteField("name", "yurisa")
	fileWrite, err := writer.CreateFormFile("upload", "main.go")
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
