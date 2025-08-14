package main

import (
	"checkpoint"
	"fmt"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello World")
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello World")
	})

	rlMiddleware := checkpoint.LimitByEndpoint(5, 15*time.Second)
	rlHandler := rlMiddleware(mux)

	http.ListenAndServe(":8080", rlHandler)
}
