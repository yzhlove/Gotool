package main

import (
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

func main() {

	handle := http.FileServerFS(os.DirFS("."))
	http.Handle("/static/", http.StripPrefix("/static/", handle))
	http.Handle("/", stand{})

	slog.Info("http server starting...", slog.Int("port", 1234))
	if err := http.ListenAndServe(":1234", nil); err != nil {
		slog.Error("http server failed.", slog.Any("error", err))
	}

}

type stand struct{}

func (stand) ServeHTTP(resp http.ResponseWriter, req *http.Request) {

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("METHOD:%s \n", req.Method))
	sb.WriteString(fmt.Sprintf("HOST:%s \n", req.Host))
	sb.WriteString("HeadInfo:\n")
	for key, value := range req.Header {
		sb.WriteString(fmt.Sprintf("%s=%s\n", key, strings.Join(value, ",")))
	}

	if err := req.ParseForm(); err != nil {
		sb.WriteString(fmt.Sprintf("parse form data failed:%v\n", err))
	} else {
		sb.WriteString("FormInfo:\n")
		for key, value := range req.Form {
			sb.WriteString(fmt.Sprintf("%s=%s\n", key, strings.Join(value, ",")))
		}
	}

	req.ParseMultipartForm(1024 * 1024 * 10)
	file, fileHead, err := req.FormFile("upload")
	if err != nil {
		sb.WriteString(fmt.Sprintf("get upload file failed:%v\n", err))
	} else {
		defer file.Close()

		sb.WriteString("UploadFileInfo:\n")
		sb.WriteString(fmt.Sprintf("Filename:%s\n", fileHead.Filename))
		sb.WriteString(fmt.Sprintf("Header:%v\n", fileHead.Header))

		data, err := io.ReadAll(file)
		if err != nil {
			sb.WriteString(fmt.Sprintf("read upload file failed:%v", err))
		} else {
			sb.WriteString("UploadFileContent:\n")
			sb.WriteString(string(data))
			sb.WriteString("\n")
		}
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		sb.WriteString(fmt.Sprintf("read body failed:%v", err))
	} else {
		sb.WriteString("BodyInfo:\n")
		sb.WriteString(string(body))
	}

	sb.WriteString("\n")
	for _, cook := range req.Cookies() {
		sb.WriteString(fmt.Sprintf("CookieInfo:%s=%s\n", cook.Name, cook.Value))
	}

	http.SetCookie(resp, &http.Cookie{
		Name:  "cookValue",
		Value: fmt.Sprint(rand.Uint64()),
		Path:  "/",
	})
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte(sb.String()))
}
