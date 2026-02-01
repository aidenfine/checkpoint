package main

import (
	"fmt"
	"net/http"

	checkpoint "github.com/aidenfine/checkpoint"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello World")
	})

	rlMiddleware := checkpoint.LimitByIp(20, 1, 1)
	rlHandler := rlMiddleware(mux)

	http.ListenAndServe(":8080", rlHandler)
}
