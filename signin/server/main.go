package main

import (
	"github.com/yzhlove/Gotool/signin/server/handler"
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/ike", handler.IkeHTTP)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
