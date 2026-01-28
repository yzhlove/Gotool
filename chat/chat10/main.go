package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/moby/moby/client"
)

func main() {
	apiClient, err := client.New(client.WithHost("unix:///app/docker.sock"))
	if err != nil {
		panic(err)
	}
	defer apiClient.Close()

	// HTTP handler：访问任意路径都返回 docker info
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		result, err := apiClient.Info(ctx, client.InfoOptions{})
		if err != nil {
			http.Error(w, "get docker info error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		if err := enc.Encode(result); err != nil {
			http.Error(w, "encode json error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	})

	addr := ":8081"
	log.Printf("listening on %s ...", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("http server error: %v", err)
	}

}
