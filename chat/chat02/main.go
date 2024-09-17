package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func main() {

	path := "/Users/yurisa/Develop/GoWork/src/WorkSpace/Gotool/chat/chat01/html/upload.html"

	file, err := os.Open(path)
	if err != nil {
		tryCatch(err)
	}
	defer file.Close()

	bodyBuf := bytes.NewBuffer([]byte{})
	bodyWriter := multipart.NewWriter(bodyBuf)
	bodyWriter.WriteField("token", "vvv123")

	fileWrite, err := bodyWriter.CreateFormFile("uploadFile", filepath.Base(path))
	if err != nil {
		tryCatch(err)
	}

	if _, err = io.Copy(fileWrite, file); err != nil {
		tryCatch(err)
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post("http://localhost:9527/upload", contentType, bodyBuf)
	if err != nil {
		tryCatch(err)
	}
	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		tryCatch(err)
	}
	fmt.Println("resp => ", string(res))
}

func tryCatch(err error) {
	if err != nil {
		panic("error: " + err.Error())
	}
}
