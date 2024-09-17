package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const DirRoot = "/Users/yurisa/Develop/GoWork/src/WorkSpace/Gotool/testFiles/"

func main() {

	http.HandleFunc("/upload", handler)
	http.ListenAndServe(":9527", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseMultipartForm(1024 * 1024 * 1024); err != nil {
		writeError(w, err)
		return
	}

	f, h, err := r.FormFile("uploadFile")
	if err != nil {
		writeError(w, err)
		return
	}
	defer f.Close()

	fmt.Printf("multipart.FileHead => %+v \n", h)
	fmt.Fprintf(w, "multipart.FileHead %+v \n", h)

	if err = os.MkdirAll(filepath.Join(DirRoot, "chat01File"), os.ModePerm); err != nil {
		writeError(w, err)
		return
	}

	file, err := os.OpenFile(filepath.Join(DirRoot, "chat01File", h.Filename), os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		writeError(w, err)
		return
	}
	defer file.Close()
	io.Copy(file, f)
}

func writeError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
