package main

import (
	"log"

	"github.com/yzhlove/Gotool/signin/client/api"
	"github.com/yzhlove/Gotool/signin/client/context"
	"github.com/yzhlove/Gotool/signin/client/http"
)

func main() {

	http.New()

	address := "http://localhost:8080"

	ctx := context.New()
	if err := api.Ike(ctx, address); err != nil {
		log.Fatal(err)
	}

}
